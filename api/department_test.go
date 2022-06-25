package api

import (
	"bytes"
	"encoding/json"
	mockdb "example/employee/server/db/mock"
	db "example/employee/server/db/sqlc"
	"example/employee/server/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomDepartment() db.Department {
	return db.Department{
		DepartmentID:   util.RandomInt64(1, 1000),
		DepartmentName: util.RandomDepartmentName(),
	}
}
func TestPostDepartmentAPI(t *testing.T) {
	department := randomDepartment()

	testCases := []struct {
		name                    string
		createDepartmentRequest CreateDepartmentRequest
		buildStubs              func(store *mockdb.MockStore)
		checkResponse           func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			createDepartmentRequest: CreateDepartmentRequest{
				DepartmentName: department.DepartmentName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDepartment(gomock.Any(), gomock.Eq(department.DepartmentName)).
					Times(1).
					Return(department, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchDepartment(t, recorder.Body, department)
			},
		},
		{
			name: "BadRequest",
			createDepartmentRequest: CreateDepartmentRequest{
				DepartmentName: department.DepartmentName,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateDepartment(gomock.Any(), gomock.Eq(department.DepartmentName)).
					Times(1).
					Return(department, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchDepartment(t, recorder.Body, department)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)

			store := mockdb.NewMockStore(controller)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			jsonBytes, err := json.Marshal(tc.createDepartmentRequest)
			url := fmt.Sprintf("/departments")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func requiredBodyMatchDepartment(t *testing.T, body *bytes.Buffer, department db.Department) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotDepartment db.Department
	err = json.Unmarshal(data, &gotDepartment)
	require.NoError(t, err)
	require.Equal(t, department, gotDepartment)
}
