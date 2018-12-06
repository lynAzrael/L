package executor

import (
	ftypes "f3d/types"
	"github.com/33cn/chain33/account"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
)

type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	localDB      dbm.Lister
	index        int
}

func NewAction(f *F3DGame, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &Action{f.GetCoinsAccount(), f.GetStateDB(), hash, fromaddr,
		f.GetBlockTime(), f.GetHeight(), dapp.ExecAddress(string(tx.Execer)), f.GetLocalDB(), index}
}

func (action *Action) GameStart(start *ftypes.StartGame) (*types.Receipt, error) {
	//  对交易的发起人进行校验

	//  获取上一轮游戏信息
	lastRound := start.LastRound

	// 上一轮游戏信息校验

	// 判断是否为第一轮游戏
	if lastRound <= 0 {
		// 如果是第一轮游戏

	} else {
		// 如果不是第一轮游戏

	}

	// 校验成功，则开始新一轮游戏

	// 游戏信息上链

	return nil, nil
}

func (action *Action) GameStop(stop *ftypes.StopGame) (*types.Receipt, error) {
	// 对交易的发起人进行校验

	//  获取当前轮次信息

	// 轮次信息校验

	// 结束当前轮次，并调用合约进行奖金的下发，同时开启新一轮游戏

	// 更新游戏轮次信息

	return nil, nil
}

func (action *Action) BuyKeys(buy *ftypes.BuyKeys) (*types.Receipt, error) {
	// 轮次信息进行校验

	// 获取当前轮次和价格

	// 检查余额是否充足

	// 更新用户持有的key相关信息

	// 更新key的价格

	// 更新轮次信息

	return nil, nil
}
