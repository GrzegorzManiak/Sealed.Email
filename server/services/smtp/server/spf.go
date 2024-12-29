package server

import (
	"context"
	"fmt"
	"github.com/wttw/spf"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

func ValidateMailFromSpf(smtpCtx *Session) spf.Result {
	checker := spf.NewChecker()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ipStr := smtpCtx.Ctx.Conn().RemoteAddr().String()
	if strings.Contains(ipStr, ":") {
		ipStr = strings.Split(ipStr, ":")[0]
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		zap.L().Debug("Failed to parse IP", zap.String("ip", ipStr))
		return spf.Result{Error: fmt.Errorf("failed to parse IP")}
	}

	return checker.SPF(
		ctx,
		ip,
		smtpCtx.From,
		smtpCtx.Ctx.Hostname(),
	)
}

func PrettyPrintSpfResult(result spf.Result) string {
	if result.Error != nil {
		return result.Error.Error()
	}

	switch result.Type {
	case spf.None:
		return "None"
	case spf.Pass:
		return "Pass"
	case spf.Fail:
		return "Fail"
	case spf.Softfail:
		return "Soft Fail"
	case spf.Neutral:
		return "Neutral"
	case spf.Temperror:
		return "Temp Error"
	case spf.Permerror:
		return "Perm Error"
	default:
		return "Unknown"
	}
}
