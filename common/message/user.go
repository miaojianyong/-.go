package message

// 用户信息 文件

// 定义用户结构体
type User struct {
	UserId     int    `json:"userId"`
	UserPwd    string `json:"userPwd"`
	UserName   string `json:"userName"`
	UserStatus int    `json:"userStatus"` // 记录用户状态
}
