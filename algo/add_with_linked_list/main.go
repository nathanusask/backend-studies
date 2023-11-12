package main

/*
You are given two non-empty linked lists representing two non-negative integers.
The digits are stored in reverse order and each of their nodes contain a single
digit. Add the two numbers and return it as a linked list.
*/

import "fmt"

// ListNode is a node in a linked list
type ListNode struct {
	Val  int
	Next *ListNode
}

// newList creates a linked list from an integer
func newList(val int) *ListNode {
	var head, cur *ListNode
	for val > 0 {
		tmp := &ListNode{Val: val % 10}
		if head == nil {
			head = tmp
			cur = head
		} else {
			cur.Next = tmp
			cur = cur.Next
		}
		val /= 10
	}
	return head
}

// String returns a string representation of a linked list
func (l *ListNode) String() string {
	var s string
	for l != nil {
		s += fmt.Sprintf("%d", l.Val)
		l = l.Next
	}
	return s
}

// addTwoNumbers adds two numbers represented by linked lists
func addTwoNumbers(head1 *ListNode, head2 *ListNode) *ListNode {
	var head, cur *ListNode
	carry := 0
	for head1 != nil || head2 != nil {
		sum := carry
		if head1 != nil {
			sum += head1.Val
			head1 = head1.Next
		}
		if head2 != nil {
			sum += head2.Val
			head2 = head2.Next
		}
		carry = sum / 10
		tmp := &ListNode{Val: sum % 10}
		if head == nil {
			head = tmp
			cur = head
		} else {
			cur.Next = tmp
			cur = cur.Next
		}
	}
	if carry > 0 {
		cur.Next = &ListNode{Val: carry}
	}
	return head
}

func main() {
	l1 := newList(789)
	l2 := newList(465)
	fmt.Println(l1)
	fmt.Println(l2)
	fmt.Println(addTwoNumbers(l1, l2))
}
