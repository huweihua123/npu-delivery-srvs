/*
 * @Author: weihua hu
 * @Date: 2024-11-25 00:09:23
 * @LastEditTime: 2024-11-25 00:55:37
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
 * @return {*proto.UserInfoResponse} userInfoRsp
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
