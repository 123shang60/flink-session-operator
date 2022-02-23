package pkg

import (
	"github.com/go-zookeeper/zk"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	ZookeeperQuorum string = "127.0.0.1:2181"
	TestPath        string = "/test"
)

func TestNormalDelete(t *testing.T) {
	conn, err := AutoConnZk(ZookeeperQuorum)
	defer conn.Close()
	assert.Nil(t, err)

	var data = []byte("test value")
	acls := zk.WorldACL(zk.PermAll)
	s, err := conn.Create(TestPath, data, 0, acls)
	assert.Nil(t, err)
	assert.Equal(t, TestPath, s)

	_, err = conn.Create(TestPath+"/test", data, 0, acls)
	assert.Nil(t, err)

	_, err = conn.Create(TestPath+"/test/test", data, 0, acls)
	assert.Nil(t, err)

	err = conn.AutoDelete(TestPath)
	assert.Nil(t, err)

	_, _, err = conn.Get(TestPath)
	assert.NotNil(t, err)
}

func TestErrorDelete(t *testing.T) {
	conn, err := AutoConnZk(ZookeeperQuorum)
	defer conn.Close()
	assert.Nil(t, err)

	err = conn.AutoDelete(TestPath)
	assert.NotNil(t, err)

	_, _, err = conn.Get(TestPath)
	assert.NotNil(t, err)
}
