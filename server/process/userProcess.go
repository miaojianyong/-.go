package process2

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/go-code/chartRoom/common/message"
	"github.com/go-code/chartRoom/server/model"
	"github.com/go-code/chartRoom/server/utils"
)

// 将下述函数关联到一个结构体中 并改成方法
type UserProcess struct {
	Conn   net.Conn
	UserId int // 增加字段 表示Conn连接是哪个用户的
}

// 定义方法 把用户状态通知所有在线用户
// 参数 用户id, 及当前登录的用户需要通知其他在线用户 我上线了
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	// 遍历 onlineUsers map,然后一个一个的发送 NotifyUserStatusMes 即用户状态信息结构体
	for id, up := range userMgr.onlineUsers {
		// 过滤掉自己
		if id == userId { // 即如果遍历的id等于传递的id就跳过
			continue
		}
		// 调用方法 开始通知所有在线用户
		up.NotifyMeOnline(userId)
	}
}

// 定义方法 通知其他所有在线用户我上线了
func (this *UserProcess) NotifyMeOnline(userId int) {
	// 组装 Message 消息结构体
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	// 返回的Data信息内容是 NotifyUserStatusMes
	// 故实例化 NotifyUserStatusMes结构体
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	// 再将 notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 然后将序列化后的data切片转换后 给mes的Data字段
	mes.Data = string(data)
	// 在对组装好的mes序列化，然后发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 创建Transfer实例，调用方法 发送
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

// 定义方法 处理用户注册请求
func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 1. 先从mes中 取出mes.Data数据，并反序列化成 RegisterMes结构体
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	// 2. 声明Message结构体变量
	var resMes message.Message
	resMes.Type = message.RegisterMesType // 赋值类型
	// 3. 声明RegisterResMes结构体变量 返回的结构体
	var registerResMes message.RegisterResMes
	// 调用方法 去redis中完成用户注册
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 400
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 404
			registerResMes.Error = "注册用户未知错误"
		}
	} else {
		registerResMes.Code = 200
	}
	// 5. 将 registerResMes序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	// 6. 将 data赋值给 resMes的Data字段
	resMes.Data = string(data)
	// 7. 在对 resMes序列化 然后发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	// 8. 发送 因为发送数据也需要验证数据包 故封装成函数
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

// 编写方法 处理登录请求
func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1. 先从mes中 取出mes.Data数据，并反序列化成 LoginMes结构体
	var LoginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &LoginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	// 2. 声明Message结构体变量
	var resMes message.Message
	resMes.Type = message.LoginResMesType // 赋值类型
	// 3. 声明LoginResMes结构体变量
	var loginResMes message.LoginResMes
	// 4. 对LoginResMes结构体 对应字段赋值
	// 如果id=100,pwd=123456 就合法 否则不合法
	/* if LoginMes.UserId == 100 && LoginMes.UserPwd == "123456" {
		loginResMes.Code = 200 // 合法
	} else {
		loginResMes.Code = 500 // 不合法 即用户不存在
		loginResMes.Error = "该用户不存在，请注册再使用"
	} */
	// 调用函数 去redis中取验证用户信息
	// 传递id 密码
	user, err := model.MyUserDao.Login(LoginMes.UserId, LoginMes.UserPwd)
	if err != nil {
		// 判断错误信息类型 返回自定义错误信息
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500 // 用户不存在
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403 // 密码错误
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505 // 未知错误
			loginResMes.Error = "服务器内部错误..."
		}
	} else {
		loginResMes.Code = 200
		// 用户登录成功 调用方法 把该用户存放到userMgr的map中
		// 给 UserProcess结构体的userId字段赋值
		this.UserId = LoginMes.UserId
		// 传递 this,即当前登录成功的用户对应的 UserProcess结构体
		userMgr.AddOnlineUser(this)
		// 调用方法 通知其他在线用户，我上线了
		this.NotifyOthersOnlineUser(LoginMes.UserId)
		// 将当前在线的用户id 存放到要返回的 登录消息结构体 对应的UsersId字段中
		// 遍历 userMgr.onlineUsers 即在线用户map
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println(user, "登录成功")
	}
	// 5. 将 loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	// 6. 将 data赋值给 resMes的Data字段
	resMes.Data = string(data)
	// 7. 在对 resMes序列化 然后发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	// 8. 发送 因为发送数据也需要验证数据包 故封装成函数
	// 调用utils中的方法，即创建Transfer变量
	tf := &utils.Transfer{ // 用地址 并给用到的字段赋值
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
