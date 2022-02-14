package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/go-code/chartRoom/common/message"
	"github.com/go-code/chartRoom/server/utils"
)

// 定义结构体 并生成方法
type UserProcess struct {
}

// 定义方法 注册用户
func (this *UserProcess) Register(userId int,
	userPwd, userName string) (err error) {
	// 1. 连接到服务器 -- 共用的代码 都可封装成函数
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer conn.Close()
	// 2. 通过conn发送信息 给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 3. 创建 RegisterMes 结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	// 4. 将 registerMes 序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 5. 把data赋给 mes.Data字段
	mes.Data = string(data)
	// 6. 将 mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 7. 把data数据发送给服务器
	// 使用utils包中定义的方法
	// 7.1 创建 Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	// 7.2 发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息失败 err=", err)
		return
	}
	// 7.3 读取服务器返回的信息
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("接收注册信息失败 err=", err)
		return
	}
	// 8. 将 mes中的Data字段 对应的消息，反序列化 RegisterResMes
	var registerResMes message.RegisterResMes // 定义结构体变量
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登录")
		os.Exit(0) // 退出程序
	} else { // 错误信息直接返回即可
		fmt.Println(registerResMes.Error)
		os.Exit(0) // 退出程序
	}
	return
}

// 和服务器端的 userProcess.go文件对应
// 登录
// 编写函数 验证登录
func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	// 开始定义协议
	// fmt.Printf("userId=%d userPwd=%s\n", userId, userPwd)
	// return nil
	// 1. 连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}
	defer conn.Close()
	// 2. 通过conn发送信息 给服务器
	// 创建 Message结构体
	var mes message.Message
	// 给 类型字段 赋值
	mes.Type = message.LoginMesType
	// 3. 创建 LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	// 把 LoginMes结构体 给 mes的Data字段
	// mes.Data = loginMes 不能直接给
	// 4. 将 LoginMes 序列化后给 Message的Data字段
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 因为序列化后的 data是 []byte切片 故转换
	mes.Data = string(data)
	// 5. 将 mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 6. 把data数据发送给服务器
	// 6.1 向发送 data的长度
	// 先获取 data的长度 转成一个表示长度的byte切片
	var pakLen uint32
	pakLen = uint32(len(data))
	var buf [4]byte // 因为uint32占4个字节
	binary.BigEndian.PutUint32(buf[0:4], pakLen)
	n, err := conn.Write(buf[:4]) // 发送4个字节
	if n != 4 || err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
	fmt.Printf("客户端发送信息长度成功,len=%d,内容=%s\n",
		len(data), string(data))
	// 6.2 发送消息本身 即用户名和密码
	_, err = conn.Write(data) // data已经序列化了
	if err != nil {
		fmt.Println("conn.Write(data) err=", err)
		return
	}
	// 接收服务器端返回的消息 调用读完数据 函数
	// 创建 Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("tf.ReadPkg() err=", err)
		return
	}
	// 将 mes中的Data字段 对应的消息，反序列化 LoginResMes
	var loginResMes message.LoginResMes // 定义结构体变量
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		// 用户登录成功 初始curUser结构体
		CurUser.Conn = conn // 把当前连接赋值给对应字段
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline
		// 用户登录成功 显示当前在线用户列表
		// 遍历 loginResMes结构体 对应的UsersId字段的切片
		fmt.Println("当前用户在线列表：")
		for _, v := range loginResMes.UsersId {
			//  在线用户列表 不显示自己的id
			if v == userId {
				continue
			}
			fmt.Printf("用户id= %v\n", v)
			// 把当前在线用户 存放到客户端维护的map中
			// 即初始化 onlineUsers map
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOffOnline, // 在线
			}
			onlineUsers[v] = user
		}
		fmt.Println("---------------------")
		// fmt.Println("登录成功")
		// 在客户端应该启动一个协程
		// 该协程保持和服务器端的通讯 如果服务器有数据推送给客户端
		// 则接收并显示在客户端的终端
		go serverProcessMes(conn)
		// 1. 调用函数 循环显示登陆成功的菜单
		for { // 也可把for循环写到ShowMenu函数里面
			ShowMenu()
		}
	} else { // 错误信息直接返回即可
		fmt.Println(loginResMes.Error)
	}
	// else if loginResMes.Code == 500 {
	// 	fmt.Println(loginResMes.Error)
	// }
	return
}
