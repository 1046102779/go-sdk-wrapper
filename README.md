## go-sdk-wrapper

减少业务开发者对dapr标准化API的协议理解，只关注自己的payload。屏蔽序列化和反序列化细节.


### Dapr与业务适配层(中间表示层).
```
// RPC适配.
import adapter "github.com/1046102779/go-sdk-wrapper"
s, err := adapter.NewService(logger.NewLogger("gateway.adapter"), ":8080")
_ = s.RegisterRPC("/tool/flow", srv.QueryFlow)

// 消息订阅适配.
import adapter "github.com/1046102779/go-sdk-wrapper"
s, err := adapter.NewService(logger.NewLogger("gateway.adapter"), ":8080")
_ = s.RegisterSubscribe(&adapter.Subscription{
     PubsubName: "push",
     Topic: "gbot-wxwork-push-prod",
}, srv.HandleMsg)

// 事件输入适配.
import adapter "github.com/1046102779/go-sdk-wrapper"
s, err := adapter.NewService(logger.NewLogger("gateway.adapter"), ":8080")
_ = s.RegisterInput("/timer/cron", srv.HandleCron)
```

### 业务代码
```golang
srv:=app.NewGatewayService(logger.NewLogger("gateway.robot"))

// QueryFlow 查询游戏工具流水.
func (g *GatewaySerivce) QueryFlow(ctx context.Context, 
    req *pb.ToolFlowReq) (rsp *pb.ErrRsp, err error) { 
     ctx = logger.NewContext(ctx, g.log)
     rsp = errs.NewErrRsp()
     defer func() { rsp = errs.HandleErrorToRsp(rsp, err) }()
     ...

     return
}

// HandleMsg 消费(kafka, pulsar...)消息.
func (g *GatewayService) HandleMsg(ctx context.Context, req *pb.PushReq)(err error) {
     ctx = logger.NewContext(ctx, g.log)
     ...

     return
}

// HandleCron 处理输入事件，比如：定时任务、钉钉消息、企微输入消息等
//
// 也可以是非CloudEvent标准的消息，比如：kafka、pulsar
func (g *GatewayService) HandleCron(ctx context.Context, req *pb.CronReq) (rsp *pb.ErrRsp, err error) {
     ctx = logger.NewContext(ctx, g.log)
     rsp = errs.NewErrRsp()
     defer func() { rsp = errs.HandleErrorToRsp(rsp, err) }()
     ...

     return
}
```

注意：**业务输入输出的协议，既可以是PB、也可以JSON.**

## 业务服务启动示例

```golang
import adapter "github.com/1046102779/go-sdk-wrapper"

s, err := adapter.NewService(logger.NewLogger("xxx.adapter"), ":8080")
s.RegisterRPC("/xxx/xxx/xx", service_handle_func)

if err= s.Start(); err !=nil{
	// TODO
}
```
