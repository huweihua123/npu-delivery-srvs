/*
 * @Author: weihua hu
 * @Date: 2024-11-28 15:30:15
 * @LastEditTime: 2024-12-01 01:06:02
 * @LastEditors: weihua hu
 * @Description:
 */
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"npu-delivery-srvs/point-srv/global"
	"npu-delivery-srvs/point-srv/model"
	"npu-delivery-srvs/point-srv/proto"
)

func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		fmt.Printf("Received message: %s\n", string(msg.Body))
	}
	return consumer.ConsumeSuccess, nil
}

type PointsServer struct {
	proto.UnimplementedPointsServer
}

/**
 * @Author: weihua hu
 * @description: 添加并冻结积分
 * @param {context.Context} ctx
 * @param {*proto.AddAndFreezePointsRequest} req
 * @return {*}
 */
func (s *PointsServer) AddAndFreezePoints(ctx context.Context, req *proto.AddAndFreezePointsRequest) (*proto.AddAndFreezePointsResponse, error) {
	// 检查用户是否存在于积分表中
	var userPoints model.UserPoints
	if err := global.DB.Where(&model.UserPoints{UserId: req.UserId}).First(&userPoints).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，插入一条默认积分为 0 的记录
			userPoints = model.UserPoints{
				UserId:       req.UserId,
				Points:       0,
				FreezePoints: 0,
			}
			if err := global.DB.Create(&userPoints).Error; err != nil {
				return &proto.AddAndFreezePointsResponse{Success: false, Message: "创建用户积分记录失败"}, nil
			}
		} else {
			return &proto.AddAndFreezePointsResponse{Success: false, Message: "查询用户积分记录失败"}, nil
		}
	}

	tx := global.DB.Begin()

	// 冻结积分逻辑
	if result := tx.Model(&model.UserPoints{}).Where(&model.UserPoints{UserId: req.UserId}).Update("freeze_points", gorm.Expr("freeze_points + ?", req.Points)); result.RowsAffected == 0 {
		tx.Rollback()
		return &proto.AddAndFreezePointsResponse{Success: false, Message: "冻结积分失败"}, nil
	}

	// 记录积分变动
	pointsTransaction := model.PointsTransaction{
		UserId:    req.UserId,
		OrderSn:   req.OrderSn,
		Change:    req.Points,
		Status:    1, // 1 表示等待支付
		Timestamp: time.Now().Unix(),
	}
	if result := tx.Create(&pointsTransaction); result.RowsAffected == 0 {
		tx.Rollback()
		return &proto.AddAndFreezePointsResponse{Success: false, Message: "记录积分变动失败"}, nil
	}

	tx.Commit()
	return &proto.AddAndFreezePointsResponse{Success: true, Message: "积分冻结成功"}, nil
}

/**
 * @Author: weihua hu
 * @description: 解冻积分
 * @param {context.Context} ctx
 * @param {*proto.UnfreezePointsRequest} req
 * @return {*}
 */
func (s *PointsServer) UnfreezePoints(ctx context.Context, req *proto.UnfreezePointsRequest) (*proto.UnfreezePointsResponse, error) {
	tx := global.DB.Begin()

	// 解冻积分逻辑
	var pointsTransaction model.PointsTransaction
	if result := tx.Model(&model.PointsTransaction{}).Where(&model.PointsTransaction{OrderSn: req.OrderSn, UserId: req.UserId}).First(&pointsTransaction); result.RowsAffected == 0 {
		tx.Rollback()
		return &proto.UnfreezePointsResponse{Success: false, Message: "未找到对应的积分变动记录"}, nil
	}

	if result := tx.Model(&model.UserPoints{}).Where(&model.UserPoints{UserId: req.UserId}).Update("points", gorm.Expr("points + ?", pointsTransaction.Change)).Update("freeze_points", gorm.Expr("freeze_points - ?", pointsTransaction.Change)); result.RowsAffected == 0 {
		tx.Rollback()
		return &proto.UnfreezePointsResponse{Success: false, Message: "解冻积分失败"}, nil
	}

	// 更新积分变动记录状态
	if result := tx.Model(&model.PointsTransaction{}).Where(&model.PointsTransaction{OrderSn: req.OrderSn, UserId: req.UserId}).Update("status", 2); result.RowsAffected == 0 {
		tx.Rollback()
		return &proto.UnfreezePointsResponse{Success: false, Message: "更新积分变动记录状态失败"}, nil
	}

	tx.Commit()
	return &proto.UnfreezePointsResponse{Success: true, Message: "积分解冻成功"}, nil
}

/**
 * @Author: weihua hu
 * @description: 订单支付失败后回滚冻结积分
 * @param {context.Context} ctx
 * @param {...*primitive.MessageExt} msgs
 * @return {consumer.ConsumeResult, error}
 */
