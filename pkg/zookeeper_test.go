package pkg

import (
	"encoding/base64"
	"github.com/123shang60/zk"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	ZookeeperQuorum string = "kerberos.example.com:2181"
	TestPath        string = "/test"
	TestInvalidPath string = "/test-1"

	Krb5Conf string = `# Configuration snippets may be placed in this directory as well
#includedir /etc/krb5.conf.d/

[logging]
 default = FILE:/var/log/krb5libs.log
 kdc = FILE:/var/log/krb5kdc.log
 admin_server = FILE:/var/log/kadmind.log

[libdefaults]
 dns_lookup_realm = false
 ticket_lifetime = 24h
 renew_lifetime = 7d
 forwardable = true
 rdns = false
 default_realm = EXAMPLE.COM
 #default_ccache_name = KEYRING:persistent:%{uid}

[realms]
 EXAMPLE.COM = {
  kdc = kerberos.example.com:88
  admin_server = kerberos.example.com:749
  default_domain = example.com
 }

[domain_realm]
 .example.com = EXAMPLE.COM
 example.com = EXAMPLE.COM`
	KeyTab string = `BQIAAABhAAIAC0VYQU1QTEUuQ09NAAl6b29rZWVwZXIAFGtlcmJlcm9zLmV4YW1wbGUuY29tAAAA
AWKONkQCABIAIO/D3KkXb3uP42EvWsWOwZoV7DqN6GinushTEMgCzsD7AAAAAgAAAFEAAgALRVhB
TVBMRS5DT00ACXpvb2tlZXBlcgAUa2VyYmVyb3MuZXhhbXBsZS5jb20AAAABYo42RAIAEQAQmfeP
cJgUOg/L1GIzgn1VuwAAAAI=`
	Pan string = `zookeeper/kerberos.example.com@EXAMPLE.COM`
)

func TestNormalDelete(t *testing.T) {
	key, _ := base64.StdEncoding.DecodeString(KeyTab)
	conn, err := AutoConnZk(ZookeeperQuorum, &KerberosConfig{
		Keytab:       key,
		Krb5:         Krb5Conf,
		PrincipalStr: Pan,
	})
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
	key, _ := base64.StdEncoding.DecodeString(KeyTab)
	conn, err := AutoConnZk(ZookeeperQuorum, &KerberosConfig{
		Keytab:       key,
		Krb5:         Krb5Conf,
		PrincipalStr: Pan,
	})
	defer conn.Close()
	assert.Nil(t, err)

	err = conn.AutoDelete(TestPath)
	assert.NotNil(t, err)

	_, _, err = conn.Get(TestPath)
	assert.NotNil(t, err)
}

func TestInvalidPathDelete(t *testing.T) {
	key, _ := base64.StdEncoding.DecodeString(KeyTab)
	conn, err := AutoConnZk(ZookeeperQuorum, &KerberosConfig{
		Keytab:       key,
		Krb5:         Krb5Conf,
		PrincipalStr: Pan,
	})
	defer conn.Close()
	assert.Nil(t, err)

	var data = []byte("test value")
	acls := zk.WorldACL(zk.PermAll)
	s, err := conn.Create(TestInvalidPath, data, 0, acls)
	assert.Nil(t, err)
	assert.Equal(t, TestInvalidPath, s)

	_, err = conn.Create(TestInvalidPath+"/test", data, 0, acls)
	assert.Nil(t, err)

	_, err = conn.Create(TestInvalidPath+"/test/test", data, 0, acls)
	assert.Nil(t, err)

	err = conn.AutoDelete(TestInvalidPath + "//")
	assert.Nil(t, err)

	_, _, err = conn.Get(TestInvalidPath)
	assert.NotNil(t, err)
}
