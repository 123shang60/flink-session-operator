package pkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitContainer(t *testing.T) {
	res := SplitContainer("kafkaClient", ",", "Client")
	assert.Equal(t, false, res)
	res = SplitContainer("Client", ",", "Client")
	assert.Equal(t, true, res)
	res = SplitContainer("kafkaClient,Client", ",", "Client")
	assert.Equal(t, true, res)
}
