package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "example/employee/server/db/mock"
	db "example/employee/server/db/sqlc"
	"example/employee/server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateJob(t *testing.T) {
	job := randomJob()

	testCases := []struct {
		name          string
		body          gin.H
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"job_title":  job.JobTitle,
				"min_salary": job.MinSalary,
				"max_salary": job.MaxSalary,
			},
			buildStub: func(store *mockdb.MockStore) {
				args := db.CreateJobParams{
					JobTitle:  job.JobTitle,
					MinSalary: job.MinSalary,
					MaxSalary: job.MaxSalary,
				}

				store.EXPECT().
					CreateJob(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(job, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchJob(t, recorder.Body, job)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"job_title":  job.JobTitle,
				"min_salary": job.MinSalary,
				"max_salary": job.MaxSalary,
			},
			buildStub: func(store *mockdb.MockStore) {
				args := db.CreateJobParams{
					JobTitle:  job.JobTitle,
					MinSalary: job.MinSalary,
					MaxSalary: job.MaxSalary,
				}

				store.EXPECT().
					CreateJob(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(db.Job{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest-Invalid body",
			body: gin.H{},
			buildStub: func(store *mockdb.MockStore) {
				args := db.CreateJobParams{
					JobTitle:  job.JobTitle,
					MinSalary: job.MinSalary,
					MaxSalary: job.MaxSalary,
				}

				store.EXPECT().
					CreateJob(gomock.Any(), gomock.Eq(args)).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			currentTest.buildStub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			//Marshal body data to JSON
			data, err := json.Marshal(currentTest.body)
			require.NoError(t, err)

			url := "/jobs"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestGetJob(t *testing.T) {
	job := randomJob()

	testCases := []struct {
		name          string
		jobID         int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			jobID: job.JobID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetJob(gomock.Any(), gomock.Eq(job.JobID)).
					Times(1).
					Return(job, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchJob(t, recorder.Body, job)
			},
		},
		{
			name:  "NotFound",
			jobID: job.JobID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetJob(gomock.Any(), gomock.Eq(job.JobID)).
					Times(1).
					Return(db.Job{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "InternalError",
			jobID: job.JobID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetJob(gomock.Any(), gomock.Eq(job.JobID)).
					Times(1).
					Return(db.Job{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "BadRequest-Invalid ID",
			jobID: 0,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetJob(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			currentTest.buildStub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/jobs/%d", currentTest.jobID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestListJobs(t *testing.T) {
	n := 5
	jobs := make([]db.Job, n)

	for i := 0; i < n; i++ {
		jobs[i] = randomJob()
	}

	type QueryParam struct {
		pageID   int32
		pageSize int32
	}

	testCases := []struct {
		name          string
		query         QueryParam
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: QueryParam{
				pageID:   1,
				pageSize: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.ListJobsParams{
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().
					ListJobs(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(jobs, nil)

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchJobs(t, recorder.Body, jobs)
			},
		},
		{
			name: "InternalError",
			query: QueryParam{
				pageID:   1,
				pageSize: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListJobs(gomock.Any(), gomock.Any()).
					Times(1).
					Return([]db.Job{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest-Invalid page ID",
			query: QueryParam{
				pageID:   -1,
				pageSize: int32(n),
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListJobs(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "BadRequest-Invalid page size",
			query: QueryParam{
				pageID:   1,
				pageSize: 10000,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListJobs(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			currentTest.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := "/jobs"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// query params
			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", currentTest.query.pageID))
			q.Add("page_size", fmt.Sprintf("%d", currentTest.query.pageSize))
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(recorder)
		})
	}
}

func TestUpdateJob(t *testing.T) {
	job := randomJob()
	updatedJobTitle := util.RandomJobTitle()

	testCases := []struct {
		name          string
		jobID         int64
		body          gin.H
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			jobID: job.JobID,
			body: gin.H{
				"job_title":  updatedJobTitle,
				"min_salary": job.MinSalary,
				"max_salary": job.MaxSalary,
			},
			buildStub: func(store *mockdb.MockStore) {

				args := db.UpdateJobParams{
					JobID:     job.JobID,
					JobTitle:  updatedJobTitle,
					MinSalary: job.MinSalary,
					MaxSalary: job.MaxSalary,
				}

				store.EXPECT().
					UpdateJob(gomock.Any(), gomock.Eq(args)).
					Times(1).
					Return(job, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchJob(t, recorder.Body, job)
			},
		},
		{
			name:  "InternalError",
			jobID: job.JobID,
			body: gin.H{
				"job_title":  updatedJobTitle,
				"min_salary": job.MinSalary,
				"max_salary": job.MaxSalary,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateJob(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Job{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "BadRequest-Invalid ID",
			jobID: -1,
			body: gin.H{
				"job_title":  updatedJobTitle,
				"min_salary": job.MinSalary,
				"max_salary": job.MaxSalary,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateJob(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "BadRequest-Invalid Body",
			jobID: job.JobID,
			body: gin.H{
				"job_title":  -1,
				"min_salary": job.MinSalary,
				"max_salary": job.MaxSalary,
			},
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateJob(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			currentTest.buildStub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			//Marshal body data to JSON
			data, err := json.Marshal(currentTest.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/jobs/%d", currentTest.jobID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func TestDeleteJob(t *testing.T) {
	job := randomJob()

	testCases := []struct {
		name          string
		jobID         int64
		buildStub     func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			jobID: job.JobID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteJob(gomock.Any(), gomock.Eq(job.JobID)).
					After(store.EXPECT().
						GetJob(gomock.Any(), gomock.Eq(job.JobID))).
					Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "BadRequest",
			jobID: -1,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteJob(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "NotFound",
			jobID: job.JobID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteJob(gomock.Any(), gomock.Eq(job.JobID)).
					Times(0).
					Return(sql.ErrNoRows).
					After(store.EXPECT().GetJob(gomock.Any(), gomock.Eq(job.JobID)).Times(1).Return(db.Job{}, sql.ErrNoRows))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:  "InternalError - from GetJob",
			jobID: job.JobID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteJob(gomock.Any(), gomock.Eq(job.JobID)).
					Times(0).
					Return(sql.ErrConnDone).
					After(store.EXPECT().GetJob(gomock.Any(), gomock.Eq((job.JobID))).Times(1).Return(db.Job{}, sql.ErrConnDone))
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:  "InternalError - from DeleteJob",
			jobID: job.JobID,
			buildStub: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteJob(gomock.Any(), gomock.Eq(job.JobID)).
					Times(1).
					After(store.EXPECT().GetJob(gomock.Any(), gomock.Eq((job.JobID))).Times(1).Return(job, nil)).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currentTest := testCases[i]

		t.Run(currentTest.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			currentTest.buildStub(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/jobs/%d", currentTest.jobID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			currentTest.checkResponse(t, recorder)
		})
	}
}

func randomJob() db.Job {
	return db.Job{
		JobID:     util.RandomInt64(1, 1000),
		JobTitle:  util.RandomJobTitle(),
		MinSalary: sql.NullInt64{util.RandomMinSalary(), true},
		MaxSalary: sql.NullInt64{util.RandomMaxSalary(), true},
	}
}

func requireBodyMatchJob(t *testing.T, responseBody *bytes.Buffer, job db.Job) {
	data, err := ioutil.ReadAll(responseBody)
	require.NoError(t, err)

	var gotJob db.Job
	err = json.Unmarshal(data, &gotJob)
	require.NoError(t, err)
	require.Equal(t, job, gotJob)
}

func requireBodyMatchJobs(t *testing.T, responseBody *bytes.Buffer, job []db.Job) {
	data, err := ioutil.ReadAll(responseBody)
	require.NoError(t, err)

	var gotJobs []db.Job
	err = json.Unmarshal(data, &gotJobs)
	require.NoError(t, err)
	require.Equal(t, job, gotJobs)
}
