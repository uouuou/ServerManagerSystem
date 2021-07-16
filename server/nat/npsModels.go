package nat

type NpsConfig struct {
	common    common    `ini:"common"`
	web       web       `ini:"web"`
	tcp       tcp       `ini:"tcp"`
	udp       udp       `ini:"udp"`
	http      httpMode  `ini:"http"`
	socks5    socks5    `ini:"socks_5"`
	passProxy passProxy `ini:"pass_proxy"`
	p2p       p2p       `ini:"p_2_p"`
	file      file      `ini:"file"`
}

// 全局配置
type common struct {
	ServerAddr       string `json:"server_addr" ini:"server_addr"`             //服务端ip/域名:port
	ConnType         string `json:"conn_type" ini:"conn_type"`                 //与服务端通信模式(tcp或kcp)
	Vkey             int    `json:"vkey" ini:"vkey"`                           //服务端配置文件中的密钥(非web)
	Username         string `json:"username" ini:"username"`                   // socks5或http(s)密码保护用户名(可忽略)
	Password         string `json:"password" ini:"password"`                   // socks5或http(s)密码保护密码(可忽略)
	Compress         bool   `json:"compress" ini:"compress"`                   // 是否压缩传输(true或false或忽略)
	Crypt            bool   `json:"crypt" ini:"crypt"`                         // 是否加密传输(true或false或忽略)
	RateLimit        int    `json:"rate_limit" ini:"rate_limit"`               // 速度限制，可忽略
	FlowLimit        int    `json:"flow_limit" ini:"flow_limit"`               // 流量限制，可忽略
	Remark           string `json:"remark" ini:"remark"`                       // 客户端备注，可忽略
	MaxConn          int    `json:"max_conn" ini:"max_conn"`                   // 最大连接数，可忽略
	AutoReconnection bool   `json:"auto_reconnection" ini:"auto_reconnection"` //断线重连
}

//域名代理
type web struct {
	Host       string `json:"host" ini:"host"`               //域名(http
	TargetAddr string `json:"target_addr" ini:"target_addr"` //内网目标，负载均衡时多个目标，逗号隔开
	HostChange string `json:"host_change" ini:"host_change"` // 请求host修改
}

//tcp隧道模式
type tcp struct {
	Mode        string `json:"mode" ini:"mode"`                 // tcp
	ServerPort  int    `json:"server_port" ini:"server_port"`   // 在服务端的代理端口
	TartgetAddr string `json:"tartget_addr" ini:"tartget_addr"` //内网目标
}

//udp隧道模式
type udp struct {
	Mode        string `json:"mode" ini:"mode"`                 // udp
	ServerPort  int    `json:"server_port" ini:"server_port"`   // 在服务端的代理端口
	TartgetAddr string `json:"tartget_addr" ini:"tartget_addr"` //内网目标
}

//http代理模式
type httpMode struct {
	Mode       string `json:"mode" ini:"mode"`               // http
	ServerPort int    `json:"server_port" ini:"server_port"` // 在服务端的代理端口
}

//socks5代理模式
type socks5 struct {
	Mode         string `json:"mode" ini:"mode"`                   // socks5
	ServerPort   int    `json:"server_port" ini:"server_port"`     // 在服务端的代理端口
	MultiAccount string `json:"multi_account" ini:"multi_account"` //socks5多账号配置文件（可选),配置后使用basic_username和basic_password无法通过认证
}

//私密代理模式
type passProxy struct {
	Mode       string `json:"mode" ini:"mode"`               // secret
	Password   int    `json:"password" ini:"password"`       // 在服务端的代理端口
	TargetAddr string `json:"target_addr" ini:"target_addr"` //内网目标
}

//p2p代理模式
type p2p struct {
	Mode       string `json:"mode" ini:"mode"`               // p2p
	Password   int    `json:"password" ini:"password"`       // 在服务端的代理端口
	TargetAddr string `json:"target_addr" ini:"target_addr"` //内网目标
}

//文件访问模式
type file struct {
	Mode       string `json:"mode" ini:"mode"`               // file
	ServerPort int    `json:"server_port" ini:"server_port"` // 在服务端的代理端口
	LocalPath  string `json:"local_path" ini:"local_path"`   //本地文件目录
	StripPre   string `json:"strip_pre" ini:"strip_pre"`     //前缀
}
