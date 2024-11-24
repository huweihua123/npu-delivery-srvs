/*
 * @Author: weihua hu
 * @Date: 2024-11-25 01:00:57
 * @LastEditTime: 2024-11-25 01:00:58
 * @LastEditors: weihua hu
 * @Description:
 */
package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
