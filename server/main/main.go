package main

import (
	"fmt"
	"net"
	"time"

	"github.com/go-code/chartRoom/server/model"
)

// 处理和客户端通讯 函数
func process(conn net.Conn) {
	defer conn.Close()
	// 创建Processor变量 调用方法 处理客户端连接
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误=", err)
		return
	}
}

// 初始化 redis连接池 UserDao结构体
func init() {
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

// 编写函数 完成多UserDao的初始化
func initUserDao() {
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	// 当服务器启动时，就去初始化redis连接池
	// initPool("localhost:6379", 16, 0, 300*time.Second)
	// // 调用函数 初始化 UserDao结构体
	// // 因为该函数中使用到pool参数故要在上述函数下方调用
	// initUserDao()
	fmt.Println("服务器在8889端口监听...~~~")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer listen.Close()
	for { // 循环监听 客户端连接服务器
		fmt.Println("等待客户端连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}
		// 连接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}
}
