# F3D游戏合约设计

- [F3D游戏合约设计](#f3d游戏合约设计)
    - [1 流程设计](#1-流程设计)
        - [1.1 创建游戏](#11-创建游戏)
            - [1.1.1 首轮游戏的创建](#111-首轮游戏的创建)
            - [1.1.2 后续游戏的创建](#112-后续游戏的创建)
        - [1.2 停止游戏](#12-停止游戏)
        - [1.3 玩家购票](#13-玩家购票)
        - [1.4 信息查询](#14-信息查询)
    - [2 接口设计](#2-接口设计)
        - [2.1 游戏创建](#21-游戏创建)
        - [2.2 游戏停止(开奖)](#22-游戏停止开奖)
        - [2.3 游戏相关信息查询](#23-游戏相关信息查询)
            - [2.3.1 GetGameRounds 获取轮次信息](#231-getgamerounds-获取轮次信息)
            - [2.3.2](#232)
        - [2.4 玩家相关信息查询](#24-玩家相关信息查询)
            - [2.4.1 GetUserBonus 获取玩家中奖信息](#241-getuserbonus-获取玩家中奖信息)
            - [2.4.2 GetUserIsLastOne 获取玩家是否中奖](#242-getuserislastone-获取玩家是否中奖)
            - [2.4.3 GetUserKeyNum 获取玩家持有key的数量](#243-getuserkeynum-获取玩家持有key的数量)
    - [3 存储](#3-存储)
        - [3.1 StateDB存储](#31-statedb存储)
        - [3.2 LocalDB存储](#32-localdb存储)
    - [4 配置](#4-配置)

## 1 流程设计

### 1.1 创建游戏
![创建游戏流程](https://github.com/lynAzrael/L/blob/master/chain33_plugin/F3D/resources/StartRound.png)

#### 1.1.1 首轮游戏的创建
第一轮游戏是由管理员创建，通过调用StartGameRound接口实现。

> Q1: 首轮是否需要在奖池放入一定额度的bty作为奖金？


> Q2:配置信息是否只在首轮传入，一旦游戏启动，所有相关配置信息不可变？


#### 1.1.2 后续游戏的创建 
后续的游戏均是在开奖完毕之后，调用StartGameRound接口进行创建。


### 1.2 停止游戏
![停止游戏(开奖)流程](https://github.com/lynAzrael/L/blob/master/chain33_plugin/F3D/resources/EndRound.png)

停止游戏即进行开奖操作，开奖操作一旦触发就不允许玩家继续购票。

通过GetGameRounds获取轮次信息，从lastOwner中获取中奖人的地址信息进行奖金的发放。如果该字段为空则说明这一轮游戏没有玩家购买key，需要将一轮游戏的时间缩短并将奖池中的全部奖金移入下一轮游戏。否则将锁定本轮游戏，调用合约进行奖金下发，并开始新一轮游戏。


>开奖规则：如果检测到上一区块的打包时间和本轮游戏中最后一个Key交易区块的打包时间间隔大于1小时，则触发开奖流程

>检测时间间隔：

>|剩余时间|下一次检测周期|
>|>10分钟|10分钟|
>|>5分钟|1分钟|
>|>1分钟|30秒|
>|<1分钟|10秒|

>剩余时间： 如果能够获取到本轮游戏的上一个区块，则剩余时间为本地时间与上一个区块打包时间的间隔；否则为本地时间与本轮游戏开始时间的间隔。

> Q3: 剩余时间是指？ 

整个过程涉及三个操作：
1、查询剩余时间
通过查询轮次信息，获取到当前轮次的剩余时间，根据上表中的逻辑进行不间断地检测
2、判断开奖

3、启动下一轮游戏


### 1.3 玩家购票
![玩家购票流程](https://github.com/lynAzrael/L/blob/master/chain33_plugin/F3D/resources/BuyKeys.png)
 
当一轮游戏开始时，玩家可以进行购票。
如果玩家购买的key数量过多导致余额不足以支付时，则认为此次交易失败。
在进行开奖的过程中，玩家不可以继续购票。

是否还能购买Keys?(本轮游戏已经结束或者正在开奖)


> Q4: 如何判断是最后一个key？

> A: 收到购买钥匙的交易时，将购买时间与上一个区块的打包时间进行比较。如果间隔不小于1小时，则更新区块链上最后一个key的购买时间。

> Q6 : 

> A


### 1.4 信息查询
#### 1.4.1 轮次信息查询
#### 1.4.2 用户信息查询

## 2 接口设计

### 2.1 游戏创建
```go
func (f *F3DGame) Exec_Start(start *ftypes.StartGame, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(f, tx, index)
	return action.GameStart(start)
}
```
### 2.2 游戏停止(开奖)

### 2.3 游戏相关信息查询
轮次信息

```proto
message RoundInfo {
    // 游戏轮次
    int64 round = 1;
    // 本轮游戏开始事件
    int64 beginTime  = 2;
    // 本轮游戏结束时间
    int64 endTime    = 3;
    // 本轮游戏目前为止最后一把钥匙持有人（游戏开奖时，这个就是中奖人）
    string lastOwner   = 4;
    // 最后一把钥匙的购买时间
    int64 lastKeyTime = 5;
    // 最后一把钥匙的价格
    int64 lastKeyPrice = 6;
    // 本轮游戏奖金池总额
    int64 bonusPool  = 7;
    // 本轮游戏参与地址数
    int64 userCount = 8;
    // 本轮游戏募集到的key个数
    int64 keyCount = 9;
}

```

用户所购Key的相关信息

```proto
message KeyInfo{
	// 游戏轮次  (是由系统合约填写后存储）
    int64 round = 1;
      
    // 本次购买key的价格 (是由系统合约填写后存储）
    int64 keyPrice = 2;  
    
    // 用户本次买的key的数量
    int64 keyNum = 3;
    
    // 用户地址 (是由系统合约填写后存储）
    string addr = 4;
    
}
```
#### 2.3.1 GetGameRounds 获取轮次信息
请求结构:

```bash
{
    "startRound":int64,
    "endRound":int64
}
```

|参数|类型|是否必填|说明|
|----|----|----|----|
|startRound|int64|是|开始轮次,<=0时表示当前轮次,否则为指定轮次|
|endRound|int64|是|结束轮次,<=0时表示当前轮次,否则为指定轮次|

>如果startRound <=0, endRound为之前的轮次，需要返回从之前的轮次到当前轮次的所有信息么？

响应:
```bash
{
    [
        {
            "round":1,
            "beginTime":1543994216,
            ...
        },
        {
            "round":2,
            "beginTime":1543995384,
            ...
        },
        ...
    ]
}
```

响应为RoundInfo结构体构成的数组。

#### 2.3.2 

### 2.4 玩家相关信息查询

#### 2.4.1 GetUserBonus 获取玩家中奖信息
请求结构：

```bash
{
    "User":string,
    "StartRound":int64,
    "EndRound":int64
}
```
|参数|类型|说明|
|----|----|----|
|User|string|游戏参与者地址|
|StartRound|int64|开始轮次|
|EndRound|int64|结束轮次|


响应:
```bash
{
    [
        {
            "round":int64,
            "keyNum":int64,
            "bonus":int64
        }
    ]
}
```

|参数|类型|说明|
|----|----|----|
|round|int64|游戏轮次|
|keyNum|int64|持有的钥匙把数|
|bonus|int64|拿到的奖金|


#### 2.4.2 GetUserIsLastOne 获取玩家是否中奖


#### 2.4.3 GetUserKeyNum 获取玩家持有key的数量


## 3 存储
### 3.1 StateDB存储
|键|值|用途|说明|
|-|-|-|-|
|mavl-f3d-user-keys:{round}:{addr}|钥匙数量|记录用户在某一轮游戏中所持有的钥匙数量|参数为用户的地址以及需要查询轮次|
|mavl-f3d-round-start|轮次开始|存储开始的轮次信息|
|mavl-f3d-round-end|轮次结束|存储结束的轮次信息|
|mavl-f3d-last-round|当前轮次|保存当前轮次||
|mavl-f3d-key-price:{round}|钥匙价格|存储最后一把钥匙的价格|
|mavl-f3d-manager-key|管理员地址|存储f3d游戏管理员的地址信息|

### 3.2 LocalDB存储
|键|值|用途|说明|
|-|-|-|-|
|LODB-f3d-round-info:{round}|RoundInfo|保存每一个轮次的游戏信息|

## 4 配置
合约配置信息：

```toml
# 本游戏合约管理员地址
manager = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
# 本游戏合约平台开发者分成地址
developer = "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"

# 超级大奖分成百分比
winnerbonus = 0.4
# 参与者分红百分比
keybonus = 0.3
# 滚动到下期奖金池百分比
poolbonus = 0.2
# 平台运营及开发者费用百分比
developerbonus = 0.1

# 本游戏一轮运行的最长周期（单位：秒）
lifetime = 3600
# 一把钥匙延长的游戏时间（单位：秒）
keytime = 30
# 一次购买钥匙最多延长的游戏时间（单位：秒）
maxkeyIncrTime = 300
# 一轮游戏没有玩家购买key时，缩减的游戏时间
nouserDecrTime = 30

# 钥匙涨价幅度（下一个人购买钥匙时在上一把钥匙基础上浮动幅度百分比），范围1-100
startKeyPrice= 0.1
# 钥匙起始价格
incrKeyPrice = 0.1
```

## 5 cmd指令


