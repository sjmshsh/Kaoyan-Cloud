package dao

import (
	"github.com/sjmshsh/grpc-gin-admin/project_user/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
)

var MysqlDB *gorm.DB

func InitMysql() {
	var builder strings.Builder
	s := []string{
		config.C.MysqlConfig.UserName,
		":",
		config.C.MysqlConfig.Password,
		"@tcp(",
		config.C.MysqlConfig.Addrs[0],
		")/",
		config.C.MysqlConfig.DbName,
		"?charset=utf8&parseTime=True&loc=Local",
	}
	for _, str := range s {
		builder.WriteString(str)
	}
	dsn := builder.String()
	log.Println(dsn)
	//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	MysqlDB = db
}
