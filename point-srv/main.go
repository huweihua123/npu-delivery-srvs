/*
 * @Author: weihua hu
 * @Date: 2024-11-29 23:00:16
 * @LastEditTime: 2024-11-30 19:32:26
 * @LastEditors: weihua hu
 * @Description:
 */

package main

import (
	"flag"
	"fmt"
	"net"
	"npu-delivery-srvs/point-srv/handler"
	"npu-delivery-srvs/point-srv/initialize"
	"npu-delivery-srvs/point-srv/proto"
	"npu-delivery-srvs/point-srv/utils"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50053, "端口号")

	initialize.InitLogger()
	initialize.InitDB()

	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	flag.Parse()

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterPointsServer(server, &handler.PointsServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}
	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 创建消费者实例
	c, err := consumer.NewPushConsumer(
		consumer.WithGroupName("points-group"),              // 消费者组，确保只有一个实例消费每条消息
		consumer.WithNameServer([]string{"127.0.0.1:9876"}), // RocketMQ NameServer 地址
	)
	if err != nil {
		zap.S().Fatalf("Error creating consumer: %v", err)
	}
	rlog.SetLogLevel("warn")

	// 订阅消息，订阅主题 "order_reback"
	err = c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoRebackPoints)
	if err != nil {
		zap.S().Fatalf("Error subscribing to topic: %v", err)
	}

	// 订阅消息，订阅主题 "order_paid"
	err = c.Subscribe("order_paid", consumer.MessageSelector{}, handler.AutoUnfreezePoints)
	if err != nil {
		zap.S().Fatalf("Error subscribing to topic: %v", err)
	}

	// 启动消费者
	err = c.Start()
	if err != nil {
		zap.S().Fatalf("Error starting consumer: %v", err)
	}

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	// if err = register_client.DeRegister(serviceId); err != nil {
	// 	zap.S().Info("注销失败:", err.Error())
	// } else {
	// 	zap.S().Info("注销成功:")
	// }
	zap.S().Info("注销成功:")

}
