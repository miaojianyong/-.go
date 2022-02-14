package main

import (
	"fmt"
	"os"

	"github.com/go-code/chartRoom/client/process"
)

var (
	userId   int    // 存储用户id
	userPwd  string // 存储用户密码
	userName string // 存储用户名 昵称
)

// 系统登录界面
func main() {
	// 接收用户的选择
	var key int
	// 判断是否继续显示菜单
	// var loop = true
	for {
		fmt.Println("-----------------欢迎登录多人聊天系统-----------------")
		fmt.Println("\t\t 1 登录聊天室")
		fmt.Println("\t\t 2 注册用户")
		fmt.Println("\t\t 3 退出系统")
		fmt.Println("\t\t 请选择(1-3)：")
		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户id：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			// loop = false
			// 1. 创建 UserProcess实例
			up := &process.UserProcess{}
			// 调用 验证登录方法 登录操作
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id：") // 后期可做成redis自增项
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码：")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名字(昵称)：")
			fmt.Scanf("%s\n", &userName)
			// 1. 创建 UserProcess实例
			up := &process.UserProcess{}
			// 调用 注册用户方法 完成注册请求
			up.Register(userId, userPwd, userName)
		case 3:
			fmt.Println("退出系统")
			// loop = false
			os.Exit(0) // 这个也能退出系统
		default:
			fmt.Println("输入有误，请重新输入！")
		}
	}
	// 根据用户输入，显示新的菜单等
	// if key == 1 { // 用户登录
	// 	// fmt.Println("请输入用户id：")
	// 	// fmt.Scanf("%d\n", &userId)
	// 	// fmt.Println("请输入用户密码：")
	// 	// fmt.Scanf("%s\n", &userPwd)
	// 	// 调用 登录验证函数 在login.go文件中
	// 	// 在同一个包 小写的函数可跨文件调用
	// 	login(userId, userPwd)
	// 	/* if err != nil {
	// 		fmt.Println("登录失败=", err)
	// 	} else {
	// 		fmt.Println("登录成功")
	// 	} */
	// } else if key == 2 { // 用户注销
	// 	fmt.Println("进行用户注册逻辑代码")
	// }
}
