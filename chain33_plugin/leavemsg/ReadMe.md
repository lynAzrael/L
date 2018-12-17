# 执行器LeaveMsg

## 1 整体结构
![](https://github.com/lynAzrael/L/blob/master/share/img/leavemsg.png)



## 2 代码分析
### 2.1 执行器的注册
在Plugin下的plugin/dapp/init/init.go文件中，需要将使用到的dapp引入。
![](https://github.com/lynAzrael/L/blob/master/share/img/init.png)

```go
import (
	...
	_ "github.com/33cn/plugin/plugin/dapp/leavemsg"
)
```

在新增执行器的目录下完成pluginmgr的注册
```go
func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.LeaveMsgX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
```

### 2.2 数据上链Exec
在执行Exec操作时，根据tx传入的funcname拼接获取到对应的Exec方法进行调用。
```go
func (d *DriverBase) Exec(tx *types.Transaction, index int) (receipt *types.Receipt, err error) {
	...
	name, value, err := d.ety.DecodePayloadValue(tx)
	if err != nil {
		return nil, err
	}
	funcmap := d.child.GetFuncMap()
	funcname := "Exec_" + name
	if _, ok := funcmap[funcname]; !ok {
		return nil, types.ErrActionNotSupport
	}
	valueret := funcmap[funcname].Func.Call([]reflect.Value{d.childValue, value, reflect.ValueOf(tx), reflect.ValueOf(index)})
	if !types.IsOK(valueret, 2) {
		return nil, types.ErrMethodReturnType
	}
	//参数1
	r1 := valueret[0].Interface()
	if r1 != nil {
		if r, ok := r1.(*types.Receipt); ok {
			receipt = r
		} else {
			return nil, types.ErrMethodReturnType
		}
	}
	//参数2
	r2 := valueret[1].Interface()
	err = nil
	if r2 != nil {
		if r, ok := r2.(error); ok {
			err = r
		} else {
			return nil, types.ErrMethodReturnType
		}
	}
	return receipt, err
}
```

因此在实现对应的方法时，也需要在前面加上Exec_的前缀
```go
func (l *LeaveMsg) Exec_XXX(xxx *xxxParam, tx *types.Transaction, index int) (*types.Receipt, error) {
	...
}
```
### 2.3 数据本地存储ExecLocal
执行ExecLocal操作将数据写入本地数据库时，与Exec类似需要根据tx传入的funcName拼接获取到对应的函数名
```go
func (d *DriverBase) callLocal(prefix string, tx *types.Transaction, receipt *types.ReceiptData, index int) (set *types.LocalDBSet, err error) {
	...
	name, value, err := d.ety.DecodePayloadValue(tx)
	if err != nil {
		return nil, err
	}
	//call action
	funcname := prefix + name
	funcmap := d.child.GetFuncMap()
	if _, ok := funcmap[funcname]; !ok {
		return nil, types.ErrActionNotSupport
	}
	valueret := funcmap[funcname].Func.Call([]reflect.Value{d.childValue, value, reflect.ValueOf(tx), reflect.ValueOf(receipt), reflect.ValueOf(index)})
	if !types.IsOK(valueret, 2) {
		return nil, types.ErrMethodReturnType
	}
	r1 := valueret[0].Interface()
	if r1 != nil {
		if r, ok := r1.(*types.LocalDBSet); ok {
			set = r
		} else {
			return nil, types.ErrMethodReturnType
		}
	}
	r2 := valueret[1].Interface()
	err = nil
	if r2 != nil {
		if r, ok := r2.(error); ok {
			err = r
		} else {
			return nil, types.ErrMethodReturnType
		}
	}
	return set, err
}
```

同样的，在实现ExecLocal函数时需要在funcName前面加上"ExecLocal_"前缀
```go
func (l *LeaveMsg) ExecLocal_xxxx(xxxx *xxxxParam, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	...
}

```
### 2.4 数据的查找

## 3 指令测试
### 3.1 创建交易
创建一个交易，其中消息的内容为"hello",发送方为14KEKbYtKKQm4wMthSK9J4La4nAiidGozt，接收方为15AsSzgynMSxhnK4gSMkFJgywKebKKUJky
```bash
curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.CreateTransaction","params":[{"execer":"leavemsg", "actionName":"Send", "payload":{"msg": "hello", "from":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt", "to":"15AsSzgynMSxhnK4gSMkFJgywKebKKUJky"}}] }'  -H 'content-type:text/plain;'   http://localhost:8801
```

响应:

```bash
{"id":1,"result":"0a086c656176656d7367121218010a0e0a0c57686f2061726520796f753f20e80730c1a3fcce8ba5a9973a3a2231456f79687135424b4b654b5a386f43646233484a673334694e42734d7162585245","error":null}
```

### 3.2 签名

```bash
 curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.SignRawTx","params":[{"addr":"14KEKbYtKKQm4wMthSK9J4La4nAiidGozt", "expire":"2h", "txHex":"0a086c656176656d7367125f18010a5b0a1148656c6c6f2c206d7920667269656e642e122231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a741a2231354173537a67796e4d5378686e4b3467534d6b464a6779774b65624b4b554a6b7920e807308f93f5f9d2d2bf81203a2231456f79687135424b4b654b5a386f43646233484a673334694e42734d7162585245"}] }' -H 'content-type:text/plain;' http://localhost:8801
```

响应:

```bash
{"id":1,"result":"0a086c656176656d7367120b10010a070a0568656c6c6f1a6e0801122102504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43871a473045022100c4ef5ab79cb1944eb7ff7cab5a735186b138dd7249d57884c7d8b5d9f93acee8022008ef35d6e8b6880dd6cd7cef130df9f4e98cf96c0ae0379340fa1c8b4e3f796e20e80728c7e494e00530dea2b8c6ce9cc983063a2231456f79687135424b4b654b5a386f43646233484a673334694e42734d7162585245","error":null}
```

### 3.3 发送交易

```bash
curl --data-binary '{"jsonrpc":"2.0", "id": 1, "method":"Chain33.SendTransaction","params":[{"data":"0a086c656176656d7367125f18010a5b0a1148656c6c6f2c206d7920667269656e642e122231344b454b6259744b4b516d34774d7468534b394a344c61346e41696964476f7a741a2231354173537a67796e4d5378686e4b3467534d6b464a6779774b65624b4b554a6b791a6d0801122102504fa1c28caaf1d5a20fefb87c50a49724ff401043420cb3ba271997eb5a43871a463044022012dd961831d23e49de5682e977921451ce4b4008f72259f7fbe7faed795d22e9022011964a20e4006206fc69864a5c6e5683cf2dc96bc5daa6581628cc69febe873a20e80728bdc099e00530f3e4ebe5b9d8a8c2263a2231456f79687135424b4b654b5a386f43646233484a673334694e42734d7162585245"}] }' -H 'content-type:text/plain;' http://localhost:8801
```

响应:

```bash
{"id":1,"result":"0x8a5aac307bfb65ca307c1b71ff567939186724d3d8a8cab2732104032be2d73b","error":null}
```