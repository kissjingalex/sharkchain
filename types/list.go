package types

import (
	"fmt"
	"reflect"
)

type List[T any] struct {
	Data []T
}

func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

func (l *List[T]) Insert(item T) {
	l.Data = append(l.Data, item)
}

func (l *List[T]) Get(index int) T {
	if index > len(l.Data)-1 {
		err := fmt.Sprintf("the given index (%d) is higher than the length (%d)", index, len(l.Data))
		panic(err)
	}
	return l.Data[index]
}

func (l *List[T]) Clear() {
	l.Data = []T{}
}

func (l *List[T]) GetIndex(v T) int {
	for i, item := range l.Data {
		if reflect.DeepEqual(v, item) {
			return i
		}
	}

	return -1
}

func (l *List[T]) Contains(v T) bool {
	index := l.GetIndex(v)

	return index >= 0
}

func (l *List[T]) Pop(index int) {
	if index < 0 {
		return
	}
	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

func (l *List[T]) Remove(v T) {
	index := l.GetIndex(v)
	l.Pop(index)
}

func (l List[T]) Last() T {
	return l.Data[l.Len()-1]
}

func (l *List[T]) Len() int {
	return len(l.Data)
}
