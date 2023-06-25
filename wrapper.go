package wrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/1046102779/go-sdk-wrapper/bind"
	"github.com/dapr/go-sdk/service/common"
	"github.com/dapr/kit/logger"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// adapterDaprAPI dapr adapter layer.
type adapterDaprAPI struct {
	log          *logger.Logger
	decodeParser protojson.UnmarshalOptions
	encodeParser protojson.MarshalOptions
	pbType       reflect.Type
}

// newAdapterDaprAPI creates dapr adapter layer instance.
func newAdapterDaprAPI(log *logger.Logger) *adapterDaprAPI {
	return &adapterDaprAPI{
		log:          log,
		decodeParser: protojson.UnmarshalOptions{DiscardUnknown: true},
		encodeParser: protojson.MarshalOptions{UseEnumNumbers: true, UseProtoNames: true},
		pbType:       reflect.TypeOf((*proto.Message)(nil)).Elem(),
	}
}

// RPC adapter service invocation api.
func (a *adapterDaprAPI) RPC(fnVal reflect.Value) common.ServiceInvocationHandler {
	_, pbType, passed := a.checkAndGetHandleFunc(fnVal)
	if !passed {
		return nil
	}
	if fnVal.Type().NumOut() != 2 {
		return nil
	}
	return func(ctx context.Context, in *common.InvocationEvent) (
		out *common.Content, err error) {
		req := &In{
			queryString: in.QueryString,
			data:        in.Data,
			typ:         pbType,
		}
		if pbType.Implements(a.pbType) {
			req.isProtoJSON = true
		}
		var rsp *Out
		if rsp, err = a.call(ctx, req, fnVal); err != nil {
			return
		}

		if rsp == nil || len(rsp.vals) != 2 {
			return
		}

		out = &common.Content{ContentType: in.ContentType}
		if req.isProtoJSON {
			rsp, _ := rsp.vals[0].Interface().(proto.Message)
			out.Data, _ = a.encodeParser.Marshal(rsp)
		} else {
			out.Data, _ = json.Marshal(rsp.vals[0].Interface())
		}

		return
	}
}

// checkAndGetHandleFunc check and get function name, input payload.
func (a *adapterDaprAPI) checkAndGetHandleFunc(fnVal reflect.Value) (
	string, reflect.Type, bool) {
	if fnVal.IsZero() {
		return "", nil, false
	}
	if fnVal.Kind() != reflect.Func {
		return "", nil, false
	}
	pbType := fnVal.Type().In(1)
	if pbType.Kind() != reflect.Ptr {
		return "", nil, false
	}
	if fnVal.Type().NumIn() != 2 {
		return "", nil, false
	}

	return GetFuncName(fnVal), pbType, true
}

// Subscribe adapter TopicEvent(pubsub.sub) api.
func (a *adapterDaprAPI) Subscribe(fnVal reflect.Value) common.TopicEventHandler {
	_, pbType, passed := a.checkAndGetHandleFunc(fnVal)
	if !passed {
		return nil
	}
	if fnVal.Type().NumOut() != 1 {
		return nil
	}

	return func(ctx context.Context, in *common.TopicEvent) (bool, error) {
		var err error
		req := &In{
			data: in.RawData,
			typ:  pbType,
		}
		if pbType.Implements(a.pbType) {
			req.isProtoJSON = true
		}
		var rsp *Out
		if rsp, err = a.call(ctx, req, fnVal); err != nil {
			return false, nil
		}

		if rsp == nil || len(rsp.vals) != 1 {
			return false, nil
		}
		if !rsp.vals[0].IsZero() {
			err = fmt.Errorf("%v", rsp.vals[0].Interface())
			return false, nil
		}

		return false, nil
	}
}

// InputBinding adapter bindings.inputbinding api.
func (a *adapterDaprAPI) InputBinding(fnVal reflect.Value) common.BindingInvocationHandler {
	_, pbType, passed := a.checkAndGetHandleFunc(fnVal)
	if !passed {
		return nil
	}
	if fnVal.Type().NumOut() != 2 {
		return nil
	}

	return func(ctx context.Context, in *common.BindingEvent) (
		rspBody []byte, err error) {
		req := &In{
			data: in.Data,
			typ:  pbType,
		}
		if pbType.Implements(a.pbType) {
			req.isProtoJSON = true
		}
		var rsp *Out
		if rsp, err = a.call(ctx, req, fnVal); err != nil {
			return nil, nil
		}

		if rsp == nil || len(rsp.vals) != 2 {
			return nil, nil
		}

		if !rsp.vals[0].IsZero() {
			return nil, fmt.Errorf("%v", rsp.vals[0].Interface())
		}

		if req.isProtoJSON {
			tmpRsp, _ := rsp.vals[0].Interface().(proto.Message)
			rspBody, _ = a.encodeParser.Marshal(tmpRsp)
		} else {
			rspBody, _ = json.Marshal(rsp.vals[0].Interface())
		}

		return rspBody, nil
	}
}

// In input request.
type In struct {
	// queryString query string.
	queryString string
	// data requestbody.
	data []byte
	// payload MicroLogic Target Protocol Payload.
	payload reflect.Value
	// typ MicroLogic Target Protocol.
	typ reflect.Type
	// isProtoJSON PB?
	isProtoJSON bool
}

// Out Output response.
type Out struct {
	vals []reflect.Value
}

// call dynamics service handler API.
func (a *adapterDaprAPI) call(ctx context.Context, req *In, fnVal reflect.Value) (
	rsp *Out, err error) {
	req.payload = reflect.New(req.typ.Elem())
	if err = a.fillRequest(req); err != nil {
		return
	}
	rsp = new(Out)
	rsp.vals = fnVal.Call([]reflect.Value{reflect.ValueOf(ctx), req.payload})

	return
}

// fillRequest fill querystring and data to MicroLogic Protocols.
func (a *adapterDaprAPI) fillRequest(req *In) error {
	var err error
	if req.isProtoJSON {
		pbReq, _ := req.payload.Interface().(proto.Message)
		if err = a.decodeParser.Unmarshal(req.data, pbReq); err != nil {
			return err
		}
	} else {
		if req.queryString != "" {
			_ = bind.Bind(req.queryString, req.payload.Interface())
		}
		if len(req.data) > 0 {
			if err = json.Unmarshal(req.data, req.payload.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetFuncName dynamicly get specific function name.
func GetFuncName(fnVal reflect.Value) string {
	// such as:
	// github.com/.../$products/cmd/gateway/app.(*UserService).User-fm
	// github.com/.../$products/cmd/gateway/app.(*GatewayService).QueryFlow-fm
	//
	pkgFuncName := runtime.FuncForPC(fnVal.Pointer()).Name()
	index := strings.LastIndex(pkgFuncName, ".")
	if index < 0 {
		return ""
	}
	receiverFuncName := pkgFuncName[index+1:]

	return strings.TrimSuffix(receiverFuncName, "-fm")
}

// GetContentData gets content.data and avoid to panic.
func GetContentData(out *common.Content) []byte {
	if out == nil {
		return []byte{}
	}

	return out.Data
}
