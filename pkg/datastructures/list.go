package datastructures

// ListNode represents single node in doubly linked list.
type ListNode struct {
	val  string
	prev *ListNode
	next *ListNode
}

// LinkedList is implementation of doubly linked list.
type LinkedList struct {
	length int
	head   *ListNode
	tail   *ListNode
}

// NewLinkedList is a constructor for LinkedList.
func NewLinkedList() *LinkedList {
	return &LinkedList{
		length: 0,
		head:   nil,
		tail:   nil,
	}
}

// get returns pointer to node located at the specified index.
// It MUST BE CALLED after validation that node at index exists.
func (ll *LinkedList) get(index int) *ListNode {
	var curr *ListNode

	if index < ll.length/2 {
		curr = ll.head
		for range index {
			curr = curr.next
		}
	} else {
		curr = ll.tail
		for range ll.length - index - 1 {
			curr = curr.prev
		}
	}

	return curr
}

// Len is getter method for linked list length.
func (ll *LinkedList) Len() int {
	return ll.length
}

// Get returns a value of the node on index.
// If index is not valid, empty string and false would be returned.
func (ll *LinkedList) Get(index int) (string, bool) {
	if index < 0 || index >= ll.length {
		return "", false
	}

	curr := ll.get(index)
	return curr.val, true
}

// PushForward inserts new node to the head of list.
// It returns length of the list after adding new node.
func (ll *LinkedList) PushForward(val string) int {
	ll.length++

	if ll.length == 1 {
		newListNode := &ListNode{val: val}
		ll.head = newListNode
		ll.tail = newListNode
	} else {
		newListNode := &ListNode{val: val, next: ll.head}
		ll.head.prev = newListNode
		ll.head = newListNode
	}

	return ll.length
}

// PushBack inserts new node to the tail of list.
// It returns length of the list after adding new node.
func (ll *LinkedList) PushBack(val string) int {
	ll.length++

	if ll.length == 1 {
		newListNode := &ListNode{val: val}
		ll.head = newListNode
		ll.tail = newListNode
	} else {
		newListNode := &ListNode{val: val, prev: ll.tail}
		ll.tail.next = newListNode
		ll.tail = newListNode
	}

	return ll.length
}

// PushAtIndex inserts new node at the place of node with index.
// It returns length of the list after adding new node or -1 in case of invalid index.
func (ll *LinkedList) PushAtIndex(index int, val string) int {
	if index > ll.length || index < 0 {
		return -1
	}

	if index == 0 {
		return ll.PushForward(val)
	}

	if index == ll.length {
		return ll.PushBack(val)
	}

	prev := ll.get(index - 1)
	next := prev.next

	// inserting itself
	newListNode := &ListNode{val: val, prev: prev, next: next}
	prev.next = newListNode
	next.prev = newListNode
	ll.length++

	return ll.length
}

// PopForward deletes the node from the head of list.
// It returns value of deleted node.
// Second value is false in case list is empty.
func (ll *LinkedList) PopForward() (string, bool) {
	if ll.length == 0 {
		return "", false
	}

	result := ll.head.val
	ll.head = ll.head.next
	if ll.head != nil {
		ll.head.prev = nil
	} else {
		ll.tail = nil
	}
	ll.length--

	return result, true
}

// PopBack deletes the node from the tail of list.
// It returns value of deleted node.
// Second value is false in case list is empty.
func (ll *LinkedList) PopBack() (string, bool) {
	if ll.length == 0 {
		return "", false
	}

	result := ll.tail.val
	ll.tail = ll.tail.prev
	if ll.tail != nil {
		ll.tail.next = nil
	} else {
		ll.head = nil
	}
	ll.length--

	return result, true
}

// PopAtIndex deletes the node from the place with the index.
// It returns value of deleted node.
// Second value is false in case list doesn't have node with this index.
func (ll *LinkedList) PopAtIndex(index int) (string, bool) {
	if index < 0 || index >= ll.length {
		return "", false
	}

	if index == 0 {
		return ll.PopForward()
	}

	if index == ll.length-1 {
		return ll.PopBack()
	}

	curr := ll.get(index)

	// deleting itself
	prev := curr.prev
	next := curr.next
	prev.next = next
	next.prev = prev

	ll.length--

	return curr.val, true
}

// LRange returns node values from indexes in range [start, stop].
func (ll *LinkedList) LRange(start, stop int) []string {
	// if indexes are negative
	if start < 0 {
		start = start + ll.length
	}
	if stop < 0 {
		stop = stop + ll.length
	}

	if stop < 0 {
		return []string{}
	}
	if start < 0 {
		if stop < 0 {
			return []string{}
		} else {
			start = 0
		}
	}

	// if first index is equal or greated than list length or stop < start,
	// there are no values to return
	if start >= ll.length || stop < start {
		return []string{}
	}

	result := []string{}

	// first node in range
	curr := ll.get(start)
	result = append(result, curr.val)
	// add to result next (stop-start) nodes
	for range stop - start {
		curr = curr.next
		if curr == nil {
			return result
		}
		result = append(result, curr.val)
	}

	return result
}
