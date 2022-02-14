package model

import (
	"encoding/json"
	"fmt"

	"github.com/go-code/chartRoom/common/message"
	"github.com/gomodule/redigo/redis"
)

// 操作 User结构体文件
// 定义全局的userDao实例，需要和redis操作时就启动
var MyUserDao *UserDao

// 定义 UserDao结构体 存在User结构体
type UserDao struct {
	// 要改结构体和 redis连接池有联系
	pool *redis.Pool
}

// 使用工厂模式(构造函数) 创建UserDao实例
// 传递 pool
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 定义方法 根据用户id 返回User实例和错误
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	// 通过id 去redis中查询用户 设置为返回字符串
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		// 判断错误信息类型
		// 如在 users哈希中，没有找到users对应的id
		if err == redis.ErrNil {
			err = ERROR_USER_NOTEXISTS // 把自定义错误信息 赋值给err
		}
		return
	}
	// 如果存在 需要把res反序列化成 User实例
	err = json.Unmarshal([]byte(res), &user) // 传递地址
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

// 定义方法 验证登录
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	// 先从UserDao的连接池中 取一个连接
	conn := this.pool.Get()
	defer conn.Close()
	// 调用根据id操作用户是否存在方法
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return // 即返回了 getUserById方法中定义的错误信息
	}
	// 这时表示用户存在 需校验密码是否正确
	// 如果中redis中取的密码 不等于 传递的密码
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD // 即自定义的密码错误
		return
	}
	return
}

// 定义方法 注册用户 传递user信息 返回错误
func (this *UserDao) Register(user *message.User) (err error) {
	// 先从UserDao的连接池中 取一个连接
	conn := this.pool.Get()
	defer conn.Close()
	// 调用根据id操作用户是否存在方法
	_, err = this.getUserById(conn, user.UserId)
	if err == nil { // 如果错误为空 表示用户找到了用户
		err = ERROR_USER_EXISTS // 返回用户已存在错误
		return
	}
	// 否则 就是有错误 即表示该id在redis中没有
	// 序列化 然后把user信息存储 到数据库
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 存放到数据库中
	_, err = conn.Do("HSet", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("向redis保存注册用户失败，err=", err)
		return
	}
	return
}
