#F3D游戏合约设计

## 1 流程设计

### 1.1 创建游戏
#### 1.1.1 首轮游戏的创建
第一轮游戏是由管理员创建，可以通过调用StartGameRound接口实现。

>首轮是否需要在奖池放入一定额度的bty作为奖金？

#### 1.1.2 后续游戏的创建 
后续的游戏均是在开奖完毕之后，自动调用StartGameRound接口进行创建。


### 1.2 停止游戏
停止游戏即进行开奖操作。

通过GetGameRounds获取轮次信息，通过lastOwner来获取中奖人的地址信息进行奖金的发放，如果该字段为空，则说明这一轮游戏没有玩家购买key。需要将一轮游戏的时间缩短，并将奖池中的全部奖金移入下一轮游戏。

### 1.3 玩家购票


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

```json
{
    "startRound":int64,
    "endRound":int64
}
```

响应:
```json
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


#### 2.3.2 

### 2.4 玩家相关信息查询

#### 2.4.1 GetUserBonus 获取玩家中奖信息
请求结构：

```json
{
	"User":string,
	"StartRound":int64,
	"EndRound":int64
}
```

响应:

#### 2.4.2 GetUserIsLastOne 获取玩家是否中奖

#### 2.4.3 GetUserKeyNum 获取玩家持有key的数量


## 3 存储

## 4 配置
合约配置信息：

```bash
# 本游戏合约管理员地址
f3d.manager = 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt

# 本游戏合约奖金池地址
f3d.bonuspool = 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv

# 超级大奖分成百分比
f3d.bonus.winner = 40

# 参与者分红百分比
f3d.bonus.key = 30

# 滚动到下期奖金池百分比
f3d.bonus.pool = 20

# 平台运营及开发者费用百分比
f3d.bonus.platform = 10

# 本游戏一轮运行的最长周期（单位：秒）
f3d.time.life = 3600

# 一把钥匙延长的游戏时间（单位：秒）
f3d.time.key = 30

# 一次购买钥匙最多延长的游戏时间（单位：秒）
f3d.time.maxkey = 300

# 没有玩家购买key时，一轮游戏缩短的时间（单位：秒）
f3d.time.nouser = 30

# 每次钥匙价格上浮的百分比
f3d.key.increase = 10
```
