# 平行链环境搭建

<!-- TOC -->

- [1.1 手工编译](#11-手工编译)
- [1.2 使用官方版本](#12-使用官方版本)

<!-- /TOC -->

## 1.1 手工编译
切换到如下路径：${GoPath}/src/github.com/33cn

GoPath为安装Go编译器之后的环境变量，如果未安装Go，请参考环境准备中的[Go环境安装](https://chain.33.cn/document/81#1.1%20Go%20%E7%8E%AF%E5%A2%83%E5%AE%89%E8%A3%85)

除GoPath之外的其他子目录如果不存在，请手工创建。

获取最新Chain33代码分支进行编译，github地址：git@github.com:33cn/plugin.git

linux下编译

windos下编译



## 1.2 使用官方版本

获取文件
```bash 
wget https://bty.oss-ap-southeast-1.aliyuncs.com/chain33/paraChain.tar.gz
```

解压
```bash
tar zxvf paraChain.tar.gz
```

修改配置文件

```bash
vi paraChain/chain33.para.toml
```


ParaRemoteGrpcClient项取值为："101.37.227.226:8802,39.97.20.242:8802,47.107.15.126:8802,jiedian2.33.cn"

mainnetJrpcAddr项取值为："http://jiedian1.33.cn:8801"


修改后如下所示：

![paraRemoteGrpcClient](https://github.com/lynAzrael/L/blob/master/share/img/paraRemoteGrpcClient.png)

![mainnetJrpcAddr](https://github.com/lynAzrael/L/blob/master/share/img/mainnetJrpcAddr.png)

启动chain33进程

```bash
cd paraChain && ./chain33 -f chain33.para.toml
```
