/*
 * @Author: weihua hu
 * @Date: 2024-11-28 01:16:01
 * @LastEditTime: 2024-11-29 23:04:27
 * @LastEditors: weihua hu
 * @Description:
 */
package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"npu-delivery-srvs/point-srv/model"
)

func main() {
	dsn := "root:123456hwh@tcp(127.0.0.1:3306)/npu_delivery_point_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	// 全局模式
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	_ = db.AutoMigrate(&model.UserPoints{})
	_ = db.AutoMigrate(&model.PointsTransaction{})
}
