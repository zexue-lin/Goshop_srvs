package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"goshop_srvs/user_srv/model"
	_ "goshop_srvs/user_srv/model"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

/*
	密文：密文不可反解
	1.对称加密
	2.非对称加密
	3.md5 信息摘要算法
	密码如果不可反解，用户找回密码的时候直接给个链接，修改密码
	md5 不可反解（不可逆），就算加密的字符差异很小，但加密后差异很大。但是简单的密码可以通过暴力破解反解。彩虹表
	用盐值加密可以
*/

func genMd5(code string) string {
	// 生成md5实例
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code) // 将字符串写入到Md5哈希对象中

	return hex.EncodeToString(Md5.Sum(nil))
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/goshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用结构体的单数形式作为表名，而不是默认的复数形式。
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("admin123", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd) // 密码字符串三个部分，算法，盐值，真正的密码密码

	// fmt.Println(len(newPassword)) // 要确保长度不能超过100，否则保存到数据库会被截断

	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("tom%d", i),
			Mobile:   fmt.Sprintf("1576543211%d", i),
			Password: newPassword,
		}
		db.Save(&user)
	}

	//_ = db.AutoMigrate(&model.User{})

	//fmt.Println("newPassword=", newPassword)
	//
	//passwordInfo := strings.Split(newPassword, "$")
	//fmt.Println(passwordInfo)
	//check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options) // 第0个是空格，第2个是salt
	//fmt.Println(check)                                                                      // true
}
