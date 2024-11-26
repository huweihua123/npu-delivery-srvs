/*
 * @Author: weihua hu
 * @Date: 2024-11-25 00:09:23
 * @LastEditTime: 2024-11-26 23:34:09
 * @LastEditors: weihua hu
 * @Description:
 */

package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"npu-delivery-srvs/user-srv/global"
	"npu-delivery-srvs/user-srv/model"
	"npu-delivery-srvs/user-srv/proto"

	"github.com/anaskhan96/go-password-encoder"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

/**
 * @Author: weihua hu
 * @description:
 * @param {*model.User} user
 * @return {*proto.UserInfoResponse}
 */
func ModelToRsp(user model.User) proto.UserInfoResponse {

	userInfoRsp := proto.UserInfoResponse{
		Id:       int32(user.ID),
		Username: user.Username,
		Password: user.Password,
		Mobile:   user.Mobile,
		Role:     int32(user.Role),
		Email:    user.Email,
		Gender:   user.Gender,
	}

	return userInfoRsp

}

/**
 * @Author: weihua hu
 * @description:
 * @param {model.UserAddress} address
 * @return {proto.UserAddressResponse}
 */
func AddressToRsp(address model.UserAddress) proto.UserAddressResponse {
	Address := proto.UserAddressResponse{
		Id:        int32(address.ID),
		UserId:    int32(address.UserId),
		IsDefault: address.IsDefault,
		Address:   address.Address,
	}

	return Address

}

/**
 * @Author: weihua hu
 * @description:
 * @param {context.Context} ctx
 * @param {*proto.CreateUserInfo} req
 * @return {*proto.UserInfoResponse}
 * @return {error}
 */
func (u *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User

	result := global.DB.Where(&model.User{Username: req.Username}).First(&user)

	if result.RowsAffected > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户名已存在")
	}

	user.Mobile = req.Mobile
	user.Username = req.Username

	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoRsp := ModelToRsp(user)
	return &userInfoRsp, nil

}

/**
 * @Author: weihua hu
 * @description:
 * @param {context.Context} ctx
 * @param {*proto.CreateUserAddressInfo} req
 * @return {*proto.UserAddressResponse} - 用户地址响应，包含新创建的地址信息
 * @return {error} - 错误信息，如果存在的话
 */
func (s *UserServer) CreateUserAddress(ctx context.Context, req *proto.CreateUserAddressInfo) (*proto.UserAddressResponse, error) {
	// 1. 验证用户是否存在
	var user model.User
	if err := global.DB.First(&user, req.UserId).Error; err != nil {
		// 用户不存在
		return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.UserId)
	}

	// 2. 如果地址是默认地址，需要处理其他地址的默认标记
	if req.IsDefault {
		// 将该用户的所有地址的 is_default 字段设置为 false
		if err := global.DB.Model(&model.UserAddress{}).
			Where(&model.UserAddress{UserId: req.UserId, IsDefault: true}).
			Update("is_default", false).Error; err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update other addresses: %v", err)
		}
	}

	// 3. 创建新的地址
	address := model.UserAddress{
		UserId:    req.UserId,
		Address:   req.Address,
		IsDefault: req.IsDefault,
	}

	if err := global.DB.Create(&address).Error; err != nil {
		// 创建地址失败
		return nil, status.Errorf(codes.Internal, "failed to create address: %v", err)
	}

	// 4. 返回新创建的地址信息
	addressResponse := AddressToRsp(address)

	return &addressResponse, nil
}
