package hello_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	hellov1 "github.com/henrywhitaker3/connect-template/gen/hello/v1"
	"github.com/henrywhitaker3/connect-template/internal/connect/common"
	"github.com/henrywhitaker3/connect-template/internal/test"
	"github.com/stretchr/testify/require"
)

func TestItSaysHello(t *testing.T) {
	app, cancel := test.App(t)
	defer cancel()

	client, cancel := test.HelloClient(t, app)
	defer cancel()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	type testCase struct {
		name string
		req  *connect.Request[hellov1.HelloRequest]
		err  error
	}

	tcs := []testCase{
		{
			name: "responds to authenticated request",
			req: test.Authenticate(t, app, connect.NewRequest(&hellov1.HelloRequest{
				Name: test.Word(),
			})),
			err: nil,
		},
		{
			name: "errors to unauthenticated request",
			req: connect.NewRequest(&hellov1.HelloRequest{
				Name: test.Word(),
			}),
			err: common.ErrUnauthenticated,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			resp, err := client.Hello(ctx, c.req)
			if c.err == nil {
				require.Nil(t, err)
				require.Equal(t, fmt.Sprintf("Hello %s!", c.req.Msg.Name), resp.Msg.Message)
				return
			}

			test.ErrorIs(t, err, c.err)
		})
	}
}
