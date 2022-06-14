package pkg

import (
	"github.com/123shang60/zk"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

type ZooKeeper struct {
	*zk.Conn
}

type KerberosConfig struct {
	Keytab       []byte
	Krb5         string
	PrincipalStr string
}

func AutoConnZk(quorum string, conf *KerberosConfig) (*ZooKeeper, error) {
	host := strings.Split(quorum, ",")
	connect, _, err := func() (*zk.Conn, <-chan zk.Event, error) {
		if conf == nil {
			return zk.Connect(host, 30*time.Second)
		} else {
			return zk.Connect(host, 30*time.Second, zk.WithSASLConfig(&zk.SASLConfig{
				Enable: true,
				KerberosConfig: &zk.KerberosConfig{
					Keytab:       conf.Keytab,
					Krb5:         conf.Krb5,
					PrincipalStr: conf.PrincipalStr,
					ServiceName:  "zookeeper",
				},
			}))
		}
	}()
	if err != nil {
		klog.Error("zk 连接失败！", err)
		return nil, err
	}
	return &ZooKeeper{
		connect,
	}, nil
}

// AutoDelete 删除指定路径下全部数据
func (zoo *ZooKeeper) AutoDelete(path string) error {
	path = strings.TrimRight(path, "/")
	klog.Infof("准备删除zk 路径: %s", path)

	_, _, err := zoo.Get(path)
	if err != nil {
		klog.Error("zk 路径检查失败！,停止zk删除流程", err)
		return err
	}
	if err := zoo.deletePath(path); err != nil {
		klog.Error("zk 递归删除失败!", err)
		return err
	}

	klog.Infof("zk 路径删除成功！%s", path)
	return nil
}

func (zoo *ZooKeeper) deletePath(path string) error {
	childs, _, err := zoo.Children(path)
	if err != nil {
		return err
	}
	for _, child := range childs {
		err := zoo.deletePath(path + "/" + child)
		if err != nil {
			return err
		}
	}

	_, stat, err := zoo.Get(path)
	if err != nil {
		return err
	}

	err = zoo.Delete(path, stat.Version)
	return err
}
