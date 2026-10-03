package combined

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func blockRuleExample() {
	x := 5
	if x > 0 {
		fmt.Println("positive")
	} // want "missing newline after block statement"
	for i := 0; i < x; i++ {
		fmt.Println(i)
	}
}
func TestCombined(t *testing.T) { // want "missing newline between function declarations"
	t.Helper()
	t.Parallel() // want "missing newline after test helper calls"
	x := compute()
	assert.Equal(t, 1, x)
	assert.True(t, x > 0) // want "missing newline after testify assertion group"
	t.Run("sub", func(t *testing.T) {}) // want "missing newline before t.Run calls" "missing newline after t.Run calls"
	fmt.Println("done")
}

func compute() int {
	return 1
}
