package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/api/internal/mocks"
	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewBaseAPI(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		handler func() *mocks.HandlerMock
		asserts func(*testing.T, *api.BaseAPI)
	}{
		{
			name: "create new base api with a route",
			handler: func() *mocks.HandlerMock {
				hm := mocks.NewHandlerMock(t)
				hm.On("URL").Return("/test")
				hm.On("Method").Return("GET")
				hm.On("Respond", mock.AnythingOfType("*gin.Context")).Return(nil)
				return hm
			},
			asserts: func(t *testing.T, ba *api.BaseAPI) {
				rr := httptest.NewRecorder()
				httpReq := httptest.NewRequest(http.MethodGet, "/test", nil)
				ba.Engine().ServeHTTP(rr, httpReq)
				require.Equal(t, rr.Code, http.StatusOK)
			},
		},
		{
			name: "create new base api with a wildcard route",
			handler: func() *mocks.HandlerMock {
				hm := mocks.NewHandlerMock(t)
				hm.On("URL").Return("/test/*")
				hm.On("Method").Return("POST")
				hm.On("Respond", mock.AnythingOfType("*gin.Context")).Return(nil)
				return hm
			},
			asserts: func(t *testing.T, ba *api.BaseAPI) {
				rr := httptest.NewRecorder()
				baseReq := httptest.NewRequest(http.MethodPost, "/test/anything", nil)
				ba.Engine().ServeHTTP(rr, baseReq)
				require.Equal(t, http.StatusOK, rr.Code)
				rr = httptest.NewRecorder()
				multilevelReq := httptest.NewRequest(http.MethodPost, "/test/anything/else", nil)
				ba.Engine().ServeHTTP(rr, multilevelReq)
				require.Equal(t, http.StatusOK, rr.Code)
				rr = httptest.NewRecorder()
				notFoundReq := httptest.NewRequest(http.MethodPost, "/asd", nil)
				ba.Engine().ServeHTTP(rr, notFoundReq)
				require.Equal(t, http.StatusNotFound, rr.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			e := gin.New()
			ba := api.NewBaseAPI(e, logger.NewTestLogger())
			ba.AddEndpoints([]api.Handler{tc.handler()})
			tc.asserts(t, ba)
		})
	}
}
