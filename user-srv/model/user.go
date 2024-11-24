/*
 * @Author: weihua hu
 * @Date: 2024-11-24 15:21:20
 * @LastEditTime: 2024-11-24 15:21:21
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

type User struct {
	BaseModel
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
	Email    string `gorm:"column:email"`
	Mobile   string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Gender   string `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role     int    `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
}
