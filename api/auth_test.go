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

	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func loginUserMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var resp LoginResponse
	err = json.Unmarshal(data, &resp)

	require.NoError(t, err)

	require.Equal(t, user.Email, resp.User.Email)
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		Status:         db.UserStatusActivated,
	}
	return
}
func TestLoginAPI(t *testing.T) {
	user, password := randomUser(t)

	pendingUser := user
	pendingUser.Status = db.UserStatusDeactivated

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				loginUserMatch(t, recorder.Body, user)
			},
		},
		{
			name: "IncorrectPassword Should Return 401",
			body: gin.H{
				"email":    user.Email,
				"password": password + "abc",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Login with not activated user Should Return 400",
			body: gin.H{
				"email":    user.Email,
				"password": password,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(pendingUser, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			testCase.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			res, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			userURL := "/auth/login"
			request, err := http.NewRequest(http.MethodPost, userURL, bytes.NewReader(res))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func randomPendingUser(t *testing.T) db.User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)
	return db.User{
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassword,
		Status:         db.UserStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
func TestActivateUserAPI(t *testing.T) {
	pendingUser := randomPendingUser(t)

	activatedUser := pendingUser
	activatedUser.Status = db.UserStatusActivated
	activatedUser.UpdatedAt = time.Now().Add(time.Minute)

	newPassword := "newPassword"
	testCases := []struct {
		name                string
		activateUserReuqest ActivateUserRequest
		buildStubs          func(store *mockdb.MockStore)
		checkResponse       func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			activateUserReuqest: ActivateUserRequest{
				Email:    pendingUser.Email,
				Password: newPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(pendingUser.Email)).
					Times(1).
					Return(pendingUser, nil)
				store.EXPECT().ActivateUser(gomock.Any(), gomock.Any()).Times(1)
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(pendingUser.Email)).
					Times(1).
					Return(activatedUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requiredBodyMatchActiveUser(t, recorder.Body, pendingUser)
			},
		},
		{
			name: "Request non-pending user should got 400",
			activateUserReuqest: ActivateUserRequest{
				Email:    activatedUser.Email,
				Password: newPassword,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(activatedUser.Email)).
					Times(1).
					Return(activatedUser, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			jsonBytes, err := json.Marshal(tc.activateUserReuqest)
			url := fmt.Sprintf("/auth/activate")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))

			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func requiredBodyMatchActiveUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var resp ActivateUserResponse
	err = json.Unmarshal(data, &resp)
	require.NoError(t, err)
	require.Equal(t, user.Email, resp.Email)
	require.Equal(t, db.UserStatusActivated, resp.Status)
	require.True(t, resp.UpdatedAt.After(user.UpdatedAt))
}
