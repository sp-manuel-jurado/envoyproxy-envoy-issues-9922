package main_test

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/socialpoint/envoyproxy-envoy-issues-9922/services/ping/pkg/sp_rpc"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestServiceServer_Ping(t *testing.T) {
	testCases := []struct{
		version string
		addr    string
	}{
		{
			version: "envoy:v1.12.0",
			addr:    "envoy-v1-12-0:10000",
		},
		{
			version: "envoy:v1.12.1",
			addr:    "envoy-v1-12-1:10000",
		},
		{
			version: "envoy:v1.12.2",
			addr:    "envoy-v1-12-2:10000",
		},
		{
			version: "envoy:v1.12.3",
			addr:    "envoy-v1-12-3:10000",
		},
		{
			version: "envoy:v1.13.0",
			addr:    "envoy-v1-13-0:10000",
		},
		{
			version: "envoy:v1.13.1",
			addr:    "envoy-v1-13-1:10000",
		},
		{
			version: "envoy:v1.14.1",
			addr:    "envoy-v1-14-1:10000",
		},
		{
			version: "envoy-dev:latest",
			addr:    "envoy-dev-latest:10000",
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		msg := fmt.Sprintf("Testing envoy version: %s", tc.version)
		t.Run(msg, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)

			conn, err := grpc.Dial(tc.addr, grpc.WithInsecure())
			r.NoError(err)

			cli := pb.NewPingServiceClient(conn)
			_, err = cli.Ping(context.Background(), &empty.Empty{})
			a.NoError(err)
		})
	}
}
