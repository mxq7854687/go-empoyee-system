package api

import (
	"example/employee/server/token"
	"example/employee/server/util"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuth(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authType string,
	email string,
	platform token.Platform,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(email, platform, duration)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authType, token)
	request.Header.Set(authHeaderKey, authHeader)
}

func addMockAuthBearer(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
) {
	token, err := tokenMaker.CreateToken(util.RandomEmail(), token.Web, time.Minute)
	require.NoError(t, err)

	authHeader := fmt.Sprintf("%s %s", authTypeBearer, token)
	request.Header.Set(authHeaderKey, authHeader)
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Happy Case",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, "testEmail", token.Web, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Token Expired",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authTypeBearer, "testEmail", token.Web, -time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "No Auth",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Auth not supported",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, "other", "testEmail", token.Web, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "Invalid Auth",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, "", "testEmail", token.Web, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		currTest := testCases[i]

		t.Run(currTest.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authPath := "/auth"
			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			currTest.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			currTest.checkResponse(t, recorder)
		})
	}
}
