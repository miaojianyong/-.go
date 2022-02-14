package main

import (
	"fmt"
	"io"
	"net"

	"github.com/go-code/chartRoom/common/message"
	"github.com/go-code/chartRoom/server/utils"
	process2 "github.com/main.go-code/chartRoom/server/process"
)

// 将下述函数关联到一个结构体中 并改成方法
type Processor struct {
	Conn net.Conn
}

// 编写 serverProcessMes函数
// 根据客户端发送消息种类不同 决定调用哪个函数来处理
func (this *Processor) serverProcessMes(mes *message.Message) (err error) {
	// 打印客户端 发送的群发消息
	fmt.Println("mes=", mes)

	switch mes.Type { // 判断消息类型
	case message.LoginMesType: // 登录消息类型
		// 创建 UserProcess结构体实例
		up := &process2.UserProcess{
			// 因为每个连接都不同 故每个类型都有重新创建
			Conn: this.Conn,
		}
		err = up.ServerProcessLogin(mes) // 处理登录
	case message.RegisterMesType: // 注册消息类型
		up := &process2.UserProcess{
			// 因为每个连接都不同 故每个类型都有重新创建
			Conn: this.Conn,
		}
		err = up.ServerProcessRegister(mes) // 调用注册方法
	case message.SmsMesType: // 群发消息类型
		// 创建SmsProcess实例，调用方法 转发信息
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

// 创建方法 循环等待客户端发送信息
func (this *Processor) process2() (err error) {
	// 读取客户端发送的信息
	for {
		// 把读取客户端发送的数据，封装成函数readPkg()
		// 调用utils中的方法读取数据，即创建Transfer变量
		tf := &utils.Transfer{ // 用地址 并给用到的字段赋值
			Conn: this.Conn,
		}
		// 读取数据
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出了，服务器端也退出")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		fmt.Println("mes=", mes)
		// 调用函数 处理不同类型数据
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
