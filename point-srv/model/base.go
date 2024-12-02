/*
 * @Author: weihua hu
 * @Date: 2024-11-28 00:45:45
 * @LastEditTime: 2024-11-28 00:45:46
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
