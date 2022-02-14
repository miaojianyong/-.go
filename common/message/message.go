package message

const ( // 定义信息类型 常量
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)
const ( // 定义用户状态 常量
	UserOnline    = iota // 在线
	UserOffOnline        // 离线
)

// 消息结构体
type Message struct {
	Type string `json:"type"` // 消息类型
	Data string `json:"data"` // 消息内容
}

// 定义登录消息 发送
type LoginMes struct {
	UserId   int    `json:"userId"`   // 用户id
	UserPwd  string `json:"userPwd"`  // 用户密码
	UserName string `json:"userName"` // 用户名
}

// 定义登录消息 返回
type LoginResMes struct {
	Code    int    `json:"code"` // 返回状态码 500表示该用户未注册 200表示登录成功
	UsersId []int  // 增加字段 存放登录用户的id
	Error   string `json:"error"` // 返回错误信息
}

// 注册消息 发送
type RegisterMes struct {
	User User `json:"user"` // 类型是User结构体
}

// 注册后 服务器响应消息
type RegisterResMes struct {
	Code  int    `json:"code"`  // 返回状态码 500表示该用户占用 200表示注册成功
	Error string `json:"error"` // 返回错误信息
}

// 用作服务器推送用户状态变化的信息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` // 用户id
	Status int `json:"status"` // 用户状态
}

// 发送信息 结构体
type SmsMes struct {
	Content string `json:"content"` // 内容
	User           // 继承 用户结构体 匿名结构体
}
