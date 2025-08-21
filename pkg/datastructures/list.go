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

// Get returns a value of the node on index.
// If index is not valid, empty string and false would be returned.
func (ll *LinkedList) Get(index int) (string, bool) {
	if index < 0 || index >= ll.length {
		return "", false
	}

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

	var prev *ListNode // should be at index-1
	if index-1 < ll.length/2 {
		prev = ll.head
		for range index - 1 {
			prev = prev.next
		}
	} else {
		prev = ll.tail
		for range ll.length - index {
			prev = prev.prev
		}
	}
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

	// deleting itself
	prev := curr.prev
	next := curr.next
	prev.next = next
	next.prev = prev

	ll.length--

	return curr.val, true
}
