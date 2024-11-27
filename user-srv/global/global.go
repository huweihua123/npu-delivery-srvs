/*
 * @Author: weihua hu
 * @Date: 2024-11-25 00:25:41
 * @LastEditTime: 2024-11-27 20:33:36
 * @LastEditors: weihua hu
 * @Description:
 */

package global

import (
	"npu-delivery-srvs/user-srv/config"

	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
)
