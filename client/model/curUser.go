package model

import (
	"net"

	"github.com/go-code/chartRoom/common/message"
)

// 当前用户的连接 信息
type CurUser struct {
	Conn net.Conn
	message.User
}
