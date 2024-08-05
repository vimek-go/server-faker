package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/api/internal/mocks"
	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func CreateTestResponseRecorder() *TestResponseRecorder {
	return &TestResponseRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func TestDynamicProxy_Respond(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		proxyMethod    string
		queryValuer    func(*gin.Context) map[string]values.Valuer
		urlValuer      func(*gin.Context) map[string]values.Valuer
		payloadValuer  func(*gin.Context) *mocks.ValuerMock
		headers        map[string]string
		asserts        func(t *testing.T, req *http.Request)
		expectedStatus int
	}{
		{
			name:        "success response with prepared url and query",
			proxyMethod: http.MethodPost,
			queryValuer: func(c *gin.Context) map[string]values.Valuer {
				mv := mocks.NewValuerMock(t)
				mv.On("Generate", c).Return("query-value", nil).Once()
				mv.On("Type").Return(enums.GenerationTypes.SingleValue()).Once()
				return map[string]values.Valuer{"query": mv}
			},
			urlValuer: func(c *gin.Context) map[string]values.Valuer {
				mv := mocks.NewValuerMock(t)
				mv.On("Generate", c).Return("url-value", nil).Once()
				return map[string]values.Valuer{"url": mv}
			},
			payloadValuer: func(c *gin.Context) *mocks.ValuerMock {
				mv := mocks.NewValuerMock(t)
				mv.On("IsNil").Return(false)
				mv.On("Generate", c).Return(map[string]interface{}{"payload": "value"}, nil).Once()
				return mv
			},
			asserts: func(t *testing.T, req *http.Request) {
				require.Contains(t, req.URL.Path, "url-value")
				require.Equal(t, req.URL.RawQuery, "query=query-value")
				reader := req.Body
				buf := make([]byte, 19)
				n, err := reader.Read(buf)
				require.NoError(t, err)
				require.Equal(t, string(buf[:n]), `{"payload":"value"}`)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failes to prepare url",
			queryValuer: func(*gin.Context) map[string]values.Valuer {
				return nil
			},
			urlValuer: func(c *gin.Context) map[string]values.Valuer {
				mv := mocks.NewValuerMock(t)
				mv.On("Generate", c).Return(map[string]interface{}{"payload": "value"}, nil).Once()
				return map[string]values.Valuer{"url": mv}
			},
			payloadValuer: func(*gin.Context) *mocks.ValuerMock {
				return mocks.NewValuerMock(t)
			},
			asserts:        func(*testing.T, *http.Request) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failes to prepare payload",
			queryValuer: func(*gin.Context) map[string]values.Valuer {
				return nil
			},
			urlValuer: func(*gin.Context) map[string]values.Valuer {
				return nil
			},
			payloadValuer: func(c *gin.Context) *mocks.ValuerMock {
				mv := mocks.NewValuerMock(t)
				mv.On("IsNil").Return(false)
				mv.On("Generate", c).Return(nil, fmt.Errorf("error")).Once()
				return mv
			},
			asserts:        func(*testing.T, *http.Request) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "success to prepare query from multiple values",
			queryValuer: func(c *gin.Context) map[string]values.Valuer {
				mv := mocks.NewValuerMock(t)
				mv.On("Generate", c).Return([]any{1, 2}, nil).Once()
				mv.On("Type").Return(enums.GenerationTypes.MultiValue()).Once()
				return map[string]values.Valuer{"query": mv}
			},
			urlValuer: func(*gin.Context) map[string]values.Valuer {
				return nil
			},
			payloadValuer: func(*gin.Context) *mocks.ValuerMock {
				mv := mocks.NewValuerMock(t)
				mv.On("IsNil").Return(true)
				return mv
			},
			asserts: func(t *testing.T, req *http.Request) {
				require.Equal(t, req.URL.RawQuery, "query=1&query=2")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "success passing headers",
			queryValuer: func(*gin.Context) map[string]values.Valuer {
				return nil
			},
			urlValuer: func(*gin.Context) map[string]values.Valuer {
				return nil
			},
			payloadValuer: func(*gin.Context) *mocks.ValuerMock {
				mv := mocks.NewValuerMock(t)
				mv.On("IsNil").Return(true)
				return mv
			},
			headers: map[string]string{"header": "value"},
			asserts: func(t *testing.T, req *http.Request) {
				require.Equal(t, req.Header.Get("header"), "value")
			},
			expectedStatus: http.StatusOK,
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
			dp := api.NewDynamicProxy(
				http.MethodPost,
				"/test",
				http.MethodPost,
				server.URL+"/:url",
				tc.urlValuer(c),
				tc.queryValuer(c),
				tc.payloadValuer(c),
				tc.headers,
				logger.NewTestLogger(),
			)
			c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)

			dp.Respond(c)
			require.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
