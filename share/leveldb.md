# LevelDB 

## 1 leveldb整体结构

![](https://github.com/lynAzrael/L/blob/master/share/img/leveldb_architesture.png)

从图中可以看出，构成LevelDb静态结构的包括六个主要部分：内存中的MemTable和Immutable MemTable以及磁盘上的几种主要文件：Current文件，Manifest文件，log文件以及SSTable文件。当然，LevelDb除了这六个主要部分还有一些辅助的文件，但是以上六个文件和数据结构是LevelDb的主体构成元素。

## 1.1 设计思想
LSM Tree(Log-structured merge tree): 结构将数据的修改增量地保存在内存中，并在达到一定大小之后批量写到磁盘上。

这种结构是将更新的内容写入到一系列的更小的所有文件中。每个文件包含了一批在一段短时间内的变化，并在被写入前都被排序以便于稍微快速的检索。这些文件都是不变的；他们从来不更新。每次更新都写入新的文件。会检查所有文件来定期合并，以降低文件的数量。

LSM本身由MemTable,Immutable MemTable,SSTable等多个部分组成，其中MemTable在内存，用于记录最近修改的数据，一般用跳跃表来组织。当MemTable达到一定大小后，将其冻结起来变成Immutable MemTable，然后开辟一个新的MemTable用来记录新的记录。而Immutable MemTable则等待转存到磁盘。

### 1.1.1 Mem Table


#### 1.1.1.2 Compaction操作(minor compaction)
当Immutable MemTable中的数据达到了一定的大小之后，会将内容写到磁盘中。

![minor compaction](https://github.com/lynAzrael/L/blob/master/share/img/leveldb_minor_compaction.png)

immutable memtable其实是一个SkipList，其中的记录是根据key有序排列的，遍历key并依次写入一个level 0 的新建SSTable文件中，写完后建立文件的index 数据，这样就完成了一次minor compaction。从上图中也可以看到，删除操作并不是真正的删除记录(因为想要获取到这个key值对应的kv数据，需要比较复杂的查找，所以在minor compaction中不会真正的删除某个记录)，而是将delete操作作为一条记录写入到文件中。

```go
go db.mCompaction()
```

### 1.1.3 SS Table
#### 2.3.1 Compaction操作(major compaction)

```go
go db.tCompaction()
```

## 2 Operation

### 2.1 PUt
levelDb的更新操作速度是非常快的，源于其内部机制决定了这种更新操作的简单性。 

![put operation](https://github.com/lynAzrael/L/blob/master/share/img/leveldb_put.png)

上图是levelDb如何更新KV数据的示意图，从图中可以看出，对于一个插入操作Put(Key,Value)来说，完成插入操作包含两个具体步骤：首先是将这条KV记录以顺序写的方式追加到之前介绍过的log文件末尾，因为尽管这是一个磁盘读写操作，但是文件的顺序追加写入效率是很高的，所以并不会导致写入速度的降低；第二个步骤是:如果写入log文件成功，那么将这条KV记录插入内存中的Memtable中，前面介绍过，Memtable只是一层封装，其内部其实是一个Key有序的SkipList列表，插入一条新记录的过程也很简单，即先查找合适的插入位置，然后修改相应的链接指针将新记录插入即可。完成这一步，写入记录就算完成了，所以一个插入记录操作涉及一次磁盘文件追加写和内存SkipList插入操作，这是为何levelDb写入速度如此高效的根本原因。

### 2.2 Get
leveldb读取数据的过程如下图所示：

![get operation](https://github.com/lynAzrael/L/blob/master/share/img/leveldb_get.png)

LevelDb首先会去查看内存中的Memtable，如果Memtable中包含key及其对应的value，则返回value值即可；如果在Memtable没有读到key，则接下来到同样处于内存中的Immutable Memtable中去读取，类似地，如果读到就返回，若是没有读到,那么只能万般无奈下从磁盘中的大量SSTable文件中查找。因为SSTable数量较多，而且分成多个Level，所以在SSTable中读数据是相当蜿蜒曲折的一段旅程。总的读取原则是这样的：首先从属于level 0的文件中查找，如果找到则返回对应的value值，如果没有找到那么到level 1中的文件中去找，如此循环往复，直到在某层SSTable文件中找到这个key对应的value为止（或者查到最高level，查找失败，说明整个系统中不存在这个Key)。

在这个过程中的查找路径是从Memtable到Immutable Memtable，再从Immutable Memtable到文件，而文件中则是从低level到高level。之所以选择这么个查询路径，是因为从信息的更新时间来说，很明显Memtable存储的是最新鲜的KV对；Immutable Memtable中存储的KV数据对的新鲜程度次之；而所有SSTable文件中的KV数据新鲜程度一定不如内存中的Memtable和Immutable Memtable的。对于SSTable文件来说，如果同时在level L和Level L+1找到同一个key，level L的信息一定比level L+1的要新。也就是说，上面列出的查找路径就是按照数据新鲜程度排列出来的，越新鲜的越先查找。

SSTable文件很多，如何快速地找到key对应的value值？在LevelDb中，level 0一直都爱搞特殊化，在level 0和其它level中查找某个key的过程是不一样的。因为level 0下的不同文件可能key的范围有重叠，某个要查询的key有可能多个文件都包含，这样的话LevelDb的策略是先找出level 0中哪些文件包含这个key（manifest文件中记载了level和对应的文件及文件里key的范围信息，LevelDb在内存中保留这种映射表）， 之后按照文件的新鲜程度排序，新的文件排在前面，之后依次查找，读出key对应的value。而如果是非level 0的话，因为这个level的文件之间key是不重叠的，所以只从一个文件就可以找到key对应的value。

最后一个问题,如果给定一个要查询的key和某个key range包含这个key的SSTable文件，那么levelDb是如何进行具体查找过程的呢？levelDb一般会先在内存中的Cache中查找是否包含这个文件的缓存记录，如果包含，则从缓存中读取；如果不包含，则打开SSTable文件，同时将这个文件的索引部分加载到内存中并放入Cache中。 这样Cache里面就有了这个SSTable的缓存项，但是只有索引部分在内存中，之后levelDb根据索引可以定位到哪个内容Block会包含这条key，从文件中读出这个Block的内容，在根据记录一一比较，如果找到则返回结果，如果没有找到，那么说明这个level的SSTable文件并不包含这个key，所以到下一级别的SSTable中去查找。

### 2.3 Compaction

## 3 源码剖析

### 3.1 基础组件

#### 3.1.1 Bloom Filter

#### 3.1.2 CRC32
CRC（Cyclic Redundancy Check）中文名是循环冗余校验，在数据存储和数据通讯领域，为了保证数据的正确，采用检错的手段。


#### 3.1.3 Murmur Hash

### 3.2 SkipList
如果一个基点存在k个向前的指针的话，则将该节点称之为k层的节点。一个跳表的层MaxLevel义为跳表中所有节点中最大的层数。

下面给出一个完整的跳表的图示：

![skip_list](https://github.com/lynAzrael/L/blob/master/share/img/skiplist_linklist_complete.png)
### 3.3 Mem Table
memtable提供了插入、删除以及查询kv的操作，但是实际上Memtable并不存在真正的删除操作，删除某个Key的Value在Memtable内是作为插入一条记录实施的，但是会打上一个Key的删除标记，真正的删除操作会在之后的Compaction过程中去掉这个KV。

LevelDb的Memtable中KV对是根据Key大小有序存储的，在系统插入新的KV时，LevelDb要把这个KV插到合适的位置上以保持这种Key有序性。其实，LevelDb的Memtable类只是一个接口类，真正的操作是通过背后的SkipList来做的，包括插入操作和读取操作等，所以Memtable的核心数据结构是一个SkipList。

SkipList不仅是维护有序数据的一个简单实现，而且相比较平衡树来说，在插入数据的时候可以避免频繁的树节点调整操作，所以写入效率是很高的，LevelDb整体而言是个高写入系统，SkipList在其中应该也起到了很重要的作用。



### 3.4 SS Table
### 3.5 日志文件

### 3.2 Current文件\Manifest文件\版本信息



###

## 附录

### SkipList的实现

skiplist结构

```c
typedef struct skipList list;

struct skipList {
	int level;
	node* header;
};
``` 

node结构

```c
typedef struct nodeStruct node;

struct nodeStruct {
	int key;
	int level;
	char value[128];
	node* forward[MAX_LEVEL];
};
```

skiplist创建

```c
list* CreateList() {
	list* sl = (list*)malloc(sizeof(list));
	sl->level = MAX_LEVEL - 1;

	node* n = (node*)malloc(sizeof(node) + MAX_LEVEL * sizeof(node*));
	sl->header = n;

	for (int i = 0; i < MAX_LEVEL - 1; i++) {
		sl->header->forward[i] = NULL;
	}
	return sl;
}
```

插入节点
```c
bool insertNode(list* l, int key, int level) {
	// 根据key值，创建一个node
	node* n = (node*)malloc(sizeof(node) + level * sizeof(node*));
	n->key = key;
	n->level = level;

	for (int i=0; i < level; i ++) {
		n->forward[i] = NULL;
	}

	// 如果节点的level超过了链表的最大level， 则返回异常
	if (n->level > l->level ) {
		return false;
	}
	
	for (int i = 0; i < n->level; i++) {
		node* header = l->header;
		node** current = &(l->header->forward[i]);
		if (*current == NULL) {
			*current = n;
		} else {
			while (*current != NULL && (*current)->key < n->key) {
				current = &((*current)->forward[i]);
			} 
			if (*current == NULL){
				*current = n;
			} else if ((*current)->key > n->key) {
				n->forward[i] = *current;
				*current = n;
			}
		}
	}
	return true;
}
```

删除节点
```c
bool deleteNode(list* l, int key) {
	node * n;
	for (int i = 0; i < l->level; i++) {
		node* header = l->header;
		node** current = &(l->header->forward[i]);

		while (*current != NULL) {
			if ((*current)->key < key) {
				// skiplist是已经排过序的
				current = &((*current)->forward[i]);				
			} else if ((*current)->key == key) {
				// 进行删除操作
				n = *current;
				*current = (*current)->forward[i];
				break;
			} else {
				// 如果当前节点的key超过了需要删除的key，则说明要删除的节点不存在.
				return true;
			}
		}
	}
	free(n);
}
```

遍历list
```c
void showList(list* l){
	for (int i = MAX_LEVEL - 1; i >= 0 ;i--){
		printf("In level %d...\n", i);
		node* n = l->header->forward[i];
		printf("header");
		while (n != NULL){
			printf("\t->\t%d", n->key);
			n = n->forward[i];
		}
		printf("\n");
	}
}
```