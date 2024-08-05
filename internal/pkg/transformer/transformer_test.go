package transformer_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/tools"
	"github.com/vimek-go/server-faker/internal/pkg/transformer"

	"github.com/stretchr/testify/require"
)

func TestTransformer_Transform(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name          string
		fileContent   string
		url           string
		responseType  string
		expectedError error
		expected      string
	}{
		{
			name:          "invalid response type",
			fileContent:   `{}`,
			responseType:  "invalid",
			expectedError: transformer.ErrInvalidResponseType,
		},
		{
			name:         "valid static response type",
			fileContent:  `{"test": "value"}`,
			url:          "/test",
			responseType: "static",
			expected: `
			{
				"endpoints": [
					{
						"url": "/test",
						"method": "GET",
						"response": {
							"status": 200,
							"type": "dynamic",
							"headers": null,
							"file": "",
							"static": null,
							"object": [
								{
									"key": "test",
									"static": {
										"value": "value"
									}
								}
							],
							"format": "json"
						},
						"proxy": null
					}
				]
			}
			`,
		},
		{
			name:         "valid dynamic array with objects response type",
			fileContent:  `[{"string": "value", "integer": 1, "boolean": true, "float": 1.1}]`,
			url:          "/array-test",
			responseType: "dynamic",
			expected: `
			{
			"endpoints": [
				{
				"url": "/array-test",
				"method": "GET",
				"response": {
					"status": 200,
					"type": "dynamic",
					"headers": null,
					"file": "",
					"static": null,
					"object": [
					{
						"array": {
						"min": 3,
						"max": 3,
						"element": [
							{
								"key": "string",
								"random": {
									"type": "string-all",
									"min": 1,
									"max": 100
								}
							},
							{
								"key": "integer",
								"random": {
									"type": "integer",
									"min": 1,
									"max": 100
								}
							},
							{
								"key": "boolean",
								"random": {
									"type": "boolean",
									"min": 0,
									"max": 0
								}
							},
							{
								"key": "float",
								"random": {
									"type": "float",
									"min": 1,
									"max": 100
								}
							}
						]
						}
					}
				],
				"format": "json"
			},
			"proxy": null
			}
		]
		}
`,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tr := transformer.New()
			dir := t.TempDir()
			file := path.Join(dir, strings.ReplaceAll(tc.name, " ", "_")+".json")
			tools.SaveToAFile(t, tc.fileContent, file)
			actual, err := tr.Transform(file, tc.url, tc.responseType)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				fmt.Println(actual)
				dst := &bytes.Buffer{}
				err := json.Compact(dst, []byte(tc.expected))
				require.NoError(t, err)
				require.JSONEq(t, dst.String(), actual)
			}
		})
	}
}
