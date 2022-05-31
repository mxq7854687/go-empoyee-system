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

			server := NewServer(store)
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

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/jobs/%d", currentTest.jobID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
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
