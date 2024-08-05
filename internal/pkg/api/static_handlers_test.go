package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestStaticHandler_Respond(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		responseFormat enums.ResponseFormat
		method         string
		response       []byte
		contentType    string
		expectedError  error
	}{
		{
			name:           "success response with json",
			responseFormat: enums.ResponseFormats.JSON(),
			method:         "GET",
			response:       []byte(`{"key": "value"}`),
		},
		{
			name:           "success response with xml",
			responseFormat: enums.ResponseFormats.XML(),
			method:         "GET",
			response:       []byte(`<key>value</key>`),
		},
		{
			name:           "success response with bytes",
			responseFormat: enums.ResponseFormats.Bytes(),
			method:         "GET",
			response:       []byte(`plain text`),
			contentType:    "text/plain",
		},
		{
			name:           "error response with empty content type",
			responseFormat: enums.ResponseFormats.Bytes(),
			expectedError:  api.ErrContentTypeEmpty,
		},
		{
			name:           "error response with unsupported response format",
			responseFormat: enums.ResponseFormat("unsupported"),
			expectedError:  api.ErrNotSupportedResponseFormat,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			sh, err := api.NewStaticHandler(
				tc.responseFormat,
				tc.method,
				"/test",
				http.StatusOK,
				tc.response,
				tc.contentType,
				logger.NewTestLogger(),
			)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				rr := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(rr)
				sh.Respond(c)
				require.Equal(t, rr.Code, http.StatusOK)
				require.Equal(t, rr.Body.String(), string(tc.response))
			}
		})
	}
}
