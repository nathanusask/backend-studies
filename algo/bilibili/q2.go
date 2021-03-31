package main

import (
	"fmt"
	"sort"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func print(head *ListNode) {
	for cur := head; cur != nil; cur = cur.Next {
		if cur.Next != nil {
			fmt.Print(cur.Val, "->")
		} else {
			fmt.Print(cur.Val, "\n")
		}
	}
}

func insertNode(head *ListNode, node *ListNode) *ListNode {
	// when the head is nil
	if head == nil {
		head = node
		node.Next = nil
		return head
	}
	// when the to-be-inserted node is the smallest
	if node.Val < head.Val {
		node.Next = head
		head = node
		return head
	}

	cur := head
	var prev *ListNode
	for cur != nil && node.Val > cur.Val {
		prev = cur
		cur = cur.Next
	}
	if prev == nil {
		// the node should be inserted to the last and there's currently only one node in the list
		cur.Next = node
		node.Next = nil
		return head
	}
	prev.Next = node
	node.Next = cur
	return head
}

func (l *ListNode) OrderAsc() {

}

func sortInList(head *ListNode) *ListNode {
	// if head == nil || head.Next == nil {
	// 	return head
	// }
	// var newHead *ListNode
	// cur := head
	// for cur != nil {
	// 	head = head.Next
	// 	newHead = insertNode(newHead, cur)
	// 	print(newHead)
	// 	cur = head
	// }

	// return newHead

	var buf []int
	for cur := head; cur != nil; cur = cur.Next {
		buf = append(buf, cur.Val)
	}

	sort.Ints(buf)
	var newHead *ListNode
	var rear *ListNode
	var cur *ListNode
	for _, val := range buf {
		cur = &ListNode{
			Val:  val,
			Next: nil,
		}
		if newHead == nil {
			newHead = cur
		} else {
			rear.Next = cur
		}
		rear = cur
	}

	return newHead
}

func main() {
	arr := []int{1, 3, 2, 4, 5}
	var head *ListNode
	var rear *ListNode
	var cur *ListNode
	for _, val := range arr {
		cur = &ListNode{
			Val:  val,
			Next: nil,
		}
		if head == nil {
			head = cur
		} else {
			rear.Next = cur
		}
		rear = cur
	}
	print(head)

	s := make(int chan)
	

	newHead := sortInList(head)
	print(newHead)
}
