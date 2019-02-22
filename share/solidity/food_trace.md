# 食品溯源案例

<!-- TOC -->autoauto- [食品溯源案例](#食品溯源案例)auto    - [1 合约的创建](#1-合约的创建)auto        - [1.1 合约的编译](#11-合约的编译)auto            - [1.1.1 Remix-ide](#111-remix-ide)auto                - [1.1.1.1 使用Ethereum官方提供的在线编译器](#1111-使用ethereum官方提供的在线编译器)auto            - [1.1.2 Intellij-Solidity](#112-intellij-solidity)auto                - [1.1.2.1 插件安装](#1121-插件安装)auto                - [1.1.2.2 solc安装](#1122-solc安装)auto                - [1.1.2.3 编译](#1123-编译)auto        - [1.2 合约的创建](#12-合约的创建)auto    - [2 合约的调用](#2-合约的调用)auto        - [2.1 信息录入](#21-信息录入)auto            - [2.1.1 猪肉信息录入](#211-猪肉信息录入)auto            - [2.1.2 信息检查](#212-信息检查)auto            - [2.1.2 食品信息录入](#212-食品信息录入)auto            - [2.1.3 质检信息录入](#213-质检信息录入)auto            - [2.1.4 上架信息录入](#214-上架信息录入)auto        - [2.2 信息查询](#22-信息查询)auto        - [2.2.1 商品信息查询](#221-商品信息查询)autoauto<!-- /TOC -->

## 1 合约的创建

### 1.1 合约的编译
使用solidity编译器进行合约的编译，可使用以下工具：

Remix-ide
Intellij-Solidity

#### 1.1.1 Remix-ide
##### 1.1.1.1 使用Ethereum官方提供的在线编译器

http://remix.ethereum.org/

> 由于Chrome浏览器不支持编译后abi、bin文件的拷贝，因此尽量使用Edge或者IE浏览器。

#### 1.1.2 Intellij-Solidity

##### 1.1.2.1 插件安装
打开IntelliJ IDEA， 在File->Settings->Plugins选项卡中，查找IntelliJ-Solidity插件进行安装。

[IntelliJ_Solidity插件查找](https://github.com/lynAzrael/L/tree/master/share/img/intelliJ_solidity.png)

##### 1.1.2.2 solc安装
如果本地已有geth节点，可以直接使用节点自带solc，

##### 1.1.2.3 编译
使用Build->Compile Solidity编译合约，编译后的结果可以在项目栏中看到。

[Compile_Solidity合约编译](https://github.com/lynAzrael/L/tree/master/share/img/compile_solidity.png)

使用生成的abi、bin文件进行合约的创建

### 1.2 合约的创建
使用chain33-cli中已有命令行进行合约的创建
./chain33-cli evm create 

## 2 合约的调用

### 2.1 信息录入

|猪肉信息|
|----|
|出栏批次|
|名称|
|重量|
|出栏日期|
|产地|

|商品信息|
|----|
|食品编号|
|体积|
|名称|
|重量|
|生产日期|
|包装日期|
|保质期|
|对应猪肉批次|
|上架日期|

#### 2.1.1 猪肉信息录入

|录入信息|
|----|
|出栏批次|
|名称|
|重量|
|出栏日期|
|产地|

使用chain33-cli进行信息录入
```bash

```

#### 2.1.2 信息检查
使用getPigNumber()函数查看信息是否成功录入

./chain33-cli evm call -c "" -f "" -e "evm合约地址" -b "abi" 

```bash
./chain33-cli evm call -c "" -f "" -e "evm合约地址" -b 
```

#### 2.1.2 食品信息录入
|录入信息||
|----|----|
|食品编号|
|名称|
|重量|
|体积|
|生产日期|
|包装日期|
|保质期|

使用addFoodInfo()进行食品信息的录入
```bash

```

#### 2.1.3 质检信息录入

#### 2.1.4 上架信息录入


### 2.2 信息查询

### 2.2.1 商品信息查询