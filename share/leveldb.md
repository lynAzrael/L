# LevelDB 

## 1 leveldb整体结构

### 1.1 LSM（Log-structured merge tree)
这种结构是将更新的内容写入到一系列的更小的所有文件中。每个文件包含了一批在一段短时间内的变化，并在被写入前都被排序以便于稍微快速的检索。这些文件都是不变的；他们从来不更新。每次更新都写入新的文件。会检查所有文件来定期合并，以降低文件的数量。



## 2 leveldb组成

### 2.1 memtable
memtable就是一个SkipList的具体实现



### 2.2 immutable memtable

#### 2.2.1 SkipList
如果一个基点存在k个向前的指针的话，那么陈该节点是k层的节点。

一个跳表的层MaxLevel义为跳表中所有节点中最大的层数。

下面给出一个完整的跳表的图示：

![](https://github.com/lynAzrael/L/blob/master/share/img/skiplist_linklist_complete.png)


#### 2.2.2 Compaction操作(minor compaction)

```go
go db.mCompaction()
```

### 2.3 SS Table
#### 2.3.1 Compaction操作(major compaction)

```go
go db.tCompaction()
```

### 2.4 Current

## 附录

### SkipList的实现

#### 结构定义

```go
	
``` 
