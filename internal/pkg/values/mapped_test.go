package values_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/enums"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/values"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestPayloadMapper_Generate(t *testing.T) {
	t.Parallel()
	const testKey = "test-key"
	testCases := []struct {
		name          string
		payload       string
		jsonPath      string
		responseKey   string
		valueKey      string
		expectedError error
		expected      any
		conv          enums.ConversionType
	}{
		{
			name:          "empty path",
			payload:       `{}`,
			jsonPath:      "",
			valueKey:      testKey,
			expectedError: values.ErrFailedLocatingElement,
			conv:          enums.ConversionTypes.None(),
		},
		{
			name: "correct path no conversion",
			payload: `{
				"name": 1323
			}`,
			jsonPath: "$.name",
			valueKey: testKey,
			expected: float64(1323),
			conv:     enums.ConversionTypes.None(),
		},
		{
			name: "correct path from int to int",
			payload: `{
				"name": 1323
			}`,
			jsonPath: "$.name",
			valueKey: testKey,
			expected: float64(1323),
			conv:     enums.ConversionTypes.Number(),
		},
		{
			name: "correct path from int to string",
			payload: `{
				"name": 1323
			}`,
			jsonPath: "$.name",
			valueKey: testKey,
			expected: "1323",
			conv:     enums.ConversionTypes.Text(),
		},
		{
			name: "correct path from string to int",
			payload: `{
				"name": "1323"
			}`,
			jsonPath: "$.name",
			valueKey: testKey,
			expected: float64(1323),
			conv:     enums.ConversionTypes.Number(),
		},
		{
			name: "correct path, conversion string to int fails",
			payload: `{
				"name": "1323as"
			}`,
			jsonPath:      "$.name",
			valueKey:      testKey,
			expectedError: values.ErrConversionFailed,
			conv:          enums.ConversionTypes.Number(),
		},
		{
			name: "path not present",
			payload: `{
				"name": 123
			}`,
			jsonPath:      "$.name2",
			expectedError: values.ErrFailedLocatingElement,
			conv:          enums.ConversionTypes.None(),
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			subject, err := values.NewMappedValuer(
				tc.responseKey,
				tc.valueKey,
				tc.jsonPath,
				"",
				enums.RequestLocations.Body(),
				nil,
				tc.conv,
				logger.NewTestLogger(),
			)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request, err = http.NewRequest(http.MethodPost, "/test?test=asd&test=asd", strings.NewReader(tc.payload))
			require.NoError(t, err)
			value, err := subject.Generate(c)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.Equal(t, tc.expected, value)
			}
		})
	}
}

func TestURLMapper_Generate(t *testing.T) {
	t.Parallel()
	const servedURL = "http://domain.com/super/page/1234"
	testCases := []struct {
		name                string
		responsKey          string
		urlKey              string
		inputURL            string
		expected            any
		expectedCreationErr error
		expectedRunningErr  error
		conversion          enums.ConversionType
	}{
		{
			name:                "empty url",
			urlKey:              "test",
			inputURL:            "",
			conversion:          enums.ConversionTypes.None(),
			expectedCreationErr: values.ErrFailedLocatingElement,
		},
		{
			name:                "empty key",
			urlKey:              "",
			inputURL:            "/super/:test",
			conversion:          enums.ConversionTypes.None(),
			expectedCreationErr: values.ErrEmptyKey,
		},
		{
			name:       "correct key",
			urlKey:     "test",
			inputURL:   "/:test/page",
			expected:   "super",
			conversion: enums.ConversionTypes.None(),
		},
		{
			name:       "correct key with conversion to string",
			urlKey:     "test",
			inputURL:   "/super/page/:test",
			expected:   "1234",
			conversion: enums.ConversionTypes.Text(),
		},
		{
			name:       "correct key with conversion to int",
			urlKey:     "test",
			inputURL:   "/super/page/:test",
			expected:   float64(1234),
			conversion: enums.ConversionTypes.Number(),
		},
		{
			name:     "correct key and response key",
			urlKey:   "test",
			inputURL: "/super/page/:test",
			expected: map[string]any{
				"best-key": "1234",
			},
			responsKey: "best-key",
			conversion: enums.ConversionTypes.None(),
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			subject, err := values.NewMappedValuer(
				tc.responsKey,
				tc.urlKey,
				"",
				tc.inputURL,
				enums.RequestLocations.URL(),
				nil,
				tc.conversion,
				logger.NewTestLogger(),
			)
			if tc.expectedCreationErr != nil {
				require.ErrorIs(t, err, tc.expectedCreationErr)
				return
			}
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, servedURL, nil)
			actual, err := subject.Generate(c)
			if tc.expectedRunningErr != nil {
				require.ErrorIs(t, err, tc.expectedRunningErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestQueryMapper_Generate(t *testing.T) {
	t.Parallel()
	index := 1
	negativeIndex := -1
	const requestURL = "http://domain.com?test=1234&test2=asd&ids[]=16&ids[]=18"
	testCases := []struct {
		name          string
		responseKey   string
		queryKey      string
		conversion    enums.ConversionType
		index         *int
		expected      any
		expectedError error
	}{
		{
			name:          "key not present in request",
			queryKey:      "not-present",
			expectedError: values.ErrFailedLocatingElement,
		},
		{
			name:     "key with no conversion",
			queryKey: "test2",
			expected: "asd",
		},
		{
			name:          "key with conversion to int, but if fails",
			queryKey:      "test2",
			conversion:    enums.ConversionTypes.Number(),
			expectedError: values.ErrConversionFailed,
		},
		{
			name:       "key with conversion to int",
			queryKey:   "test",
			expected:   float64(1234),
			conversion: enums.ConversionTypes.Number(),
		},
		{
			name:        "key with conversion to string and response key",
			queryKey:    "test",
			responseKey: "best-key",
			expected: map[string]any{
				"best-key": "1234",
			},
			conversion: enums.ConversionTypes.Text(),
		},
		{
			name:       "key with conversion to string and an index",
			queryKey:   "ids[]",
			index:      &index,
			expected:   "18",
			conversion: enums.ConversionTypes.None(),
		},
		{
			name:          "index out of range",
			queryKey:      "ids[]",
			index:         &negativeIndex,
			expectedError: values.ErrFailedLocatingElement,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			subject, err := values.NewMappedValuer(
				tc.responseKey,
				tc.queryKey,
				"",
				"",
				enums.RequestLocations.Query(),
				tc.index,
				tc.conversion,
				logger.NewTestLogger(),
			)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(recorder)
			c.Request = httptest.NewRequest(http.MethodPost, requestURL, nil)
			actual, err := subject.Generate(c)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestQueryMapper_IsNill_Type(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		pm            *values.PayloadMapper
		expectedIsNil bool
		expectedType  enums.GenerationType
	}{
		{
			name:          "nil",
			pm:            nil,
			expectedIsNil: true,
		},
		{
			name:          "not nil",
			pm:            &values.PayloadMapper{},
			expectedIsNil: false,
			expectedType:  enums.GenerationTypes.SingleValue(),
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.expectedIsNil {
				require.True(t, tc.pm.IsNil())
			} else {
				require.False(t, tc.pm.IsNil())
				require.Equal(t, tc.expectedType, tc.pm.Type())
			}
		})
	}
}
