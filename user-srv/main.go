/*
 * @Author: weihua hu
 * @Date: 2024-11-25 00:30:18
 * @LastEditTime: 2024-11-27 21:44:13
 * @LastEditors: weihua hu
 * @Description:
 */

package main

import (
	"flag"
	"fmt"
	"net"
	"npu-delivery-srvs/user-srv/global"
	"npu-delivery-srvs/user-srv/handler"
	"npu-delivery-srvs/user-srv/initialize"
	"npu-delivery-srvs/user-srv/proto"
	"npu-delivery-srvs/user-srv/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	*Port = 53415

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 初始化 Consul 客户端配置
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	// 创建 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		zap.S().Fatal("Failed to create Consul client:", err)
	}

	// 设置健康检查
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.31.192:%d", *Port),
		Timeout:                        "10s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "1m",
	}

	// 创建注册到 Consul 的服务
	serviceId := uuid.NewV4().String()
	registration := &api.AgentServiceRegistration{
		Name:    global.ServerConfig.Name,
		ID:      serviceId,
		Port:    *Port,
		Tags:    []string{"user", "srv"},
		Address: "192.168.31.192",
		Check:   check,
	}

	zap.S().Infof("Registering service to Consul with ID: %s on port %d", registration.ID, *Port)

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Fatal("Failed to register service with Consul:", err)
	}

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err = client.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")

}
