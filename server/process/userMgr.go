package process2

import "fmt"

// UserMgr结构体实例 在服务器有且只有一个
// 在其他地方，都会用到 故定义成全局变量
var userMgr *UserMgr

type UserMgr struct {
	// 在线用户map 用户id字段 用户登录注册结构体
	onlineUsers map[int]*UserProcess
}

// 初始化userMgr
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 定义方法 onlineUsers字段对应的map中 添加数据方法
// 参数 处理用户登录注册对应的结构体
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	// 把对应的 用户登录注册信息 添加到map中
	this.onlineUsers[up.UserId] = up
}

// 定义方法 从onlineUsers字段对应的map中 删除数据方法
// 根据传递的用户id删除 用户
func (this *UserMgr) DelOnlineUser(userId int) {
	// 从onlineUsers map中删除 指定id的用户
	delete(this.onlineUsers, userId)
}

// 定义方法 查询onlineUsers字段对应的map的所有信息
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers // 返回map
}

// 定义方法 根据id返回 map中该id对应的数据 和 错误
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[userId]
	if !ok { // map中取不到 就返回错误信息
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
