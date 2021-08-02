// +build integration

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
		assert.Equal(t, expected, actual)
	}

}
