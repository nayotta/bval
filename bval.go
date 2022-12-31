package bval

import "sync"

type Value[T comparable] struct {
	mutex sync.RWMutex
	value T

	onSet     func(old, new T)
	onChanged func(old, new T)
}

type option[T comparable] func(v *Value[T])

func WithInitValue[T comparable](initValue T) option[T] {
	return func(v *Value[T]) {
		v.value = initValue
	}
}

func WithOnSet[T comparable](onSet func(old, new T)) option[T] {
	return func(v *Value[T]) {
		v.onSet = onSet
	}
}

func WithOnChanged[T comparable](onChanged func(old, new T)) option[T] {
	return func(v *Value[T]) {
		v.onChanged = onChanged
	}
}

func New[T comparable](opts ...option[T]) *Value[T] {
	v := new(Value[T])

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *Value[T]) set(value T) {
	old := v.value
	v.value = value

	if v.onSet != nil {
		v.onSet(old, value)
	}

	if v.onChanged != nil && old != value {
		v.onChanged(old, value)
	}
}

func (v *Value[T]) Set(value T) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.set(value)
}

func (v *Value[T]) Operate(calcFunc func(old T) T) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.set(calcFunc(v.value))
}

func (v *Value[T]) Get() T {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.value
}
