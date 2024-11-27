/*
 * @Author: weihua hu
 * @Date: 2024-11-25 00:30:18
 * @LastEditTime: 2024-11-27 16:32:53
 * @LastEditors: weihua hu
 * @Description:
 */

package main

import (
	"flag"
	"fmt"
	"net"
	"npu-delivery-srvs/user-srv/handler"
	"npu-delivery-srvs/user-srv/initialize"
	"npu-delivery-srvs/user-srv/proto"
	"npu-delivery-srvs/user-srv/utils"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	*Port = 53415

	initialize.InitDB()
	initialize.InitLogger()

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// if err = client.Agent().ServiceDeregister(serviceID); err != nil {
	// 	zap.S().Info("注销失败")
	// }
	zap.S().Info("注销成功")

}
