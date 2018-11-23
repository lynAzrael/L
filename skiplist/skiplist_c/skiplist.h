#ifndef _SKIP_LIST_H_
#define _SKIP_LIST_H_

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

#define bool int
#define true 1
#define false 0

#define MAX_LEVEL 16
#define MAX_VALUE_LENGTH 128

typedef struct nodeStruct node;
typedef struct skipList list;

struct nodeStruct {
	int key;
	int level;
	char value[128];
	node* forward[MAX_LEVEL];
};

struct skipList {
	int level;
	node* header;
};

node* CreateNode(int level, int key ,char value[MAX_VALUE_LENGTH]) {
	node* n = (node*)malloc(sizeof(node) + level * sizeof(node*));
	n->key = key;
	memcpy(n->value, value, MAX_VALUE_LENGTH - 1);
    // n->value = value;
	n->level = level;

	for (int i =0 ; i < level  ; i++) {
		n->forward[i] = NULL;
	}

	return n;
}

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

#endif
