/*
 * @Author: weihua hu
 * @Date: 2024-11-28 00:45:16
 * @LastEditTime: 2024-11-30 19:49:06
 * @LastEditors: weihua hu
 * @Description:
 */
package model

// 用户积分信息
type UserPoints struct {
	BaseModel
	UserId       int32 `gorm:"type:int;index"` // 用户ID
	Points       int32 `gorm:"type:int"`       // 可用积分
	FreezePoints int32 `gorm:"type:int"`       // 冻结积分
}

func (UserPoints) TableName() string {
	return "user_points"
}

// 积分变动记录
type PointsTransaction struct {
	BaseModel
	UserId    int32  `gorm:"type:int;index"`                    // 用户ID
	OrderSn   string `gorm:"type:varchar(200);column:order_sn"` // 订单号
	Change    int32  `gorm:"type:int"`                          // 积分变动数量（正数表示增加，负数表示减少）
	Status    int32  `gorm:"type:int"`                          // 状态：1 表示等待支付 2 表示支付成功 3 表示失败
	Timestamp int64  `gorm:"type:bigint"`                       // 变动时间戳
}

func (PointsTransaction) TableName() string {
	return "points_transactions"
}
