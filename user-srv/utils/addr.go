/*
 * @Author: weihua hu
 * @Date: 2024-11-25 01:02:42
 * @LastEditTime: 2024-11-25 01:02:42
 * @LastEditors: weihua hu
 * @Description:
 */
package utils

import (
	"net"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}