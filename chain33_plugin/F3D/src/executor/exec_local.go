package executor

import (
	"fmt"
	ftypes "f3d/types"
	"github.com/33cn/chain33/types"
)

func (f *F3DGame) execLocal(receipt *types.Receipt) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receipt.GetTy() != types.ExecOk {
		return dbSet, nil
	}
	for _, log := range receipt.Logs {
		switch log.Ty {
		// 根据传入的日志类型，判断是否需要更新数据库中的轮次信息
		case ftypes.LogTypeF3DGameBuyKeys, ftypes.LogTypeF3DGameStart, ftypes.LogTypeF3DGameStop:

		}
	}
	return dbSet, nil
}

func (f *F3DGame) ExecLocal_Start(game *ftypes.StartGame, transaction types.Transaction, receipt *types.Receipt, index int) (*types.LocalDBSet, error) {
	return f.execLocal(receipt)
}

func (f *F3DGame) ExecLocal_Stop(stop *ftypes.StopGame, tx types.Transaction, receipt *types.Receipt, index int) (*types.LocalDBSet, error) {
	return f.execLocal(receipt)
}

func (f *F3DGame) ExecLocal_Buy(buy *ftypes.BuyKeys, tx types.Transaction, receipt *types.Receipt, index int) (*types.LocalDBSet, error) {
	return f.execLocal(receipt)
}

func getLocalDbRoundsInfoPrefix(round int64) []byte {
	return []byte(fmt.Sprint("LODB-f3d-round-info:{%d}", round))
}
