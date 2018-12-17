/*
 * Copyright Fuzamei Corp. 2018 All Rights Reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package commands

import (
	"fmt"
	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	ptypes "github.com/33cn/plugin/plugin/dapp/f3d/ptypes"
	"github.com/spf13/cobra"
	"time"
)

func F3DCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "f3d",
		Short: "F3D contracts operation",
		Args:  cobra.MinimumNArgs(1),
	}

	cmd.AddCommand(
		F3DGameCmd(),
		F3DInfoCmd(),
	)

	return cmd
}

func F3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "game",
		Short: "game operation for f3d.",
	}
	cmd.AddCommand(
		createF3DGameCmd(),
		luckyDrawF3DGameCmd(),
		buyKeyF3DGameCmd(),
	)
	return cmd
}

func F3DInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "info operation for f3d.",
	}
	cmd.AddCommand(
		recordInfoCmd(),
		roundInfoCmd(),
		//roundsInfoCmd(),
		lastRoundInfoCmd(),
		addrInfoCmd(),
	)
	return cmd
}

func addF3DGameFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("round", "r", 0, "Game Round")
	cmd.MarkFlagRequired("round")
}

func createF3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a new f3d game.",
		Run:   createNewGame,
	}
	addF3DGameFlags(cmd)
	return cmd
}

func createNewGame(cmd *cobra.Command, args []string) {
	round, _ := cmd.Flags().GetInt64("round")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	if round > 0 {
		start(round, rpcLaddr)
	}
}

func start(round int64, rpcLaddr string) {
	// 开启新一轮
	params := ptypes.GameStartReq{Round: round}

	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DStartTx", params, &res)
	ctx.RunWithoutMarshal()
}

func luckyDrawF3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "draw",
		Short: "Send the bonus to the last user that buys the key",
		Run:   runLuckyDraw,
	}
	addF3DGameLuckyDrawFlags(cmd)
	return cmd
}

func addF3DGameLuckyDrawFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("round", "r", 0, "Game Round")
	cmd.Flags().BoolP("force", "f", false, "force lucky draw")
	cmd.MarkFlagRequired("round")
}

func runLuckyDraw(cmd *cobra.Command, args []string) {
	round, _ := cmd.Flags().GetInt64("round")
	force, _ := cmd.Flags().GetBool("force")
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	if force {
		// 直接开奖
		luckDraw(round, rpcLaddr)
	} else {
		// 检测是否到了开奖时间
		//go func() {
		var interval time.Duration
		for {
			roundInfo := getLastRoundInfo(rpcLaddr)
			if !roundCheck(round, roundInfo) {
				return
			}
			if remainTimeCheck(roundInfo, &interval) {
				time.Sleep(time.Duration(ptypes.GetF3dTimeDelay()) * time.Second)
				luckDraw(round, rpcLaddr)
				break
			}
			fmt.Println("It 's not time to luckydraw , interval:", interval*time.Second, ", now:", time.Now().Unix())
			time.Sleep(interval * time.Second)
			continue
		}
		//}()
	}
	// 开奖之后，创建新一轮的游戏
	time.Sleep(time.Duration(ptypes.GetF3dTimeDelay()) * time.Second)
	start(round+1, rpcLaddr)
}

func luckDraw(round int64, rpcLaddr string) {
	params := ptypes.GameDrawReq{Round: round}

	var res string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DLuckyDrawTx", params, &res)
	ctx.RunWithoutMarshal()
}

func getLastRoundInfo(rpcAddr string) *ptypes.RoundInfo {
	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	params.FuncName = ptypes.FuncNameQueryLastRoundInfo
	params.Payload = types.MustPBToJSON(&ptypes.QueryF3DLastRound{})
	resp = &ptypes.RoundInfo{}

	ctx := jsonclient.NewRPCCtx(rpcAddr, "Chain33.Query", params, resp)
	ctx.Run()

	return resp.(*ptypes.RoundInfo)
}

func roundCheck(round int64, info *ptypes.RoundInfo) bool {
	if round > 0 && round == info.Round {
		return true
	}
	return false
}

func remainTimeCheck(info *ptypes.RoundInfo, interval *time.Duration) bool {
	currentTime := time.Now().Unix()
	remainTime := info.RemainTime + info.UpdateTime - currentTime
	fmt.Println("Remain time is :", remainTime)

	if remainTime > 600 {
		*interval = 600
	} else if remainTime > 300 {
		*interval = 60
	} else if remainTime > 60 {
		*interval = 30
	} else if remainTime > 0 {
		*interval = 10
	} else {
		return true
	}

	return false
}

func buyKeyF3DGameCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy",
		Short: "Buy some keys",
		Run:   buyKeys,
	}
	addNumberFlags(cmd)
	return cmd
}

func buyKeys(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	num, _ := cmd.Flags().GetInt64("num")

	params := ptypes.GameBuyKeysReq{
		Num: num,
	}

	var resp string
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "f3d.F3DBuyKeysTx", params, &resp)
	ctx.RunWithoutMarshal()
}

func addNumberFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("num", "n", 1, "the number of keys that you want to buy.")
}

func recordInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "Show the buy record about the user matched by address and round.",
		Run:   buyRecordInfoQuery,
	}
	addRecordInfoQueryFlags(cmd)
	return cmd
}

func addRecordInfoQueryFlags(cmd *cobra.Command) {
	cmd.Flags().Int64P("round", "r", 0, "Game Round")
	cmd.Flags().Int64P("index", "i", 0, "index")
	cmd.Flags().StringP("addr", "a", "", "user addr")
	cmd.MarkFlagRequired("round")
}

func roundInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "round",
		Short: "Show the round info matched by round.",
		Run:   roundInfoQuery,
	}
	addF3DGameFlags(cmd)
	return cmd
}

//func roundsInfoCmd() *cobra.Command {
//	cmd := &cobra.Command{
//		Use:   "rounds",
//		Short: "Show the round info matched by round.",
//		Run:   roundsInfoQuery,
//	}
//	addRoundsInfoQueryFlag(cmd)
//	return cmd
//}

//func addRoundsInfoQueryFlag(cmd *cobra.Command) {
//	cmd.Flags().Int64P("startRound", "s", 0, "start round")
//	cmd.Flags().Int32P("direction", "d", 1, "query direction, 0: desc  1:asc")
//	cmd.Flags().Int32P("count", "c", 0, "query amount")
//
//	cmd.MarkFlagRequired("startRound")
//}

func lastRoundInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last_round",
		Short: "Show the last round info.",
		Run:   lastRoundInfoQuery,
	}

	return cmd
}

func addrInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addr info",
		Short: "Show addr info.",
		Run:   addrInfoQuery,
	}

	addAddrInfoQueryFlags(cmd)
	return cmd
}

func addAddrInfoQueryFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("addr", "a", "", "addr info")
	cmd.Flags().Int64P("round", "r", 0, "round info")

	cmd.MarkFlagRequired("addr")
	cmd.MarkFlagRequired("round")
}

func roundInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	round, _ := cmd.Flags().GetInt64("round")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryF3DByRound{
		Round: round,
	}
	fmt.Println(req)

	params.FuncName = ptypes.FuncNameQueryRoundInfoByRound
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.RoundInfo{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}

//func roundsInfoQuery(cmd *cobra.Command, args []string) {
//	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
//	startRound, _ := cmd.Flags().GetInt64("startRound")
//	direction, _ := cmd.Flags().GetInt32("direction")
//	count, _ := cmd.Flags().GetInt32("count")
//
//	var params rpctypes.Query4Jrpc
//	var resp interface{}
//
//	params.Execer = ptypes.F3DX
//	req := ptypes.QueryF3DListByRound{
//		StartRound: startRound,
//		Direction:  direction,
//		Count:      count,
//	}
//	params.FuncName = ptypes.FuncNameQueryRoundsInfoByRounds
//	params.Payload = types.MustPBToJSON(&req)
//	resp = &ptypes.ReplyF3DList{}
//
//	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
//	ctx.Run()
//}

func lastRoundInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryF3DLastRound{}
	params.FuncName = ptypes.FuncNameQueryLastRoundInfo
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.RoundInfo{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}

func buyRecordInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	round, _ := cmd.Flags().GetInt64("round")
	addr, _ := cmd.Flags().GetString("addr")
	index, _ := cmd.Flags().GetInt64("index")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryBuyRecordByRoundAndAddr{
		Round: round,
		Addr:  addr,
		Index: index,
	}
	params.FuncName = ptypes.FuncNameQueryBuyRecordByRoundAndAddr
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.ReplyBuyRecord{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}

func addrInfoQuery(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	round, _ := cmd.Flags().GetInt64("round")
	addr, _ := cmd.Flags().GetString("addr")

	var params rpctypes.Query4Jrpc
	var resp interface{}

	params.Execer = ptypes.F3DX
	req := ptypes.QueryKeyCountByRoundAndAddr{
		Round: round,
		Addr:  addr,
	}
	params.FuncName = ptypes.FuncNameQueryKeyCountByRoundAndAddr
	params.Payload = types.MustPBToJSON(&req)
	resp = &ptypes.AddrInfo{}

	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Chain33.Query", params, resp)
	ctx.Run()
}
