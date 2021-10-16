package hw04lrucache

import "errors"

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) (*ListItem, error)
	PushBack(v interface{}) (*ListItem, error)
	Remove(i *ListItem) error
	MoveToFront(i *ListItem) error
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	First *ListItem
	Last  *ListItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.Last
}

func (l *list) Back() *ListItem {
	return l.First
}

func (l *list) PushFront(v interface{}) (*ListItem, error) {
	if v == nil {
		return nil, errors.New("nil input not allowed")
	}

	elem := &ListItem{
		Value: v,
		Next:  l.Last,
		Prev:  nil,
	}
	if l.Last == nil {
		l.First = elem
	} else {
		l.Last.Prev = elem
	}
	l.Last = elem
	l.len++
	return elem, nil
}

func (l *list) PushBack(v interface{}) (*ListItem, error) {
	if v == nil {
		return nil, errors.New("nil input not allowed")
	}

	elem := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.First,
	}
	if l.First == nil {
		l.Last = elem
	} else {
		l.First.Next = elem
	}
	l.First = elem
	l.len++
	return elem, nil
}

func (l *list) Remove(i *ListItem) error {
	if i == nil {
		return errors.New("empty element")
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.First = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.Last = i.Next
	}

	l.len--

	return nil
}

func (l *list) MoveToFront(i *ListItem) error {
	if err := l.Remove(i); err != nil {
		return err
	}
	_, _ = l.PushFront(i.Value)
	return nil
}

func NewList() List {
	return new(list)
}
