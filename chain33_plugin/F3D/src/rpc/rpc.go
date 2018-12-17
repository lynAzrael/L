/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package rpc

import (
	"github.com/33cn/chain33/types"
	ptypes "github.com/33cn/plugin/plugin/dapp/f3d/ptypes"
)

import (
	"context"
)

func (c *channelClient) Start(ctx context.Context, start *ptypes.F3DStart) (*types.UnsignTx, error) {
	val := &ptypes.F3DAction{
		Ty:    ptypes.F3dActionStart,
		Value: &ptypes.F3DAction_Start{Start: start},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	data, err := types.FormatTxEncode(ptypes.F3DX, tx)
	if err != nil {
		return nil, err
	}
	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) LuckyDraw(ctx context.Context, luckyDraw *ptypes.F3DLuckyDraw) (*types.UnsignTx, error) {
	val := &ptypes.F3DAction{
		Ty:    ptypes.F3dActionDraw,
		Value: &ptypes.F3DAction_Draw{Draw: luckyDraw,},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	data, err := types.FormatTxEncode(ptypes.F3DX, tx)
	if err != nil {
		return nil, err
	}

	return &types.UnsignTx{Data: data}, nil
}

func (c *channelClient) BuyKeys(ctx context.Context, buyKeys *ptypes.F3DBuyKey) (*types.UnsignTx, error) {
	val := &ptypes.F3DAction{
		Ty:    ptypes.F3dActionBuy,
		Value: &ptypes.F3DAction_Buy{Buy: buyKeys,},
	}
	tx := &types.Transaction{
		Payload: types.Encode(val),
	}
	data, err := types.FormatTxEncode(ptypes.F3DX, tx)
	if err != nil {
		return nil, err
	}

	return &types.UnsignTx{Data: data}, nil
}
