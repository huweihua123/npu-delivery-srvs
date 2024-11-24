/*
 * @Author: weihua hu
 * @Date: 2024-11-24 19:31:48
 * @LastEditTime: 2024-11-24 19:42:11
 * @LastEditors: weihua hu
 * @Description:
 */
package model

// 食堂模型
type Canteen struct {
	BaseModel
	Name string `gorm:"type:varchar(255);not null" json:"name"` // 食堂名称
}

// 商家模型
type Merchant struct {
	BaseModel
	Name      string `gorm:"type:varchar(255);not null" json:"name"` // 商家名称
	CanteenID int32  `gorm:"type:int;not null" json:"canteen_id"`    // 所属食堂ID（逻辑关联）
}

// 商品种类模型
type Category struct {
	BaseModel
	Name       string `gorm:"type:varchar(255);not null" json:"name"` // 种类名称
	MerchantID int32  `gorm:"type:int;not null" json:"merchant_id"`   // 所属商家ID（逻辑关联）
}

// 商品模型
type Product struct {
	BaseModel
	Name        string  `gorm:"type:varchar(255);not null" json:"name"`   // 商品名称
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"` // 商品价格
	CategoryID  int32   `gorm:"type:int;not null" json:"category_id"`     // 所属种类ID（逻辑关联）
	MerchantID  int32   `gorm:"type:int;not null" json:"merchant_id"`     // 所属商家ID（逻辑关联）
	Sales       int32   `gorm:"type:int;default:0" json:"sales"`          // 销量
	Description string  `gorm:"type:varchar(500)" json:"description"`     // 商品描述（可选）
}
