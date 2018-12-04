package rpc

import (
	rpctypes "github.com/33cn/chain33/rpc/types"

)

// 对外提供服务的RPC接口总体定义
type Jrpc struct {
	cli *channelClient
}

// RPC接口的本地实现
type channelClient struct {
	rpctypes.ChannelClient
}

func Init(name string, s rpctypes.RPCServer) {
	cli := &channelClient{}
	// 为了简单起见，这里只注册Jrpc，如果提供grpc的话也在这里注册
	cli.Init(name, s, &Jrpc{cli: cli}, nil)
}
