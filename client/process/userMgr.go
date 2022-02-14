package process

import (
	"fmt"

	"github.com/go-code/chartRoom/client/model"
	"github.com/go-code/chartRoom/common/message"
)

// 客户端的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)

// 初始化CurUser结构体，声明成全局的
var CurUser model.CurUser // 当用户登录成功后，去初始化该结构体

// 在客户端显示当前在线用户
func outputOnlineUser() {
	// 变量map
	fmt.Println("当前在线用户列表")
	for id, user := range onlineUsers {
		fmt.Printf("用户%v\t用户id%v\n", user, id)
	}
}

// 调用函数 接收用户状态结构体 返回该结构体
func updateUserStatus(nus *message.NotifyUserStatusMes) {
	// 判断map中是否用该用户
	user, ok := onlineUsers[nus.UserId]
	if !ok { // 如果没有
		// 就创建user 并赋值id
		user = &message.User{
			UserId: nus.UserId,
		}
	}
	// 如果有 就设置状态
	user.UserStatus = nus.Status
	// 然后把 user设置好的用户状态信息 存放到map
	onlineUsers[nus.UserId] = user
	// 调用函数 显示当前在线用户列表
	outputOnlineUser()
}
