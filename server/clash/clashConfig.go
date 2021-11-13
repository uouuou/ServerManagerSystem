package clash

import (
	"fmt"
	"github.com/Dreamacro/clash/adapter/provider"
	"github.com/Dreamacro/clash/component/auth"
	"github.com/Dreamacro/clash/component/fakeip"
	"github.com/Dreamacro/clash/component/trie"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/dns"
	"github.com/Dreamacro/clash/log"
	T "github.com/Dreamacro/clash/tunnel"
	mid "github.com/uouuou/ServerManagerSystem/middleware"
	"gopkg.in/yaml.v2"
	"net"
	"os"
)

// General config
type General struct {
	Inbound
	Controller
	Mode      T.TunnelMode `json:"mode" json:"mode"`
	LogLevel  log.LogLevel `json:"log-level" json:"log-level"`
	IPv6      bool         `json:"ipv6" json:"ipv6"`
	Interface string       `json:"interface-name" json:"interface-name"`
}

// Inbound config
type Inbound struct {
	Port           int      `json:"port" json:"port"`
	SocksPort      int      `json:"socks-port" json:"socks-port"`
	RedirPort      int      `json:"redir-port" json:"redir-port"`
	TProxyPort     int      `json:"tproxy-port" json:"t-proxy-port"`
	MixedPort      int      `json:"mixed-port" json:"mixed-port"`
	Authentication []string `json:"authentication" json:"authentication"`
	AllowLan       bool     `json:"allow-lan" json:"allow-lan"`
	BindAddress    string   `json:"bind-address" json:"bind-address"`
}

// Controller config
type Controller struct {
	ExternalController string `json:"-"`
	ExternalUI         string `json:"-"`
	Secret             string `json:"-"`
}

// DNS config
type DNS struct {
	Enable            bool             `yaml:"enable,omitempty" json:"enable"`
	IPv6              bool             `yaml:"ipv6,omitempty" json:"ipv6"`
	NameServer        []dns.NameServer `yaml:"nameserver,omitempty" json:"nameserver"`
	Fallback          []dns.NameServer `yaml:"fallback,omitempty" json:"fallback"`
	FallbackFilter    FallbackFilter   `yaml:"fallback-filter,omitempty" json:"fallback-filter"`
	Listen            string           `yaml:"listen,omitempty" json:"listen"`
	EnhancedMode      dns.EnhancedMode `yaml:"enhanced-mode,omitempty" json:"enhanced-mode"`
	DefaultNameserver []dns.NameServer `yaml:"default-nameserver,omitempty" json:"default-nameserver"`
	FakeIPRange       *fakeip.Pool
	Hosts             *trie.DomainTrie
}

// FallbackFilter config
type FallbackFilter struct {
	GeoIP  bool         `yaml:"geoip,omitempty" json:"geoip"`
	IPCIDR []*net.IPNet `yaml:"ipcidr,omitempty" json:"ipcidr"`
	Domain []string     `yaml:"domain,omitempty" json:"domain"`
}

// Experimental config
type Experimental struct{}

// Config is clash config manager
type Config struct {
	General      *General
	DNS          *DNS
	Experimental *Experimental
	Hosts        *trie.DomainTrie
	Rules        []C.Rule
	Users        []auth.AuthUser
	Proxies      map[string]C.Proxy
	Providers    map[string]provider.ProxyProvider
}

type RawDNS struct {
	Enable            bool              `yaml:"enable,omitempty" json:"enable"`
	IPv6              bool              `yaml:"ipv6,omitempty" json:"ipv6"`
	UseHosts          bool              `yaml:"use-hosts,omitempty" json:"use-hosts"`
	NameServer        []string          `yaml:"nameserver,omitempty" json:"nameserver"`
	Fallback          []string          `yaml:"fallback,omitempty" json:"fallback"`
	FallbackFilter    RawFallbackFilter `yaml:"fallback-filter,omitempty" json:"fallback-filter"`
	Listen            string            `yaml:"listen,omitempty" json:"listen"`
	EnhancedMode      dns.EnhancedMode  `yaml:"enhanced-mode,omitempty" json:"enhanced-mode"`
	FakeIPRange       string            `yaml:"fake-ip-range,omitempty" json:"fake-ip-range"`
	FakeIPFilter      []string          `yaml:"fake-ip-filter,omitempty" json:"fake-ip-filter"`
	DefaultNameserver []string          `yaml:"default-nameserver,omitempty" json:"default-nameserver"`
}

type RawFallbackFilter struct {
	GeoIP  bool     `yaml:"geoip,omitempty" json:"geoip"`
	IPCIDR []string `yaml:"ipcidr,omitempty" json:"ipcidr"`
	Domain []string `yaml:"domain,omitempty" json:"domain"`
}

type RawConfig struct {
	Port               int          `yaml:"port,omitempty" json:"port"`
	SocksPort          int          `yaml:"socks-port,omitempty" json:"socks-port"`
	RedirPort          int          `yaml:"redir-port,omitempty" json:"redir-port"`
	TProxyPort         int          `yaml:"tproxy-port,omitempty" json:"tproxy-port"`
	MixedPort          int          `yaml:"mixed-port,omitempty" json:"mixed-port"`
	Authentication     []string     `yaml:"authentication,omitempty" json:"authentication"`
	AllowLan           bool         `yaml:"allow-lan,omitempty" json:"allow-lan"`
	BindAddress        string       `yaml:"bind-address,omitempty" json:"bind-address"`
	Mode               T.TunnelMode `yaml:"mode,omitempty" json:"mode"`
	LogLevel           log.LogLevel `yaml:"log-level,omitempty" json:"log-level"`
	IPv6               bool         `yaml:"ipv6,omitempty" json:"ipv6"`
	ExternalController string       `yaml:"external-controller,omitempty" json:"external-controller"`
	ExternalUI         string       `yaml:"external-ui,omitempty" json:"external-ui"`
	Secret             interface{}  `yaml:"secret,omitempty" json:"secret"`
	Interface          string       `yaml:"interface-name,omitempty" json:"interface-name"`

	ProxyProvider map[string]map[string]interface{} `yaml:"proxy-providers,omitempty" json:"proxy-provider"`
	Hosts         map[string]string                 `yaml:"hosts,omitempty" json:"hosts"`
	DNS           RawDNS                            `yaml:"dns,omitempty" json:"dns"`
	Experimental  Experimental                      `yaml:"experimental,omitempty" json:"experimental"`
	Proxy         []map[string]interface{}          `yaml:"proxies,omitempty" json:"proxies"`
	ProxyGroup    []map[string]interface{}          `yaml:"proxy-groups,omitempty" json:"proxy-groups"`
	Rule          []string                          `yaml:"rules,omitempty" json:"rules"`
}

func ReadConfig() {
	var rawConfig RawConfig
	config, err := os.ReadFile(mid.Dir + "/config/clash.yaml")
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(config, &rawConfig)
	if err != nil {
		mid.Log.Error(err.Error())
	}
	if rawConfig.BindAddress == "" {
		rawConfig.BindAddress = "0.0.0.0"
	}
	if rawConfig.ExternalUI != "" {
		rawConfig.ExternalUI = ""
	}
	data, _ := yaml.Marshal(rawConfig)
	_ = os.WriteFile(mid.Dir+"/config/configClash.yaml", data, 0777)
}
