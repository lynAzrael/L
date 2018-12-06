package executor

import (
	ftypes "f3d/types"
	driver "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
)

var (

	driverName = ftypes.F3DX
)

// 初始化时通过反射获取本执行器的方法列表
func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&F3DGame{}))
}

func Init(name string, sub []byte) {
	driver.Register(driverName, NewF3DGame, types.GetDappFork(driverName, "Enable"))
}

type F3DGame struct {
	driver.DriverBase
}

func NewF3DGame() driver.Driver {
	t := &F3DGame{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

func (l *F3DGame) GetDriverName() string {
	return driverName
}

func GetName() string {
	return NewF3DGame().GetName()
}

