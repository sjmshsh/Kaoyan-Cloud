package discovery


import (
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc/resolver"
	"strings"
)

type Server struct {
	Name    string `json:"name"`    // 名字为服务的名字(用来进行服务的发现)
	Addr    string `json:"addr"`    // 服务的地址(存储服务的地址)
	Version string `json:"version"` // 服务的版本(方便服务的版本迭代)
	Weight  int64  `json:"weight"`  // 服务的权重(后续用来降级熔断)
}

// BuildPrefix 定义服务名字前缀的函数
func BuildPrefix(server Server) string {
	if server.Version == "" {
		return fmt.Sprintf("/%s/", server.Name)
	}
	return fmt.Sprintf("/%s/%s/", server.Name, server.Version)
}

// BuildRegisterPath 定义注册的地址函数
func BuildRegisterPath(server Server) string {
	return fmt.Sprintf("%s%s", BuildPrefix(server), server.Addr)
}

// ParseValue 将值反序列化成一个注册Server服务
func ParseValue(value []byte) (Server, error) {
	server := Server{}
	if err := json.Unmarshal(value, &server); err != nil {
		return server, err
	}
	return server, nil
}

// SplitPath 分割路径，后续用作Server地址的更新
func SplitPath(path string) (Server, error) {
	server := Server{}
	str := strings.Split(path, "/")
	if len(str) == 0 {
		return server, errors.New("invalid path")
	}
	server.Addr = str[len(str)-1]
	return server, nil
}

// Exist 判断这个服务地址是否已经存在，防止服务访问冲突
func Exist(l []resolver.Address, addr resolver.Address) bool {
	for i := range l {
		// 如果找到了，说明是存在的，返回true
		if l[i].Addr == addr.Addr {
			return true
		}
	}
	return false
}

// Remove 移除服务
func Remove(s []resolver.Address, addr resolver.Address) ([]resolver.Address, bool) {
	for i := range s {
		// 这个删除服务的方式好原始啊，呵呵呵
		if s[i].Addr == addr.Addr {
			s[i] = s[len(s) - 1]
			return s[:len(s) - 1], true
		}
	}
	return nil, false
}