package discovery

import (
	"context"
	"encoding/json"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Register 定义服务实例的实例，用来存储全部的实例信息，并且维持各个服务之间的执行，防止宕机
type Register struct {
	EtcdAddrs   []string                                // etcd的地址，例如http://etcd-1.etcd-headless.devops.svc.cluster.local:2379,服务地址
	DialTimeout int                                     // DialTimeout is the timeout for failing to establish a connection.
	closeCh     chan struct{}                           // 是否关闭
	leasesID    clientv3.LeaseID                        // 租约
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 心跳检测
	srvInfo     Server                                  // 服务
	srvTTL      int64                                   // 服务的TTL
	cli         *clientv3.Client                        // 客户端
	logger      *zap.Logger                             // 日志
}

// NewRegister 创建一个服务对象
func NewRegister(etcdAddrs []string, logger *zap.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Register 注册服务到etcd中
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		// 判断服务地址的正确性
		return nil, errors.New("invalid ip address")
	}
	// 对服务进行注册
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	// 服务信息的注册
	r.srvInfo = srvInfo
	// 服务的存活时间
	r.srvTTL = ttl

	if err = r.register(); err != nil {
		return nil, err
	}
	// 初始化一个切片来判断这个服务连接是否关闭
	r.closeCh = make(chan struct{})
	// 异步进行心跳检测
	go r.keepAlive()

	return r.closeCh, nil
}

// Stop 关闭服务连接
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(
		context.Background(),
		BuildRegisterPath(r.srvInfo),
	)
	return err
}

// 存活检测
func (r *Register) keepAlive() {
	// 超时器的用法，总不能让这个东西一直阻塞着
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh: // 是否存在这个服务
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed, error")
			}
			// 撤销租约
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed, error:")
			}
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed, error:")
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed, error: ")
				}
			}
		}
	}
}

func (r *Register) register() error {
	// 设置超时时间，访问etcd有超时控制
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	// 注册一个新的租约
	leaseResp, err := r.cli.Grant(ctx, r.srvTTL)
	if err != nil {
		return err
	}

	// 赋值租约的ID
	r.leasesID = leaseResp.ID

	// 对这个cli进行心跳检测
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), r.leasesID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	// 将服务写到ETCD中
	_, err = r.cli.Put(
		context.Background(),
		BuildRegisterPath(r.srvInfo),
		string(data),
		clientv3.WithLease(r.leasesID),
	)

	return err
}

// GetServerInfo 获取服务注册的信息
func (r *Register) GetServerInfo() (Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegisterPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}

	server := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &server); err != nil {
			return server, err
		}
	}
	return server, err
}
