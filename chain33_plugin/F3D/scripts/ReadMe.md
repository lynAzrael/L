#游戏合约相关测试框架

<!-- TOC -->

- [游戏合约相关测试框架](#游戏合约相关测试框架)
    - [1 背景](#1-背景)
    - [2 配置文件](#2-配置文件)
        - [2.1 操作声明](#21-操作声明)
        - [2.2 操作运行](#22-操作运行)
        - [2.3 操作创建](#23-操作创建)
        - [2.4 状态信息检查](#24-状态信息检查)
        - [2.5 公共状态信息](#25-公共状态信息)
    - [3 实现](#3-实现)
        - [3.1 执行ops预置条件](#31-执行ops预置条件)
        - [3.2 执行ops](#32-执行ops)
        - [3.3 op状态检查](#33-op状态检查)
        - [3.4 循环](#34-循环)
    - [4 测试](#4-测试)
        - [4.1 F3D测试配置](#41-f3d测试配置)

<!-- /TOC -->

## 1 背景
以Chain33为基础进行开发的游戏合约，进行稳定性和性能测试时需要使用脚本进行长时间的持续调用。如果每一个游戏合约都开发一套测试脚本，成本就有些高。因此考虑完成一个基本的游戏合约测试框架，在后续针对具体合约进行测试时，只需要根据游戏逻辑的不同，组装响应的测试指令即可完成测试。

## 2 配置文件
在配置文件中，填写游戏支持的操作信息以及对应chain33内部实现的rpc接口名称以及查询所需字段信息
### 2.1 操作声明
```bash=
# 根据游戏规则的不同，配置不同的策略
[Op]
ops="CreateUser,Start,Buy,Stop"
```

ops表示游戏支持的动作
其中CreateUser，Start，Buy，Stop表示游戏支持创建用户、开始、购买、停止的操作。
后续需要针对每一个操作，配置一个Section。

### 2.2 操作运行
```bash=
[Run]
preset="CreateUser"
implement="Start,Buy,Stop"
runtimes="10"
```

Run表示游戏执行的逻辑

|字段|说明|
|----|----|
|preset|预置条件，在运行之前需要执行|
|implemen|具体执行逻辑，按照配置的先后顺序依次执行|

### 2.3 操作创建
```bash=
[op]
method="rpcMethodName"
param={"param1":"value1", "param2":"value2"...}
times="repeadtime"
needrange="true"
check="true"
```
此处的section名称op表示一个操作,在实际游戏合约可以是Start,Stop等操作。
动作的执行是通过调用合约已经写好的rpc接口进行实现。

|字段|说明|
|----|----|
|method|合约内部已经实现的rpc接口名称|
|param|调用合约需要传入的参数字段信息|
|times|操作重复的次数|
|needrange|操作是否需要用户传入一个范围信息|
|check|操作执行前是否需要状态检查|

### 2.4 状态信息检查
```bash=
[op_Check]
method="rpcMethodName"
param={"param1":"value1", "param2":"value2"...}
expectField=["field1", "updateTime"]
check=$(expectField.remainTime+expectField.updateTime-CommonField.localtime)
expectVal="0"
```

存在某种场景：只有当某个状态触发时，才允许后续的操作进行下去，因此需要一个状态检查。

如果一个op操作的check选项为true时，则会根据对应的op_Check 中配置的规则进行校验

|字段|说明|
|----|----|
|expectField|期望从状态查询响应中获取到的字段，用于后续check逻辑的检查|
|check|操作校验的逻辑|
|expectVal|校验的期望值，如果满足expectVal则表示校验成功。|

### 2.5 公共状态信息
```bash=
[CommonField]
locatime="common.GetLocalTime"
```
CommonField表示公共字段信息, 例如系统当前时间、系统Cpu数量等。

其中各个变量的字段值根据配置的函数进行获取。

## 3 实现
### 3.1 执行ops预置条件
1 从Run标签的preset中查找需要执行的预置op

2 校验是否支持该op，即是否能在配置文件中查找到该op对应的section信息

3 如果支持该操作，则执行该op

### 3.2 执行ops
1 判断预置条件是否执行完毕

2 从implement中获取所有需要执行的op，校验是否支持

3 校验通过，则根据配置策略进行状态检查以及具体执行(执行顺序与implement中配置顺序保持一致)

### 3.3 op状态检查
1 从op中获取对应section(${op}_Check)中check字段的取值，如果为true则需要进行状态的检查

2 根据op_Check中配置的策略获取状态信息，如果计算结果与预期值相符，则认为检查成功，可以进行op操作.

### 3.4 循环
1 当一轮操作(implement中的所有op)均执行完毕之后，会再次从头开始执行，循环的次数由runtimes控制; preset不会循环执行，因为预置条件只需要执行一次即可。

2 在一轮操作中，每个op也有可能执行多次。循环次数由各个op对应的times控制。


## 4 测试
### 4.1 F3D测试配置
```bash=
# 根据游戏规则的不同，配置不同的策略
[Run]
preset=["CreateUser"]
implement=["Start,Buy,Stop"]
runtimes="10"

[Start]
method="f3d.F3DStartTx"
param={"round": "inputParam"}
check="true"

[Start_Check]
method="Chain33.Query"
param={"execer": "f3d", "funcName": "QueryLastRoundInfo", "payload": ""}
expectField=["remainTime"]
check=$(expectField.remainTime)
expectVal=""

[Stop]
method="f3d.F3DLuckyDrawTx"
param={"round": "inputParam"}
check="true"

[Stop_Check]
method= "Chain33.Query"
param={"execer": "f3d", "funcName": "QueryLastRoundInfo", "payload": ""}
expectField=["remainTime", "updateTime"]
check=$(expectField.remainTime+expectField.updateTime-CommonField.localtime)
expectVal="0"

[Buy]
method="f3d.F3DBuyKeysTx"
param={"num": "inputParam"}
times="1000"
needRange="true"
check="true"

[Buy_Check]
method="Chain33.Query"
param={"execer": "f3d", "funcName": "QueryLastRoundInfo", "payload": ""}
expectField=["remainTime"]
check=$(expectField.remainTime)
expectVal=""

# Chain33交易的签名、发送等操作，一般不需要修改
[SaveSeed]
method="Chain33.SaveSeed"
param={"seed": "lens involve pudding midnight climb depend alcohol sibling carpet ghost garment child faith upper runway", "passwd": "fzm"}

[ImportKey]
method="Chain33.ImportPrivkey"
param={"privkey": "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944", "label": "manager"}

[Unlock]
method="Chain33.UnLock"
params={"passwd": "fzm"}

[Sign]
method="Chain33.SignRawTx"
param={"addr": "inputParam", "txHex": "inputParam",  "expire": "0"}

[Send]
method="Chain33.SendTransaction"
param={"token": "BTY",  "data": "inputParam"}

[CreateUser]
method="Chain33.NewAccount"
param={"label": "inputParam"}
needRange="true"

[CommonField]
locatime
```
