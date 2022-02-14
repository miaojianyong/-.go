package process2

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/go-code/chartRoom/common/message"
	"github.com/go-code/chartRoom/server/utils"
)

// 处理信息类型

type SmsProcess struct{}

// 编写方法 转发 群发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	// 取出 mes中的内容，即序列化后的Data数据
	// 即反序列化 mse中的Data,然后给smsMes 要使用里面的id
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Unmarshal err=", err.Error())
		return
	}
	// 在序列化mes 用来发送
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes() json.Marshal err=", err.Error())
		return
	}
	// 遍历服务器端维护的 在线用户map 转发消息
	for id, up := range userMgr.onlineUsers {
		// 过滤自己
		if id == smsMes.UserId {
			continue
		}
		// 调用方法 发送消息
		this.SendMesToEachOlineUser(data, up.Conn)
	}
}

// 编写方法 转发给每个在线用户
// 参数接收 序列化后的数据 和 连接
func (this *SmsProcess) SendMesToEachOlineUser(data []byte, conn net.Conn) {
	// 创建Transfer实例 发送数据
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发信息失败 err=", err.Error())
		return
	}
}
