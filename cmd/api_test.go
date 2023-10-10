// BEGIN: 2d7f8a6c7b3d
package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApi(t *testing.T) {
	cmd := Api()

	assert.NotNil(t, cmd)
	assert.Equal(t, "api", cmd.Use)
	assert.Equal(t, "Run the api server for krcrdr", cmd.Short)
}

// END: 2d7f8a6c7b3d
