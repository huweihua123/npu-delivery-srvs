/*
 * @Author: weihua hu
 * @Date: 2024-11-27 21:22:28
 * @LastEditTime: 2024-11-27 21:23:02
 * @LastEditors: weihua hu
 * @Description:
 */

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	cwd, err := os.Getwd()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Current working directory:", cwd)
}