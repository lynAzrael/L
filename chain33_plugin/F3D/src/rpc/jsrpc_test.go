// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"strings"
	"testing"

	commonlog "github.com/33cn/chain33/common/log"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util/testnode"
	pty "github.com/33cn/plugin/plugin/dapp/f3d/ptypes"
	"github.com/stretchr/testify/assert"

	_ "github.com/33cn/chain33/system"
	_ "github.com/33cn/plugin/plugin"
)

func init() {
	commonlog.SetLogLevel("error")
}

func TestJRPCChannel(t *testing.T) {
	// 启动RPCmocker
	mocker := testnode.New("--notset--", nil)
	defer func() {
		mocker.Close()
	}()
	mocker.Listen()

	jrpcClient := mocker.GetJSONC()
	assert.NotNil(t, jrpcClient)

	testCases := []struct {
		fn func(*testing.T, *jsonclient.JSONClient) error
	}{
		{fn: testGameCreateRawTxCmd},
		{fn: testGameLuckyDrawTxCmd},
		{fn: testGameBuyKeysTxCmd},
		{fn: testRoundInfoCmd},
		{fn: testRoundListCmd},
	}
	for index, testCase := range testCases {
		err := testCase.fn(t, jrpcClient)
		if err == nil {
			continue
		}
		assert.NotEqualf(t, err, types.ErrActionNotSupport, "test index %d", index)
		if strings.Contains(err.Error(), "rpc: can't find") {
			assert.FailNowf(t, err.Error(), "test index %d", index)
		}
	}
}

func testGameCreateRawTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &pty.GameStartReq{}
	var res string
	return jrpc.Call("f3d.F3DStartTx", params, &res)
}

func testGameLuckyDrawTxCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	params := &pty.GameDrawReq{}
	var res string
	return jrpc.Call("f3d.F3DLuckyDrawTx", params, &res)
}

func testGameBuyKeysTxCmd(t *testing.T, jrpc *jsonclient.JSONClient)error {
	params := &pty.GameBuyKeysReq{}
	var res string
	return jrpc.Call("f3d.F3DBuyKeysTx", params, &res)
}

func testRoundInfoCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.QueryF3DByRound{}
	params.FuncName = pty.FuncNameQueryRoundInfoByRound
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplyF3D{}
	return jrpc.Call("Chain33.Query", params, rep)
}

func testRoundListCmd(t *testing.T, jrpc *jsonclient.JSONClient) error {
	var rep interface{}
	var params rpctypes.Query4Jrpc
	req := &pty.QueryF3DListByRound{}
	params.FuncName = pty.FuncNameQueryRoundsInfoByRounds
	params.Payload = types.MustPBToJSON(req)
	rep = &pty.ReplyF3DList{}

	return jrpc.Call("Chain33.Query", params, rep)
}
