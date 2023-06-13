package wrapper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRegisterRPC 测试RPC注册.
func TestRegisterRPC(t *testing.T) {
	s, _ := NewService(nil, ":8080")
	uts := []struct {
		Router      string
		HandlerFunc interface{}
		Expected    bool
	}{
		{
			Router:      "/router",
			HandlerFunc: nil,
			Expected:    false,
		},
		{
			Router:      "/router",
			HandlerFunc: struct{}{},
			Expected:    false,
		},
		{
			Router: "router",
			HandlerFunc: func(ctx context.Context, req *struct {
				ID int `json:"id"`
			}) (rsp *struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}, err error) {
				return nil, nil
			},
			Expected: true,
		},
	}
	for _, ut := range uts {
		err := s.RegisterRPC(ut.Router, ut.HandlerFunc)
		if ut.Expected {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
