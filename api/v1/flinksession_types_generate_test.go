package v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildNodeSelector(t *testing.T) {
	maps := make(map[string]string)

	res := buildNodeSelector(maps)

	assert.Equal(t, "", res)

	maps["flink"] = "run"

	res = buildNodeSelector(maps)

	assert.Equal(t, "flink:run", res)

	maps["disk"] = "ssd"

	res = buildNodeSelector(maps)

	assert.Equal(t, "flink:run,disk:ssd", res)
}
