package process

import (
	"encoding/json"
	"fmt"

	"github.com/go-code/chartRoom/common/message"
)

// 管理信息
// 编写方法 显示群发信息
func outputGroupMes(mes *message.Message) {
	// 1. 反序列化mes 取出 mes.Data内容 给SmsMes
	// 实例化SmsMes结构体
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("outputGroupMes() json.Unmarshal err=", err.Error())
		return
	}
	// 2. 显示信息
	info := fmt.Sprintf("用户id:\t%d 对大家说:\t%s\n",
		smsMes.UserId, smsMes.Content)
	fmt.Println(info)
}
