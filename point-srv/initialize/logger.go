/*
 * @Author: weihua hu
 * @Date: 2024-11-30 00:42:44
 * @LastEditTime: 2024-11-30 00:42:45
 * @LastEditors: weihua hu
 * @Description:
 */

package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
