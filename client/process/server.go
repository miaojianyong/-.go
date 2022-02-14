package process

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/go-code/chartRoom/client/utils"
	"github.com/go-code/chartRoom/common/message"
)

// 持续连接服务器 接收服务器返回信息 并显示在页面上

// 显示登陆成功后的界面
func ShowMenu() {
	fmt.Println("-----------------恭喜xxx登陆成功-----------------")
	fmt.Println("\t\t 1. 显示在线用户列表-----------------")
	fmt.Println("\t\t 2. 发送信息-----------------")
	fmt.Println("\t\t 3. 信息列表-----------------")
	fmt.Println("\t\t 4. 退出系统-----------------")
	fmt.Println("请选择(1-4)：")
	var key int
	var content string // 发送消息
	// 创建smsProcess实例 方便调用该结构体的方法
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		outputOnlineUser()
	case 2:
		fmt.Println("请输入对大家数的内容：")
		fmt.Scanf("%s\n", &content)
		// 调用方法 群发消息
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("退出系统~")
		os.Exit(0)
	default:
		fmt.Println("输入有误，请重新输入")
	}
}

// 定义和服务器保持通讯的函数
func serverProcessMes(conn net.Conn) {
	// 创建Transfer实例，调用读取服务器数据方法
	tf := &utils.Transfer{
		Conn: conn,
	}
	// 循环接收服务器推送的数据
	for {
		fmt.Println("客户端正在读取服务器发送的数据")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() err=", err)
			return
		}
		fmt.Println("mes=", mes)
		// 获取服务器发送来的消息类型
		switch mes.Type {
		// 有人上线
		case message.NotifyUserStatusMesType:
			// 取出返回的 用户状态信息内容 Data
			var nus message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &nus)
			// 调用函数 显示当前用户列表
			updateUserStatus(&nus)
		case message.SmsMesType: // 群发信息类型
			// 调用函数 显示消息
			outputGroupMes(&mes)
		default:
			fmt.Println("服务器返回了未知消息类型")
		}
	}
}
