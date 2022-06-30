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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func randomePendingUser(t *testing.T) db.User {
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
	pendingUser := randomePendingUser(t)

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

			server := NewServer(store)
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
