package api_test

import (
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

func TestDynamicHandler_Respond(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		valuer         func(*gin.Context) *mocks.ValuerMock
		responseFormat enums.ResponseFormat
		expectedStatus int
		expectedError  error
	}{
		{
			name:           "success response with json",
			responseFormat: enums.ResponseFormats.JSON(),
			valuer: func(c *gin.Context) *mocks.ValuerMock {
				v := mocks.NewValuerMock(t)
				v.On("Generate", c).Return(map[string]string{"key": "value"}, nil)
				return v
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error response with failed locating element",
			responseFormat: enums.ResponseFormats.JSON(),
			valuer: func(c *gin.Context) *mocks.ValuerMock {
				v := mocks.NewValuerMock(t)
				v.On("Generate", c).Return(nil, values.ErrFailedLocatingElement)
				return v
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "error response with conversion failed",
			responseFormat: enums.ResponseFormats.JSON(),
			valuer: func(c *gin.Context) *mocks.ValuerMock {
				v := mocks.NewValuerMock(t)
				v.On("Generate", c).Return(nil, values.ErrConversionFailed)
				return v
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "erorr creating dynamic handler",
			valuer: func(*gin.Context) *mocks.ValuerMock {
				return nil
			},
			expectedError:  api.ErrNotSupportedFormat,
			responseFormat: enums.ResponseFormat("0"),
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rr := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rr)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
			dh, err := api.NewDynamicHandler(
				tc.responseFormat,
				http.MethodGet,
				"/test",
				http.StatusOK,
				tc.valuer(c),
				logger.NewTestLogger(),
			)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				dh.Respond(c)
				require.Equal(t, rr.Code, tc.expectedStatus)
			}
		})
	}
}