func AutoRebackPoints(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderRebackMessage struct {
		OrderSn string `json:"order_sn"`
	}

	for i := range msgs {
		var orderRebackMessage OrderRebackMessage
		err := json.Unmarshal(msgs[i].Body, &orderRebackMessage)
		if err != nil {
			zap.S().Errorf("解析json失败： %v\n", msgs[i].Body)
			return consumer.ConsumeRetryLater, err // 返回消费失败
		}

		// 根据 OrderSn 查询相关记录
		var pointsTransaction model.PointsTransaction
		if err := global.DB.Where(&model.PointsTransaction{OrderSn: orderRebackMessage.OrderSn}).First(&pointsTransaction).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				zap.S().Errorf("未找到对应的积分变动记录： %v\n", orderRebackMessage.OrderSn)
				continue
			} else {
				zap.S().Errorf("查询积分变动记录失败： %v\n", err)
				return consumer.ConsumeRetryLater, err
			}
		}

		if pointsTransaction.Status == 3 {
			zap.S().Infof("积分已回滚，无需重复处理： %v\n", orderRebackMessage.OrderSn)
			continue // 跳过已处理的消息
		}

		// 获取 UserId 和其他必要的信息
		userId := pointsTransaction.UserId
		points := pointsTransaction.Change

		fmt.Println("幂等失败")

		// 回滚冻结积分逻辑
		tx := global.DB.Begin()
		if result := tx.Model(&model.UserPoints{}).Where(&model.UserPoints{UserId: userId}).Update("freeze_points", gorm.Expr("freeze_points - ?", points)); result.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Errorf("回滚冻结积分失败： %v\n", userId)
			return consumer.ConsumeRetryLater, nil
		}

		// 更新积分变动记录状态
		if result := tx.Model(&model.PointsTransaction{}).Where(&model.PointsTransaction{OrderSn: orderRebackMessage.OrderSn}).Update("status", 3); result.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Errorf("更新积分变动记录状态失败： %v\n", orderRebackMessage.OrderSn)
			return consumer.ConsumeRetryLater, nil
		}

		tx.Commit()
	}

	return consumer.ConsumeSuccess, nil
}

/**
 * @Author: weihua hu
 * @description: 支付成功后解冻积分
 * @param {context.Context} ctx
 * @param {...*primitive.MessageExt} msgs
 * @return {consumer.ConsumeResult, error}
 */
func AutoUnfreezePoints(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderPaidMessage struct {
		OrderSn string `json:"order_sn"`
	}

	for i := range msgs {
		var orderPaidMessage OrderPaidMessage
		err := json.Unmarshal(msgs[i].Body, &orderPaidMessage)
		if err != nil {
			zap.S().Errorf("解析json失败： %v\n", msgs[i].Body)
			return consumer.ConsumeRetryLater, err // 返回消费失败
		}
		fmt.Println("OrderSn:", orderPaidMessage.OrderSn)

		// 根据 OrderSn 查询相关记录
		var pointsTransaction model.PointsTransaction
		if err := global.DB.Where("order_sn = ?", orderPaidMessage.OrderSn).First(&pointsTransaction).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				zap.S().Errorf("未找到对应的积分变动记录： %v\n", orderPaidMessage.OrderSn)
				return consumer.ConsumeRetryLater, err // 返回消费失败
			} else {
				zap.S().Errorf("查询积分变动记录失败： %v\n", err)
				return consumer.ConsumeRetryLater, err // 返回消费失败
			}
		}

		// 检查幂等性
		if pointsTransaction.Status == 2 {
			zap.S().Infof("积分已解冻，无需重复处理： %v\n", orderPaidMessage.OrderSn)
			continue // 跳过已处理的消息
		}

		zap.S().Info("PointsTransaction: ", pointsTransaction)

		// 获取 UserId 和其他必要的信息
		userId := pointsTransaction.UserId
		points := pointsTransaction.Change

		fmt.Println("userId: ", userId)
		fmt.Println("points: ", points)

		// 解冻积分逻辑
		tx := global.DB.Begin()
		if result := tx.Model(&model.UserPoints{}).Where(&model.UserPoints{UserId: userId}).Update("points", gorm.Expr("points + ?", points)).Update("freeze_points", gorm.Expr("freeze_points - ?", points)); result.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Errorf("解冻积分失败： %v\n", userId)
			return consumer.ConsumeRetryLater, nil // 返回消费失败
		}

		// 更新积分变动记录状态
		if result := tx.Model(&model.PointsTransaction{}).Where(&model.PointsTransaction{OrderSn: orderPaidMessage.OrderSn, UserId: userId}).Update("status", 2); result.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Errorf("更新积分变动记录状态失败： %v\n", orderPaidMessage.OrderSn)
			return consumer.ConsumeRetryLater, nil // 返回消费失败
		}

		tx.Commit()
	}

	return consumer.ConsumeSuccess, nil
}

func (s *PointsServer) GetPointsDetails(ctx context.Context, req *proto.GetPointsDetailsRequest) (*proto.GetPointsDetailsResponse, error) {
	var pointsTransactions []model.PointsTransaction
	if err := global.DB.Where("user_id = ? AND status = ?", req.UserId, 2).Order("timestamp").Find(&pointsTransactions).Error; err != nil {
		return nil, err
	}

	for _, pt := range pointsTransactions {
		fmt.Println("pt: ", pt)
	}

	var pointsDetails []*proto.PointsDetail
	var currentBalance int32

	// 获取用户当前的总积分
	var userPoints model.UserPoints
	if err := global.DB.Where("user_id = ?", req.UserId).First(&userPoints).Error; err != nil {
		return nil, err
	}
	currentBalance = userPoints.Points

	// 逆序遍历积分变动记录，计算每次变动后的积分余额
	for i := len(pointsTransactions) - 1; i >= 0; i-- {
		pt := pointsTransactions[i]
		pointsDetails = append([]*proto.PointsDetail{
			{
				UserId:    pt.UserId,
				OrderSn:   pt.OrderSn,
				Change:    pt.Change,
				Balance:   currentBalance,
				Timestamp: pt.Timestamp,
			},
		}, pointsDetails...)
		currentBalance -= pt.Change
	}

	return &proto.GetPointsDetailsResponse{
		PointsDetails: pointsDetails,
	}, nil
}
