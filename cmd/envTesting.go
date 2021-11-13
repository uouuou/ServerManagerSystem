package cmd

import (
	"fmt"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	mod "github.com/uouuou/ServerManagerSystem/models"
	"github.com/uouuou/ServerManagerSystem/server/clash"
	"github.com/uouuou/ServerManagerSystem/util"
	"os"
	"strings"
)

var NftVersion string
var UnboundVersion string
var ClashVersion string

// EnvTesting 检测环境中的基础程序的存在性
func EnvTesting() map[string]string {
	//检测nft是否安装并输出版本号
	info, err := mod.ExecCommandWithResult("nft -v")
	if err != nil {
		mod.InstallPack("nftables")
		info, err = mod.ExecCommandWithResult("nft -v")
		if err != nil {
			info = "unknown version"
		}
		ntfVersion := strings.Fields(strings.TrimSpace(info))
		NftVersion = fmt.Sprintf("%v %v", ntfVersion[0], ntfVersion[1])
	} else {
		ntfVersion := strings.Fields(strings.TrimSpace(info))
		NftVersion = fmt.Sprintf("%v %v", ntfVersion[0], ntfVersion[1])
	}

	//检测unbound的安装
	infoUnbound, err := mod.ExecCommandWithResult("unbound -v")
	if err != nil {
		mod.InstallPack("unbound")
		infoUnbound, err = mod.ExecCommandWithResult("unbound -v")
		if err != nil {
			infoUnbound = "now unbound is not install unknown version"
		}
		unboundVersion := strings.Fields(strings.TrimSpace(infoUnbound))
		UnboundVersion = fmt.Sprintf("%v %v", unboundVersion[5], unboundVersion[6])
	} else {
		unboundVersion := strings.Fields(strings.TrimSpace(infoUnbound))
		UnboundVersion = fmt.Sprintf("%v %v", unboundVersion[5], unboundVersion[6])
	}
	//检测clash版本号
	infoClash, err := mod.ExecCommandWithResult("clash -v")
	if err != nil {
		if clash.UpdateClash("premium") {
			mid.Log.Info("Clash初始化成功")
			infoClash, err = mod.ExecCommandWithResult("clash -v")
			if err != nil {
				infoClash = "unknown version"
			}
			clashVersion := strings.Fields(strings.TrimSpace(infoClash))
			ClashVersion = fmt.Sprintf("%v %v", clashVersion[0], clashVersion[1])
		} else {
			mid.Log.Error("Clash初始化失败")
		}
	} else {
		clashVersion := strings.Fields(strings.TrimSpace(infoClash))
		ClashVersion = fmt.Sprintf("%v %v", clashVersion[0], clashVersion[1])
	}
	//开放SM管理端端口
	mod.OpenPort(util.Port)
	mod.OpenPort(util.RpcPort)
	var m = map[string]string{
		"unbound":  UnboundVersion,
		"nftables": NftVersion,
		"clash":    ClashVersion,
	}
	//检测Nps
	if mod.NpsVersion() == "" {
		mid.Log.Info("Nps Not Install")
	} else {
		mid.Log.Info("Nps Info:" + mod.NpsVersion())
	}
	return m
}

func EnvNftables() {
	nftables := `
ip tuntap add user root mode tun utun
ip link set utun up
ip address replace 172.31.255.253/30 dev utun
ip route replace default dev utun table 114
ip rule del fwmark 114514 lookup 114
ip rule add fwmark 114514 lookup 114

nft -f - << EOF
define LOCAL_SUBNET = {127.0.0.0/8, 224.0.0.0/4, 192.168.0.0/16, 10.0.0.0/8, 172.16.0.0/12, 169.254.0.0/16, 240.0.0.0/4}
define TUN_DEVICE = utun
define FORWARD_DNS_REDIRECT = {127.0.0.1:1053}
table inet clash
flush table inet clash
table inet clash {
    chain local {
        type route hook output priority 0; policy accept;
        
        ip protocol != { tcp, udp } accept
        
        ip daddr \$LOCAL_SUBNET accept
        
        ct state new ct mark set 114514
        ct mark 114514 mark set 114514
    }
    
    chain forward {
        type filter hook prerouting priority 0; policy accept;
        
        ip protocol != { tcp, udp } accept
    
        iif utun accept
        ip daddr \$LOCAL_SUBNET accept
        
        mark set 114514
    }
    
    chain local-dns-redirect {
        type nat hook output priority 0; policy accept;
        
        ip protocol != { tcp, udp } accept
        
        
        udp dport 53 dnat ip to $FORWARD_DNS_REDIRECT
        tcp dport 53 dnat ip to $FORWARD_DNS_REDIRECT
    }
    
    chain forward-dns-redirect {
        type nat hook prerouting priority 0; policy accept;
        
        ip protocol != { tcp, udp } accept
        
        udp dport 53 dnat ip to $FORWARD_DNS_REDIRECT
        tcp dport 53 dnat ip to $FORWARD_DNS_REDIRECT
    }
}

EOF

sysctl -w net/ipv4/ip_forward=1
`
	_ = os.WriteFile(mid.Dir+"/1.sh", []byte(nftables), 0755)

}
