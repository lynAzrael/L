#游戏合约相关测试框架

## 1 背景
以Chain33为基础进行开发的游戏合约，进行稳定性和性能测试时需要使用脚本进行长时间的持续调用。如果每一个游戏合约都开发一套测试脚本，成本就有些高。因此考虑完成一个基本的游戏合约测试框架，在后续针对具体合约进行测试时，只需要根据游戏逻辑的不同，组装响应的测试指令即可完成测试。

## 2 配置文件
在配置文件中，填写游戏支持的操作信息以及对应chain33内部实现的rpc接口名称以及查询所需字段信息
### 2.1 op声明
```bash=
# 根据游戏规则的不同，配置不同的策略
[Op]
ops="CreateUser,Start,Buy,Stop"
```

ops表示游戏支持的动作.
其中CreateUser，Start，Buy，Stop表示游戏支持创建用户、开始、购买、停止的操作。
后续需要针对每一个操作，配置一个Section。

### 2.2 op创建
```bash=
[op]
method="rpcMethodName"
param={"param1":"value1", "param2":"value2"...}
times="repeadtime"
needrange="true"
check="true"
```
此处的section名称op表示一个操作,实际可以为Start,Stop操作。
动作的执行，是通过调用合约已经写好的rpc接口进行实现。
|字段|说明|
|method|合约内部已经实现的rpc接口名称|
|param|调用合约需要传入的参数字段信息|
|times|操作重复的次数|
|needrange|操作是否需要用户传入一个范围信息|
|check|操作执行前是否需要校验|

### 2.3 op合法性校验
```bash=
[op_Check]
method="rpcMethodName"
param={"param1":"value1", "param2":"value2"...}
expectField=["field1", "updateTime"]
check=$(expectField.remainTime+expectField.updateTime-CommonField.localtime)
expectVal="0"
```

如果一个op操作的check选项为true时，则会根据对应的op_Check 中配置的规则进行校验
|字段|说明|
|expectField|期望从状态查询响应中获取到的字段，用于后续check逻辑的检查|
|check|操作校验的逻辑|
|expectVal|校验的期望值，如果满足expectVal则表示校验成功。|