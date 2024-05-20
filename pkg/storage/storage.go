package storage

import "sync"

// Storage - объект для хранения данных внутри памяти
type Storage[T any, S ~[]T] interface {
	Add(...T)
	Flush() S
	Len() int
}

type storage[T any, S ~[]T] struct {
	col       S
	mux       *sync.RWMutex
	allocSize int
}

// New Example:
//
// storage := storage.New[*fact.Fact, fact.Collection](100)
//
// storage.Add(&fact.Fact{Value:100})
func New[T any, S ~[]T](allocSize int) Storage[T, S] {
	return &storage[T, S]{
		col:       make(S, 0, allocSize),
		allocSize: allocSize,
		mux:       &sync.RWMutex{},
	}
}

// Add Добавляет элемент в хранилище
func (s *storage[T, S]) Add(elem ...T) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.col = append(s.col, elem...)

}

// Flush вынимает всё из хранилища
func (s *storage[T, S]) Flush() S {
	s.mux.Lock()
	defer s.mux.Unlock()

	defer func() {
		s.col = make(S, 0, s.allocSize)
	}()

	return s.col
}

// Len получаем кол-во элементов в хранилище
func (s *storage[T, S]) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return len(s.col)
}
