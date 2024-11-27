/*
 * @Author: weihua hu
 * @Date: 2024-11-27 20:41:10
 * @LastEditTime: 2024-11-27 21:25:59
 * @LastEditors: weihua hu
 * @Description:
 */

package initialize

import (
	"encoding/json"
	"fmt"
	"npu-delivery-srvs/user-srv/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	v := viper.New() //创建viper实例

	// 读取配置文件
	v.SetConfigName("config-debug") // 配置文件名称(不需要带后缀)
	v.SetConfigType("yaml")         // 配置文件类型
	v.AddConfigPath("./user-srv")   // 配置文件路径(这里使用相对路径)

	if err := v.ReadInConfig(); err != nil { // 处理读取配置文件的错误
		panic(err)
	}

	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}

	zap.S().Infof("配置信息: &v", global.NacosConfig)

	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId: global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以场景多个client,它们有不同的NamespaceId
		TimeoutMs:   5000,
		LogDir:      "tmp/nacos/log",
		CacheDir:    "tmp/nacos/cache",
		LogLevel:    "debug",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group:  global.NacosConfig.Group})

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	}

	zap.S().Infof("读取nacos配置成功： %s", content)
	fmt.Printf("%+v\n", global.ServerConfig)
}
