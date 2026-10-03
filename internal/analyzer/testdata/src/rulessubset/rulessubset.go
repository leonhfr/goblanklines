package rulessubset

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mixed(t *testing.T) {
	t.Helper() // want "missing newline after test helper calls"
	fmt.Println("setup")
	assert.Equal(t, 1, 1)
	assert.True(t, true)
	fmt.Println("done")
}
