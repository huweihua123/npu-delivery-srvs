/*
 * @Author: weihua hu
 * @Date: 2024-11-24 19:31:23
 * @LastEditTime: 2024-11-24 19:31:23
 * @LastEditors: weihua hu
 * @Description:
 */

package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int" json:"id"` // 主键
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`      // 创建时间
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`   // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`                // 软删除时间
	IsDeleted bool           `json:"-"`                             // 是否删除
}
