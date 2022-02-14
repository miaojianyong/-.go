package main

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// 连接redis数据库 初始化代码

// 定义全局的pool
var pool *redis.Pool

// 编写函数创建 redis连接池
// 1> redis服务器地址
// 2> 最大空闲连接数
// 3> 最大连接数
// 4> 最大空闲时间 时间类型
func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {
	pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address)
		},
	}
}
