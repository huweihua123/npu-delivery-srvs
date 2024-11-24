/*
 * @Author: weihua hu
 * @Date: 2024-11-24 15:21:44
 * @LastEditTime: 2024-11-24 15:24:52
 * @LastEditors: weihua hu
 * @Description:
 */

package main

import (
	"npu-delivery-srvs/user-srv/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456hwh@tcp(127.0.0.1:3306)/npu_delivery_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{})
}
