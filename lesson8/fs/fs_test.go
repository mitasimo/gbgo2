package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArgsArrTo(t *testing.T) {
	actual := StringArgsArrTo("123", "456")
	expected := []string{"123", "456"}
	assert.Equal(t, expected, actual)
}
