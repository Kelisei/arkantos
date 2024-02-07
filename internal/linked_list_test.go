package internal_test

import (
	"arkantos/internal"
	"testing"
)

func TestAppendLinkedList(t *testing.T) {
	t.Run("Append 1 element", func(t *testing.T) {
		list := internal.LinkedList[string]{}
		list.Append("Goku")
		if list.Size == 0 {
			t.Fatal("Element wasn't added to list")
		}
		if list.FirstNode.Data != "Goku" {
			t.Fatal("The elemen't added doesn't have the correct value")
		}
	})

	t.Run("Append 2 elements", func(t *testing.T) {
		list := internal.LinkedList[string]{}
		list.Append("Goku")
		list.Append("Vegeta")
		if list.Size == 0 {
			t.Fatal("Elements wasn't added to list")
		}
		if list.LastNode.Data != "Vegeta" {
			t.Fatal("The last element doesn't have the correct value")
		}
		if list.LastNode.Previous != list.FirstNode {
			t.Fatal("Link from last element to first element wasn't created")
		}
		if list.FirstNode.Next != list.LastNode {
			t.Fatal("Link from first element to last element wasn't created")
		}
	})
}

func TestPopLinkedList(t *testing.T) {
	t.Run("Pop with 0 elements", func(t *testing.T) {
		list := internal.LinkedList[string]{}
		_, err := list.Pop()
		if err == nil {
			t.Fatal("Didn't return an error when popping in a empty list")
		}
	})

	t.Run("Pop with 1 elements", func(t *testing.T) {
		list := internal.LinkedList[string]{}
		list.Append("Goku")
		value, err := list.Pop()
		if err != nil {
			t.Fatal("Failed to pop element out of list")
		}
		if *value != "Goku" {
			t.Fatal("Didn't return a correct value when popping out of list")
		}
		if list.Size != 0 {
			t.Fatal("Failed to update list size")
		}
	})

	t.Run("Pop with 2 elements", func(t *testing.T) {
		list := internal.LinkedList[string]{}
		list.Append("Goku")
		list.Append("Vegeta")
		value, err := list.Pop()
		if err != nil {
			t.Fatal("Failed to pop element out of list")
		}
		if *value != "Vegeta" {
			t.Fatal("Didn't return a correct value when popping out of list")
		}
		if list.Size != 1 {
			t.Fatal("Failed to update list size")
		}
	})
}
