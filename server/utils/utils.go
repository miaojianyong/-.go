package utils

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/go-code/chartRoom/common/message"
)

// 将下述函数关联到一个结构体中 并改成方法
// 因为其他包用到该包中的方法 故方法名改为大写
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte // 传输时 用的缓冲
}

// 读取客户端发送的数据
// 接收连接 返回mesage结构体 和 错误
func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	fmt.Println("读取客户端发送来的数据...")
	// buf := make([]byte, 8096)
	n, err := this.Conn.Read(this.Buf[:4]) // 只读4个字节
	if n != 4 || err != nil {              // n,就是返回读取的字节数 不判断也行
		// fmt.Println("conn.Read err=", err)
		// 自定义错误
		// err = errors.New("read pak header error")
		return // 返回的是 err
	}
	fmt.Println("读到的buf=", this.Buf[:4])
	// 把读取到的 buf[:4]字节 转成 uint32类型
	// 即使用下述函数 包byte切片 转成uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	// 根据 pkgLen读取消息内容 即从conn中读取pkgLen字节的数据存放到buf中
	n, err = this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		err = errors.New("read pak body error")
		return
	}
	// 把pkgLen反序列化成 message.Message类型
	// 注：这里一定要使用地址 &mes
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

// 编写 writePkg函数，验证发送的数据包
func (this *Transfer) WritePkg(data []byte) (err error) {
	// 1.  先发送数据包的长度
	var pakLen uint32
	pakLen = uint32(len(data))
	// var buf [4]byte // 因为uint32占4个字节
	binary.BigEndian.PutUint32(this.Buf[0:4], pakLen)
	n, err := this.Conn.Write(this.Buf[:4]) // 发送4个字节
	if n != 4 || err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
	// 2. 再发送消息本身
	n, err = this.Conn.Write(data) // 发送4个字节
	if n != int(pakLen) || err != nil {
		fmt.Println("conn.Write err=", err)
		return
	}
	return
}
