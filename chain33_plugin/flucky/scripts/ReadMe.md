# 游戏合约测试框架

<!-- TOC -->

- [游戏合约测试框架](#游戏合约测试框架)
    - [1 背景](#1-背景)
    - [2 配置文件](#2-配置文件)
        - [2.1 运行配置](#21-运行配置)
        - [2.2 动作配置](#22-动作配置)
        - [2.3 状态检查配置](#23-状态检查配置)
        - [2.4 公共状态信息配置](#24-公共状态信息配置)
    - [3 实现](#3-实现)
        - [3.1 执行ops预置条件](#31-执行ops预置条件)
        - [3.2 执行ops](#32-执行ops)
        - [3.3 op状态检查](#33-op状态检查)
        - [3.4 循环](#34-循环)
    - [4 测试](#4-测试)
        - [4.1 flucky测试配置](#41-flucky测试配置)

<!-- /TOC -->

## 1 背景
以Chain33为基础进行开发的游戏合约，进行稳定性和性能测试时需要使用脚本进行长时间的持续调用。如果每一个游戏合约都开发一套测试脚本，成本就有些高。因此考虑完成一个基本的游戏合约测试框架，在后续针对具体合约进行测试时，只需要根据游戏逻辑的不同，组装响应的测试指令即可完成测试。

## 2 配置文件
在配置文件中，填写游戏支持的操作信息以及对应chain33内部实现的rpc接口名称以及查询所需字段信息

### 2.1 运行配置

```toml
[Run]
preset=["SaveSeed,Unlock,ImportKey,CreateUser"]
implement=["Buy"]
runtimes=1
```

Run表示游戏执行的逻辑

|字段|说明|
|----|----|
|preset|预置条件，在运行之前需要执行|字符串数组|
|implement|具体执行逻辑，按照配置的先后顺序依次执行|
|runtimes|implement循环执行的次数|

### 2.2 动作配置

```toml
[op]
method="rpcMethodName"
param={"param1":"value1", "param2":"value2"...}
times=1
check="true"
```
此处的section名称op表示一个操作,在实际游戏合约可以是Start,Stop等操作。

动作的执行是通过调用合约已经写好的rpc接口进行实现。

|字段|类型|说明|备注|
|----|----|----|----|
|method|字符串|chain33内部已经实现的op动作对应rpc接口||
|param|json字符串|调用rpc时，需要传入的参数字段信息|当value字段为"inputParam"时，表示取值需要从外部传入|
|times|整形|操作重复地次数||
|check|字符串|执行钱是否需要进行状态检查||


### 2.3 状态检查配置

```toml
[op_Check]
interval=10
times=3
method="rpcMethodName"
param={"param1":"value1", "param2":"value2"...}
expectField=["field1", "updateTime"]
rule="$(expectField.remainTime+expectField.updateTime-CommonField.localtime)"
expectVal="0"
```

存在某种场景：只有当某个状态触发时，才允许后续的操作进行下去，因此需要一个状态检查

如果一个op操作的check选项为true时，则会根据对应的op_Check 中rule配置的规则进行校验

$表示需要进行数学运算，对后续()中的内容进行四则运算.

|字段|说明|
|----|----|
|interval|状态检查的间隔|
|times|状态检查的次数|
|method|字符串|chain33内部已经实现的状态信息查询接口||
|param|json字符串|调用rpc时，需要传入的参数字段信息|当value字段为"inputParam"时，表示取值需要从外部传入|
|expectField|期望从状态查询响应中获取到的字段，用于后续check逻辑的检查|
|rule|操作校验的逻辑|
|expectVal|校验的期望值，如果满足expectVal则表示校验成功|

### 2.4 公共状态信息配置

```toml
[CommonField]
locatime="common.GetLocalTime"
startheight="common.GetStartHeight"
currentheight="common.GetCurrentHeight"
starthash="Start.Resp['result']"
```
CommonField表示公共字段信息, 例如系统当前时间、当前区块高度等

其中各个变量的字段值根据配置的函数进行获取

## 3 实现
### 3.1 执行ops预置条件
1 从Run标签的preset中获取预置操作

2 校验是否支持（是否能在配置文件中查找到该操作对应的section信息，以及所需的method、param信息）

3 校验通过，则按照配置顺序依次执行;校验失败则退出

4 执行结束后，记录执行的状态信息; 成功则开始执行后续implement操作，失败则退出

### 3.2 执行ops
1 判断预置条件是否成功执行

2 获取implement中所有需要执行的操作，校验配置文件信息

3 校验通过，判断是否需要进行状态检查; 不通过则等待10s后再次校验，连续校验三次不通过则退出;

4 状态检查通过后，根据配置具体执行操作

>Note: implement操作执行顺序与配置顺序保持一致

### 3.3 op状态检查
1 根据op_Check中配置的策略查询状态信息

2 从响应中获取期望字段(expectField配置的字段列表), 并根据Check中的规则计算最终结果

3 将计算结果与预期值(expectVal)对比。相同则认为状态检查通过，可以继续进行后续操作；否则，认为检查失败，继续等待一次状态检查.

### 3.4 循环
1 当一轮操作(implement中的所有操作)均执行完毕之后，会再次从头开始执行，循环的次数由runtimes控制; preset不会循环执行，因为预置条件只需要执行一次即可

2 在一轮操作中，每个op也有可能执行多次。循环次数由各个操作对应的times控制

## 4 测试
### 4.1 flucky测试配置

```bash
# 根据游戏规则的不同，配置不同的策略
[Run]
preset=["SaveSeed","Unlock","ImportKey","CreateUser"]
implement=["Buy"]
runtimes=10

[Buy]
method="Chain33.CreateTransaction"
param={"execer": "flucky", "actionName": "Bet", "payload": {"amount": 10}}
times=5
needRange="false"
check="false"

[Buy_Check]
method="Chain33.Query"
param={"execer": "flucky", "funcName": "QueryLastRoundInfo", "payload": {}}
expectField=["remainTime"]
check=expectField.remainTime
symbol="lt"
expectVal=0

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
times=100

[CommonField]
localtime="GetLocalTime"
```
