package executor

import (
	"fmt"
	"github.com/33cn/chain33/types"
)
import ftypes "f3d/types"

func (f *F3DGame) Exec_Start(start *ftypes.StartGame, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(f, tx, index)
	return action.GameStart(start)
}

func (f *F3DGame) Exec_Stop(stop *ftypes.StopGame, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(f, tx, index)
	return action.GameStop(stop)
}

func (f *F3DGame) Exec_BuyKey(buy *ftypes.BuyKeys, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(f, tx, index)
	return action.BuyKeys(buy)
}

func getKeyNumberPrefix(round int64, addr string) []byte {
	return []byte(fmt.Sprintf("mavl-f3d-user-keys:{%d}:{%s}", round, addr))
}

func getRoundsStartPrefix() []byte {
	return []byte(fmt.Sprintf("mavl-f3d-round-start"))
}

func getRoundsStopPrefix() []byte {
	return []byte(fmt.Sprintf("mavl-f3d-round-end"))
}

func getRoundsCurrentPrefix() []byte {
	return []byte(fmt.Sprintf("mavl-f3d-last-round"))
}

func getKeyCurrentPricePrefix(round int64) []byte {
	return []byte(fmt.Sprintf("mavl-f3d-key-price:{%d}", round))
}
