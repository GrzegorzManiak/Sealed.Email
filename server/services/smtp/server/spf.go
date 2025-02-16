package server

import (
	"blitiri.com.ar/go/spf"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"net"
	"strings"
)

func GetRemoteConnectionIp(smtpCtx *Session) string {
	ipStr := smtpCtx.Ctx.Conn().RemoteAddr().String()
	if strings.Contains(ipStr, ":") {
		ipStr = strings.Split(ipStr, ":")[0]
	}
	return ipStr
}

func ValidateMailFromSpf(ipStr, from, host string) (spf.Result, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return spf.None, fmt.Errorf("failed to parse IP: %s", ipStr)
	}
	return spf.CheckHostWithSender(ip, host, from)
}

func SpfShouldBlock(result spf.Result) bool {
	return (result == spf.Fail || result == spf.SoftFail) && config.Smtp.SpfHardFail
}
