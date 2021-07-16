package models

import (
	"bytes"
	"encoding/binary"
)

const (
	UnixStaTimestamp = 2208988800
)

//Ntp
/**
NTP协议 http://www.ntp.org/documentation.html
@author mengdj@outlook.com
*/
type Ntp struct {
	//1:32bits
	Li        uint8 //2 bits
	Vn        uint8 //3 bits
	Mode      uint8 //3 bits
	Stratum   uint8
	Poll      uint8
	Precision uint8
	//2:
	RootDelay           int32
	RootDispersion      int32
	ReferenceIdentifier int32
	//64位时间戳
	ReferenceTimestamp uint64 //指示系统时钟最后一次校准的时间
	OriginateTimestamp uint64 //指示客户向服务器发起请求的时间
	ReceiveTimestamp   uint64 //指服务器收到客户请求的时间
	TransmitTimestamp  uint64 //指示服务器向客户发时间戳的时间
}

func NewNtp() (p *Ntp) {
	//其他参数通常都是服务器返回的
	p = &Ntp{Li: 0, Vn: 3, Mode: 3, Stratum: 0}
	return p
}

//GetBytes
/**
构建NTP协议信息
*/
func (thisFun *Ntp) GetBytes() []byte {
	//注意网络上使用的是大端字节排序
	buf := &bytes.Buffer{}
	head := (thisFun.Li << 6) | (thisFun.Vn << 3) | ((thisFun.Mode << 5) >> 5)
	_ = binary.Write(buf, binary.BigEndian, head)
	_ = binary.Write(buf, binary.BigEndian, thisFun.Stratum)
	_ = binary.Write(buf, binary.BigEndian, thisFun.Poll)
	_ = binary.Write(buf, binary.BigEndian, thisFun.Precision)
	//写入其他字节数据
	_ = binary.Write(buf, binary.BigEndian, thisFun.RootDelay)
	_ = binary.Write(buf, binary.BigEndian, thisFun.RootDispersion)
	_ = binary.Write(buf, binary.BigEndian, thisFun.ReferenceIdentifier)
	_ = binary.Write(buf, binary.BigEndian, thisFun.ReferenceTimestamp)
	_ = binary.Write(buf, binary.BigEndian, thisFun.OriginateTimestamp)
	_ = binary.Write(buf, binary.BigEndian, thisFun.ReceiveTimestamp)
	_ = binary.Write(buf, binary.BigEndian, thisFun.TransmitTimestamp)
	//[27 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	return buf.Bytes()
}

func (thisFun *Ntp) Parse(bf []byte, useUnixSec bool) {
	var (
		bit8  uint8
		bit32 int32
		bit64 uint64
		rb    *bytes.Reader
	)
	//貌似这binary.Read只能顺序读，不能跳着读，想要跳着读只能使用切片bf
	rb = bytes.NewReader(bf)
	_ = binary.Read(rb, binary.BigEndian, &bit8)
	//向右偏移6位得到前两位LI即可
	thisFun.Li = bit8 >> 6
	//向右偏移2位,向右偏移5位,得到前中间3位
	thisFun.Vn = (bit8 << 2) >> 5
	//向左偏移5位，然后右偏移5位得到最后3位
	thisFun.Mode = (bit8 << 5) >> 5
	_ = binary.Read(rb, binary.BigEndian, &bit8)
	thisFun.Stratum = bit8
	_ = binary.Read(rb, binary.BigEndian, &bit8)
	thisFun.Poll = bit8
	_ = binary.Read(rb, binary.BigEndian, &bit8)
	thisFun.Precision = bit8

	//32bits
	_ = binary.Read(rb, binary.BigEndian, &bit32)
	thisFun.RootDelay = bit32
	_ = binary.Read(rb, binary.BigEndian, &bit32)
	thisFun.RootDispersion = bit32
	_ = binary.Read(rb, binary.BigEndian, &bit32)
	thisFun.ReferenceIdentifier = bit32

	//以下几个字段都是64位时间戳(NTP都是64位的时间戳)
	_ = binary.Read(rb, binary.BigEndian, &bit64)
	thisFun.ReferenceTimestamp = bit64
	_ = binary.Read(rb, binary.BigEndian, &bit64)
	thisFun.OriginateTimestamp = bit64
	_ = binary.Read(rb, binary.BigEndian, &bit64)
	thisFun.ReceiveTimestamp = bit64
	_ = binary.Read(rb, binary.BigEndian, &bit64)
	thisFun.TransmitTimestamp = bit64
	//转换为unix时间戳,先左偏移32位拿到64位时间戳的整数部分，然后ntp的起始时间戳 1900年1月1日 0时0分0秒 2208988800
	if useUnixSec {
		thisFun.ReferenceTimestamp = (thisFun.ReceiveTimestamp >> 32) - UnixStaTimestamp
		if thisFun.OriginateTimestamp > 0 {
			thisFun.OriginateTimestamp = (thisFun.OriginateTimestamp >> 32) - UnixStaTimestamp
		}
		thisFun.ReceiveTimestamp = (thisFun.ReceiveTimestamp >> 32) - UnixStaTimestamp
		thisFun.TransmitTimestamp = (thisFun.TransmitTimestamp >> 32) - UnixStaTimestamp
	}
}
