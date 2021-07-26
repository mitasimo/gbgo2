package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {

	expected := []string{
		"../tests/f1",
		"../tests/f2",
		"../tests/f4",
		"../tests/sum/s1",
		"../tests/sum/s2",
		"../tests/sum/s3",
	}

	actual, err := iterateFilesInDirerctory("../tests", true)

	if assert.NoError(t, err) {
		// for _, path := range actual {
		// 	assert.Contains(t, expected, path)
		// }
		assert.Equal(t, expected, actual)
	}

}
