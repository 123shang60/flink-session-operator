package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommand(t *testing.T) {
	var c Command
	c.FieldConfig("kubernetes.cluster-id", "flink")
	res := c.Build()

	assert.Equal(t, `-Dkubernetes.cluster-id=flink `, res)
}

func TestCommandNil(t *testing.T) {
	var c Command
	c.FieldConfig("kubernetes.cluster-id", "")
	res := c.Build()

	assert.Equal(t, ``, res)
}
