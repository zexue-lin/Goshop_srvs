package global

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

// import的时候，init()方法会自动执行，不用显式的调用它
func init() {
	dsn := "root:root@tcp(127.0.0.1:3306)/goshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"

	// 【日志配置】设置全局的logger，这个logger在执行每个sql语句的时候都会打印每一行sql
	// sql才是最重要的，std -> standard 标准的
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log 级别
			//IgnoreRecordNotFoundError: true,        // 查询没有找到记录时，不会记录错误日志。
			//ParameterizedQueries:      true,        // 不在日志中显示查询参数值，只显示占位符
			Colorful: true, // 彩色打印
		},
	)

	// 全局模式
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用结构体的单数形式作为表名，而不是默认的复数形式。
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

}
