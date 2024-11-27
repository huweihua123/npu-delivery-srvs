/*
 * @Author: weihua hu
 * @Date: 2024-11-27 12:52:18
 * @LastEditTime: 2024-11-27 19:07:23
 * @LastEditors: weihua hu
 * @Description:
 */

package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
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
	Age int `json:"age"`
}

/**
 * @Author: weihua hu
 * @description: 实现 Scanner 接口（从数据库读取）
 * @param {interface{}} value
 * @return {error}
 */
func (ei *ExtraInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	// 确保 value 是 []byte 类型（字节切片）
	switch v := value.(type) {
	case []byte:
		// 如果是 []byte 类型，直接将其转换为字符串进行反序列化
		return json.Unmarshal(v, ei)
	case string:
		// 如果是 string 类型，将其转换为 []byte 进行反序列化
		return json.Unmarshal([]byte(v), ei)
	default:
		err := fmt.Errorf("expected string or []byte, got %T", value)
		zap.S().Errorf("Scan failed: %v", err)
		return err
	}

}

/**
 * @Author: weihua hu
 * @description: 实现 Valuer 接口（将结构体写入数据库）
 * @return {driver.Value, error}
 */
func (ei ExtraInfo) Value() (driver.Value, error) {
	return json.Marshal(ei)
}

type User struct {
	BaseModel
	Username  string    `gorm:"column:username"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	Mobile    string    `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Gender    string    `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女, male表示男'"`
	Role      int       `gorm:"column:role;default:1;type:int comment '1表示普通用户, 2表示管理员'"`
	ExtraInfo ExtraInfo `gorm:"column:extra_info;type:json"`
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
