package server

import (
	"blitiri.com.ar/go/spf"
	"go.uber.org/zap"
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
		zap.L().Debug("Failed to parse IP", zap.String("ip", ipStr))
	}

	result, err := spf.CheckHostWithSender(ip, host, from)
	if err != nil {
		zap.L().Debug("Failed to validate SPF", zap.Error(err))
		return result, err
	}

	return result, nil
}

func SpfShouldBlock(result spf.Result) bool {
	return result == spf.Fail || result == spf.SoftFail
}
