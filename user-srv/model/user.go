/*
 * @Author: weihua hu
 * @Date: 2024-11-24 15:21:20
 * @LastEditTime: 2024-11-26 23:26:05
 * @LastEditors: weihua hu
 * @Description:
 */

package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32     `gorm:"primeryKey"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type ExtraInfo struct {
}

type User struct {
	BaseModel
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	Mobile    string    `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Gender    string    `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role      int       `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
	ExtraInfo ExtraInfo `gorm:"column:extra_info;type:json"` // 新增的 extra_info 字段
}

func (User) TableName() string {
	return "users"
}

type UserAddress struct {
	BaseModel
	UserId    int32  `gorm:"index:idx_user_id;column:user_id;not null"` // 外键，关联用户表
	Address   string `gorm:"column:address;type:varchar(255);not null"` // 地址字段
	IsDefault bool   `gorm:"column:is_default;default:false"`           // 是否为默认地址
}

func (UserAddress) TableName() string {
	return "user_addresses"
}
