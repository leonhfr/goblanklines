package testify

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func groupMissingBlankAfter(t *testing.T) {
	x := compute()
	assert.Equal(t, 1, x)
	assert.True(t, x > 0) // want "missing newline after testify assertion group"
	fmt.Println("done")
}

func groupWithBlankAfter(t *testing.T) {
	x := compute()
	assert.Equal(t, 1, x)
	assert.True(t, x > 0)

	fmt.Println("done")
}

func mixedAssertRequireOneGroup(t *testing.T) {
	x := compute()
	assert.Equal(t, 1, x)
	require.NoError(t, nil)
	assert.True(t, x > 0) // want "missing newline after testify assertion group"
	fmt.Println("done")
}

func groupFirstInBlockStillNeedsAfter(t *testing.T) {
	if true {
		assert.Equal(t, 1, 1)
		assert.True(t, true) // want "missing newline after testify assertion group"
		fmt.Println("inner")
	}
}

func groupLastInBlockNoAfterNeeded(t *testing.T) {
	fmt.Println("setup")
	assert.Equal(t, 1, 1)
	assert.True(t, true)
}

func groupViaAssertionsObject(t *testing.T) {
	a := assert.New(t)
	x := compute()
	a.Equal(1, x)
	a.True(x > 0) // want "missing newline after testify assertion group"
	fmt.Println("done")
}

func requireAssertionsObject(t *testing.T) {
	r := require.New(t)
	x := compute()
	r.Equal(1, x)
	r.True(x > 0) // want "missing newline after testify assertion group"
	fmt.Println("done")
}

func conditionalAssertNotGrouped(t *testing.T) {
	x := compute()
	if !assert.Equal(t, 1, x) {
		return
	}
	fmt.Println("done")
}

func compute() int {
	return 1
}
