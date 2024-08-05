package tools

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func SaveToAFile(t *testing.T, content, filePath string) {
	f, err := os.Create(filePath)
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("failed to close file: %v", err)
		}
	}()
	require.NoError(t, err)

	_, err = f.WriteString(content)
	require.NoError(t, err)
}
