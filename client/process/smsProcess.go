package process

import (
	"encoding/json"
	"fmt"

	"github.com/go-code/chartRoom/common/message"
	"github.com/go-code/chartRoom/server/utils"
)

type SmsProcess struct{}

// 对应方法 发送群聊消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	// 创建mes
	var mes message.Message
	mes.Type = message.SmsMesType
	// 创建smsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content               // 内容
	smsMes.UserId = CurUser.UserId         // 是哪个用户
	smsMes.UserStatus = CurUser.UserStatus // 状态
	// 序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err.Error())
		return
	}
	mes.Data = string(data)
	// 对mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err.Error())
		return
	}
	// 调用方法 将序列化后的mes 发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes tf.WritePkg err=", err.Error())
		return
	}
	return
}
