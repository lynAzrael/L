/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package rpc

import (
	"context"
	"encoding/hex"
	"github.com/33cn/chain33/types"
	ptypes "github.com/33cn/plugin/plugin/dapp/f3d/ptypes"
)

// F3DStartTx 创建游戏开始交易
func (c *Jrpc) F3DStartTx(parm *ptypes.GameStartReq, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}
	start := &ptypes.F3DStart{
		Round: parm.Round,
	}
	reply, err := c.cli.Start(context.Background(), start)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *Jrpc) F3DLuckyDrawTx(parm *ptypes.GameDrawReq, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}

	luckyDraw := &ptypes.F3DLuckyDraw{
		Round: parm.Round,
	}

	reply, err := c.cli.LuckyDraw(context.Background(), luckyDraw)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *Jrpc) F3DBuyKeysTx(parm *ptypes.GameBuyKeysReq, result *interface{}) error {
	if parm == nil {
		return types.ErrInvalidParam
	}

	buyKey := &ptypes.F3DBuyKey{
		KeyNum: parm.Num,
	}

	reply, err := c.cli.BuyKeys(context.Background(), buyKey)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil

}
