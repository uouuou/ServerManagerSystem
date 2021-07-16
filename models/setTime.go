package models

import (
	"log"
	"net"
	"time"
)

// SetTime 使用阿里云NTP服务器校准时间
func SetTime() {
	var (
		ntp    *Ntp
		buffer []byte
		err    error
		ret    int
	)
	//链接阿里云NTP服务器,NTP有很多免费服务器可以使用time.windows.com
	conn, err := net.Dial("udp", "ntp1.aliyun.com:123")
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		_ = conn.Close()
	}()
	ntp = NewNtp()
	_, _ = conn.Write(ntp.GetBytes())
	buffer = make([]byte, 2048)
	ret, err = conn.Read(buffer)
	if err == nil {
		if ret > 0 {
			ntp.Parse(buffer, true)
			tm := time.Unix(int64(ntp.TransmitTimestamp), 0)
			UpdateSystemDate(tm.Format("2006-01-02 15:04:05"))
		}
	}
}
