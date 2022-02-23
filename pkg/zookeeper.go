package pkg

import (
	"github.com/go-zookeeper/zk"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

type ZooKeeper struct {
	*zk.Conn
}

func AutoConnZk(quorum string) (*ZooKeeper, error) {
	host := strings.Split(quorum, ",")
	connect, _, err := zk.Connect(host, time.Second*30)
	if err != nil {
		klog.Error("zk 连接失败！", err)
		return nil, err
	}
	return &ZooKeeper{
		connect,
	}, nil
}

// TODO: 递归删除
func (zoo *ZooKeeper) AutoDelete(path string) error {
	klog.Infof("准备删除zk 路径: %s", path)
	//children, status, err := zoo.Children(path)
	//if err != nil {
	//	klog.Error("zk 子路径检查失败！,停止zk删除流程", err)
	//}
	//
	//for _, child := range children {
	//	err = zoo.Delete(child, status.Version)
	//	if err != nil {
	//		klog.Error(err)
	//	}
	//}

	_, stat, err := zoo.Get(path)
	if err != nil {
		klog.Error("zk 路径检查失败！,停止zk删除流程", err)
		return err
	}
	if err = zoo.Delete(path, stat.Version); err != nil {
		klog.Error("zk 路径删除失败，跳过删除逻辑", err)
		return err
	}

	klog.Infof("zk 路径删除成功！%s", path)
	return nil
}
