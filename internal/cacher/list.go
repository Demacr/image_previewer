package cacher //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	front    *listItem
	back     *listItem
	capacity int
}

func (l list) Len() int {
	return l.capacity
}

func (l list) Front() *listItem {
	return l.front
}

func (l list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	newElem := listItem{v, l.front, nil}
	if l.front != nil {
		l.front.Prev = &newElem
		l.front = &newElem
	} else {
		l.front = &newElem
		l.back = &newElem
	}
	l.capacity++
	return &newElem
}

func (l *list) PushBack(v interface{}) *listItem {
	newElem := listItem{v, nil, l.back}
	if l.back != nil {
		l.back.Next = &newElem
		l.back = &newElem
	} else {
		l.front = &newElem
		l.back = &newElem
	}
	l.capacity++
	return &newElem
}

func (l *list) Remove(i *listItem) {
	switch {
	case l.front == i && l.back == i:
		l.front = nil
		l.back = nil
	case l.front == i:
		l.front = l.front.Next
		l.front.Prev = nil
	case l.back == i:
		l.back = l.back.Prev
		l.back.Next = nil
	default:
		i.Next.Prev = i.Prev
		i.Prev.Next = i.Next
	}
	l.capacity--
}

func (l *list) MoveToFront(i *listItem) {
	switch {
	case l.front == i:
		return
	case l.back == i:
		l.back.Prev.Next = nil
		l.back = l.back.Prev
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
}

// NewList returns address of new list{}
//
func NewList() List {
	return &list{}
}
