package parser_test

import (
	"path"
	"strings"
	"testing"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/parser"
	"github.com/vimek-go/server-faker/internal/pkg/parser/internal/mocks"
	"github.com/vimek-go/server-faker/internal/pkg/tools"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLoader_LoadConfig(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	testValidJSON := `
	{
		"endpoints": [
			{
				"url": "/test/test",
				"method": "GET",
				"response": {
					"status": 200,
					"type": "static",
					"file": "test-test.json",
					"format": "json"
				}
			}
		]
	}`
	testCases := []struct {
		name          string
		jsonConfig    string
		factory       func(string) *mocks.FactoryMock
		expected      []api.Handler
		expectedError error
	}{
		{
			name: "multiple validation errors",
			factory: func(string) *mocks.FactoryMock {
				return mocks.NewFactoryMock(t)
			},
			jsonConfig: `
			{
				"endpoints": [
					{
						"url": "/test",
						"response": {}
					}
				]
			}`,
			expectedError: parser.ErrValidation,
		},
		{
			name: "no validation errors, handler creation error",
			factory: func(dir string) *mocks.FactoryMock {
				factory := mocks.NewFactoryMock(t)
				factory.On("CreateEndpoint", mock.AnythingOfType("dto.Endpoint"), dir).Return(nil, parser.ErrNotHandled)
				return factory
			},
			jsonConfig:    testValidJSON,
			expectedError: parser.ErrNotHandled,
		},
		{
			name: "no validation errors, handler creation success",
			factory: func(dir string) *mocks.FactoryMock {
				factory := mocks.NewFactoryMock(t)
				factory.On("CreateEndpoint", mock.AnythingOfType("dto.Endpoint"), dir).Return(&mocks.HandlerMock{}, nil)
				return factory
			},
			jsonConfig: testValidJSON,
			expected:   []api.Handler{&mocks.HandlerMock{}},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			file := path.Join(dir, strings.ReplaceAll(tc.name, " ", "_")+".json")
			tools.SaveToAFile(t, tc.jsonConfig, file)
			loader := parser.NewLoader(tc.factory(dir), logger.NewTestLogger())
			handlers, err := loader.LoadConfig(file)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				require.NotNil(t, handlers)
			}
		})
	}
}
