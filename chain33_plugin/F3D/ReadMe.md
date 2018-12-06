#F3D游戏合约设计

## 1 流程设计

### 1.1 创建游戏
![创建游戏流程](https://github.com/lynAzrael/L/blob/master/chain33_plugin/F3D/resources/StartRound.png)

#### 1.1.1 首轮游戏的创建
第一轮游戏是由管理员创建，可以通过调用StartGameRound接口实现。

>1.首轮是否需要在奖池放入一定额度的bty作为奖金？
>
>2.配置信息是否只在首轮传入，一旦游戏启动，所有相关配置信息不可变？
>

#### 1.1.2 后续游戏的创建 
后续的游戏均是在开奖完毕之后，调用StartGameRound接口进行创建。


### 1.2 停止游戏
![停止游戏(开奖)流程](https://github.com/lynAzrael/L/blob/master/chain33_plugin/F3D/resources/EndRound.png)

停止游戏即进行开奖操作。

通过GetGameRounds获取轮次信息，从lastOwner中获取中奖人的地址信息进行奖金的发放，如果该字段为空，则说明这一轮游戏没有玩家购买key。需要将一轮游戏的时间缩短，并将奖池中的全部奖金移入下一轮游戏。

>开奖规则：如果检测到上一区块的打包时间和本轮游戏中最后一个Key交易区块的打包时间间隔大于1小时，则触发开奖流程

### 1.3 玩家购票
![玩家购票流程](https://github.com/lynAzrael/L/blob/master/chain33_plugin/F3D/resources/BuyKeys.png)

### 1.4 信息查询


## 2 接口设计

### 2.1 游戏创建

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
#### 2.3.1 GetGameRounds 获取轮次信息
请求结构:

```bash
{
    "startRound":int64,
    "endRound":int64
}
```

|参数|类型|说明|
|----|----|----|
|startRound|int64|开始轮次|
|endRound|int64|结束轮次|

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

### 3.2 LocalDB存储
|键|值|用途|说明|
|-|-|-|-|
|LODB-f3d-round-info:{round}|RoundInfo|保存每一个轮次的游戏信息|

## 4 配置
合约配置信息：

```toml
[Addr]
# 本游戏合约管理员地址
manager = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
# 本游戏合约奖金池地址
bonuspool = "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"

[Bonus]
# 超级大奖分成百分比
winner = 40
# 参与者分红百分比
keyowner = 30
# 滚动到下期奖金池百分比
pool = 20
# 平台运营及开发者费用百分比
platform = 10

[Game]
# 本游戏一轮运行的最长周期（单位：秒）
life = 3600
# 一把钥匙延长的游戏时间（单位：秒）
key = 30
# 一次购买钥匙最多延长的游戏时间（单位：秒）
maxkey = 300
# 没有玩家购买key时，一轮游戏缩短的时间（单位：秒）
nouser = 30
# 每次钥匙价格上浮的百分比
increase = 10
```
