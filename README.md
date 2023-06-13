## go-sdk-wrapper

Reduce business developers understanding of the dapr standardized API, and only focus on their own payload. Mask serialization and deserialization details.


### Dapr(go-sdk) and business adaptation layer (intermediate presentation layer)

```
// RPC adapter.
import adapter "github.com/1046102779/go-sdk-wrapper"
s, err := adapter.NewService(logger.NewLogger("gateway.adapter"), ":8080")
_ = s.RegisterRPC("/tool/flow", srv.QueryFlow)

// pubsub.sub adapter.
import adapter "github.com/1046102779/go-sdk-wrapper"
s, err := adapter.NewService(logger.NewLogger("gateway.adapter"), ":8080")
_ = s.RegisterSubscribe(&adapter.Subscription{
     PubsubName: "push",
     Topic: "gbot-wxwork-push-prod",
}, srv.HandleMsg)

// bindings.inputbinding adapter.
import adapter "github.com/1046102779/go-sdk-wrapper"
s, err := adapter.NewService(logger.NewLogger("gateway.adapter"), ":8080")
_ = s.RegisterInput("/timer/cron", srv.HandleCron)
```

### MicroLogic Codes
```golang
srv:=app.NewGatewayService(logger.NewLogger("gateway.robot"))

// QueryFlow querys game flow.
func (g *GatewaySerivce) QueryFlow(ctx context.Context, 
    req *pb.ToolFlowReq) (rsp *pb.ErrRsp, err error) { 
     ctx = logger.NewContext(ctx, g.log)
     rsp = errs.NewErrRsp()
     defer func() { rsp = errs.HandleErrorToRsp(rsp, err) }()
     ...

     return
}

// HandleMsg consumes (kafka, pulsar...) payload.
func (g *GatewayService) HandleMsg(ctx context.Context, req *pb.PushReq)(err error) {
     ctx = logger.NewContext(ctx, g.log)
     ...

     return
}

// HandleCron handle inputbindings event. such as：cronjob、dingtalk、wxwork...
func (g *GatewayService) HandleCron(ctx context.Context, req *pb.CronReq) (rsp *pb.ErrRsp, err error) {
     ctx = logger.NewContext(ctx, g.log)
     rsp = errs.NewErrRsp()
     defer func() { rsp = errs.HandleErrorToRsp(rsp, err) }()
     ...

     return
}
```

Notice：**The protocol of business input and output, which can be either Protobuffer or JSON.**

## MicroLogic Demo

```golang
import adapter "github.com/1046102779/go-sdk-wrapper"

s, err := adapter.NewService(logger.NewLogger("xxx.adapter"), ":8080")
s.RegisterRPC("/xxx/xxx/xx", serviceHandleFunc)

if err= s.Start(); err !=nil{
	// TODO
}
```
