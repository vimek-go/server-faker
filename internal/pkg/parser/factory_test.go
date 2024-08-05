package parser_test

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/parser"
	"github.com/vimek-go/server-faker/internal/pkg/parser/dto"
	"github.com/vimek-go/server-faker/internal/pkg/tools"

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

func TestFactory_CreateResponseEndpoint(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		endpoint      dto.Endpoint
		asserts       func(*testing.T, api.ResponseHandler)
		setup         func(*testing.T, string)
		expectedError string
	}{
		{
			name: "static response endpoint creation",
			endpoint: dto.Endpoint{
				Method: http.MethodPost,
				URL:    "/test",
				Response: &dto.Response{
					Type:   enums.ResponseTypes.Static(),
					Status: http.StatusOK,
					File:   "test.json",
					Format: enums.ResponseFormats.JSON(),
				},
			},
			setup: func(t *testing.T, dir string) {
				filepath := path.Join(dir, "test.json")
				tools.SaveToAFile(t, `{"test": "value"}`, filepath)
			},
			asserts: func(t *testing.T, handler api.ResponseHandler) {
				require.NotNil(t, handler)
				require.Equal(t, http.MethodPost, handler.Method())
				require.Equal(t, "/test", handler.URL())
				require.Equal(t, http.StatusOK, handler.ReturnCode())
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				handler.Respond(c)
				require.Equal(t, http.StatusOK, rr.Code)
				require.Equal(t, `{"test": "value"}`, rr.Body.String())
			},
		},
		{
			name: "error on static response endpoint creation",
			endpoint: dto.Endpoint{
				Method: http.MethodPost,
				URL:    "/test",
				Response: &dto.Response{
					Type:   enums.ResponseTypes.Static(),
					Status: http.StatusOK,
					File:   "test.json",
					Format: enums.ResponseFormats.JSON(),
				},
			},
			setup: func(*testing.T, string) {
				// file not created
			},
			expectedError: "no such file or directory",
		},
		{
			name: "dynamic response endpoint creation",
			endpoint: dto.Endpoint{
				Method: http.MethodGet,
				URL:    "/test",
				Response: &dto.Response{
					Type:   enums.ResponseTypes.Dynamic(),
					Status: http.StatusCreated,
					Format: enums.ResponseFormats.JSON(),
					Object: dto.Params{
						dto.Param{
							Key: "array",
							Array: &dto.Array{
								Min: 1,
								Max: 1,
								Element: []dto.Param{
									{
										Static: &dto.Static{
											Value: "value",
										},
									},
								},
							},
						},
					},
				},
			},
			setup: func(*testing.T, string) {},
			asserts: func(t *testing.T, handler api.ResponseHandler) {
				require.NotNil(t, handler)
				require.Equal(t, http.MethodGet, handler.Method())
				require.Equal(t, "/test", handler.URL())
				require.Equal(t, http.StatusCreated, handler.ReturnCode())
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				handler.Respond(c)
				require.Equal(t, http.StatusCreated, rr.Code)
				require.Equal(t, `{"array":["value"]}`, rr.Body.String())
			},
		},
		{
			name: "dynamic response endpoint with object creation",
			endpoint: dto.Endpoint{
				Method: http.MethodGet,
				URL:    "/test",
				Response: &dto.Response{
					Type:   enums.ResponseTypes.Dynamic(),
					Status: http.StatusCreated,
					Format: enums.ResponseFormats.JSON(),
					Object: dto.Params{
						dto.Param{
							Key: "name",
							Static: &dto.Static{
								Value: "name",
							},
						},
						dto.Param{
							Key: "value",
							Static: &dto.Static{
								Value: "value",
							},
						},
					},
				},
			},
			setup: func(*testing.T, string) {},
			asserts: func(t *testing.T, handler api.ResponseHandler) {
				require.NotNil(t, handler)
				require.Equal(t, http.MethodGet, handler.Method())
				require.Equal(t, "/test", handler.URL())
				require.Equal(t, http.StatusCreated, handler.ReturnCode())
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				handler.Respond(c)
				require.Equal(t, http.StatusCreated, rr.Code)
				require.Equal(t, `{"name":"name","value":"value"}`, rr.Body.String())
			},
		},
		{
			name: "dynamic response endpoint with object creation error",
			endpoint: dto.Endpoint{
				Method: http.MethodGet,
				URL:    "/test",
				Response: &dto.Response{
					Type:   enums.ResponseTypes.Dynamic(),
					Status: http.StatusCreated,
					Format: enums.ResponseFormats.JSON(),
					Object: dto.Params{
						dto.Param{
							Key: "name",
							Static: &dto.Static{
								Value: "name",
							},
						},
						dto.Param{
							Static: &dto.Static{
								Value: "value",
							},
						},
					},
				},
			},
			setup:         func(*testing.T, string) {},
			expectedError: "empty param key",
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := parser.NewFactory(nil, logger.NewTestLogger())
			dir := t.TempDir()
			tc.setup(t, dir)
			handler, err := f.CreateResponseEndpoint(tc.endpoint, dir)
			if len(tc.expectedError) > 0 {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				tc.asserts(t, handler)
			}
		})
	}
}

func TestFactory_CreateProxyEndpoint(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusOK)
	}))
	testCases := []struct {
		name          string
		endpoint      dto.Endpoint
		asserts       func(*testing.T, api.Handler)
		expectedError error
	}{
		{
			name: "static proxy endpoint creation",
			endpoint: dto.Endpoint{
				Method: http.MethodPost,
				URL:    "/test",
				Proxy: &dto.Proxy{
					Type:   enums.ResponseTypes.Static(),
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
			asserts: func(t *testing.T, handler api.Handler) {
				require.NotNil(t, handler)
				require.Equal(t, http.MethodPost, handler.Method())
				require.Equal(t, "/test", handler.URL())
				rr := CreateTestResponseRecorder()
				c, _ := gin.CreateTestContext(rr)
				c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
				handler.Respond(c)
				require.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			name: "dynamic proxy endpoint creation",
			endpoint: dto.Endpoint{
				Method: http.MethodPatch,
				URL:    "/test",
				Proxy: &dto.Proxy{
					Type:   enums.ResponseTypes.Dynamic(),
					Method: http.MethodGet,
					URL:    server.URL,
					URLParams: dto.Params{
						dto.Param{
							Key: "url",
							Static: &dto.Static{
								Value: "url",
							},
						},
					},
					Query: dto.Params{
						dto.Param{
							Key: "query",
							Mapped: &dto.Mapped{
								From:  "query",
								Param: "query",
							},
						},
					},
					Object: dto.Params{
						dto.Param{
							Key: "object",
							Random: &dto.Random{
								Min:  1,
								Max:  1,
								Type: enums.RandomKinds.Boolean().String(),
							},
						},
					},
				},
			},
			asserts: func(t *testing.T, handler api.Handler) {
				require.NotNil(t, handler)
				require.Equal(t, http.MethodPatch, handler.Method())
				require.Equal(t, "/test", handler.URL())
				rr := CreateTestResponseRecorder()
				c, _ := gin.CreateTestContext(rr)
				c.Request = httptest.NewRequest(http.MethodPatch, "/test?query=test", nil)
				handler.Respond(c)
				require.Equal(t, http.StatusOK, rr.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			f := parser.NewFactory(nil, logger.NewTestLogger())
			handler, err := f.CreateProxyEndpoint(tc.endpoint)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				tc.asserts(t, handler)
			}
		})
	}
}
