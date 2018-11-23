#include "skiplist.h"

int main() { 
	list* l = CreateList();
	(void)insertNode(l, 1, 2);
	(void)insertNode(l, 2, 3);
	(void)insertNode(l, 3, 5);
	(void)insertNode(l, 5, 7);
	(void)insertNode(l, 6, 10);
	(void)insertNode(l, 9, 2);
	(void)insertNode(l, 4, 3);
	(void)insertNode(l, 7, 1);

	showList(l);

	(void)deleteNode(l, 6);

	showList(l);
	
	return 0;
}
