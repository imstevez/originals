package middlewares

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/client"
)

// LogCli go-micro服务调用日志中间件
type LogCli struct {
	client.Client
}

// Call go-micro服务调用日志中间件方法
func (lc *LogCli) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	fmt.Printf("[service] service: %-24s | endpoint: %-24s\n", req.Service(), req.Endpoint())
	return lc.Client.Call(ctx, req, rsp)
}
