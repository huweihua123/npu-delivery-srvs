/*
 * @Author: weihua hu
 * @Date: 2024-11-25 00:27:05
 * @LastEditTime: 2024-11-25 00:39:04
 * @LastEditors: weihua hu
 * @Description:
 */

package initialize

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"npu-delivery-srvs/user-srv/global"
)

func InitDB() {
	// 连接数据库
	dsn := "root:123456hwh@tcp(127.0.0.1:3306)/npu_delivery_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
}
