package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStaticProxy_Respond(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		proxyMethod string
		headers     map[string]string
		asserts     func(t *testing.T, req *http.Request)
	}{
		{
			name:        "success response with headers",
			proxyMethod: http.MethodGet,
			headers:     map[string]string{"header": "value"},
			asserts: func(t *testing.T, req *http.Request) {
				require.Equal(t, req.Header.Get("header"), "value")
				require.Equal(t, req.Method, http.MethodGet)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rr := CreateTestResponseRecorder()
			c, _ := gin.CreateTestContext(rr)
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Send response to be tested
				tc.asserts(t, req)
				rw.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			sp := api.NewStaticProxy(
				http.MethodPost,
				"/test",
				http.MethodGet,
				server.URL,
				tc.headers,
				logger.NewTestLogger(),
			)
			c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
			sp.Respond(c)
			require.Equal(t, c.Request.Method, http.MethodPost)
		})
	}
}
