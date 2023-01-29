package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/sjmshsh/grpc-gin-admin/project_common/logs"
	"github.com/spf13/viper"
	"log"
	"os"
)

var C = InitConfig()

type Config struct {
	viper        *viper.Viper
	SC           *ServerConfig
	GC           *GrpcConfig
	EtcdConfig   *EtcdConfig
	RabbitConfig *RabbitConfig
	MysqlConfig  *MysqlConfig
}

type ServerConfig struct {
	Name string
	Addr string
}

type MysqlConfig struct {
	Addrs    []string
	UserName string
	Password string
	DbName   string
}

type RabbitConfig struct {
	Addr     []string
	Password string
	UserName string
	Virtual  string
}

type GrpcConfig struct {
	Name    string
	Addr    string
	Version string
	Weight  int64
}

type EtcdConfig struct {
	Addr []string
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	ec.Addr = addrs
	c.EtcdConfig = ec
}

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.weight")
	c.GC = gc
}

func (c *Config) ReadRedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.DB"),
	}
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("app")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	err := conf.viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	conf.ReadServerConfig()
	conf.ReadRedisConfig()
	conf.InitZapLog()
	conf.ReadGrpcConfig()
	conf.ReadEtcdConfig()
	conf.ReadRabbitConfig()
	conf.ReadMysqlConfig()
	return conf
}

func (c *Config) InitZapLog() {
	// 从配置文件中读取日志配置，初始化项目
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.debugFileName"),
		InfoFileName:  c.viper.GetString("zap.infoFileName"),
		WarnFileName:  c.viper.GetString("zap.warnFileName"),
		MaxSize:       c.viper.GetInt("maxSize"),
		MaxAge:        c.viper.GetInt("maxAge"),
		MaxBackups:    c.viper.GetInt("maxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) ReadMysqlConfig() {
	ms := &MysqlConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("mysql.addrs", &addrs)
	if err != nil {
		log.Println(err)
	}
	ms.Addrs = addrs
	var username string
	err = c.viper.UnmarshalKey("mysql.user", &username)
	if err != nil {
		log.Println(err)
	}
	ms.UserName = username
	var password string
	err = c.viper.UnmarshalKey("mysql.password", &password)
	if err != nil {
		log.Println(err)
	}
	ms.Password = password
	var dbname string
	err = c.viper.UnmarshalKey("mysql.dbname", &dbname)
	if err != nil {
		log.Println(err)
	}
	ms.DbName = dbname
	c.MysqlConfig = ms
}

func (c *Config) ReadRabbitConfig() {
	rb := &RabbitConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("rabbit.addrs", &addrs)
	if err != nil {
		log.Println(err)
	}
	rb.Addr = addrs
	var username string
	err = c.viper.UnmarshalKey("rabbit.username", &username)
	if err != nil {
		log.Println(err)
	}
	rb.UserName = username
	var password string
	err = c.viper.UnmarshalKey("rabbit.password", &password)
	if err != nil {
		log.Println(err)
	}
	rb.Password = password
	var virtual string
	err = c.viper.UnmarshalKey("rabbit.virtual", &virtual)
	if err != nil {
		log.Println(err)
	}
	rb.Virtual = virtual
	c.RabbitConfig = rb
}
