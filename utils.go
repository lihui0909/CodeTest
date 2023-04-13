package main

import (
	"log"
)

// Debugging
const Debug = true

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

//自定义set结构体，实现创建、添加、判断包含、求元素个数、去掉元素方法
type Set struct {
	m map[interface{}]struct{}
	//	lock sync.Mutex
}

func New() *Set {
	s := &Set{}
	s.m = make(map[interface{}]struct{})
	return s
}

func (s *Set) Add(items ...interface{}) error {
	for _, item := range items {
		if s.Contains(item) {
			continue
		}
		s.m[item] = struct{}{}
	}
	return nil
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}
func (s *Set) Size() int {
	return len(s.m)
}

func (s *Set) Remove(items ...interface{}) error {
	for _, item := range items {
		if !s.Contains(item) {
			continue
		}
		delete(s.m, item)
	}
	return nil
}
