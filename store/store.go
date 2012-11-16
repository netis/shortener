package store

import (
	"sync"
	"fmt"
	"strconv"
)

type URLStore struct {
	urls map[string]string
	mu   sync.RWMutex
}

func (s *URLStore) Get(key string) string {
	//s.mu.RLock()
	//defer s.mu.RUnlock()
	url := s.urls[key]
	s.mu.RUnlock()
	return url
}

func (s *URLStore) Set(key string, url string) bool {
	fmt.Printf("locking stuff\n")
	//s.mu.Lock()
	fmt.Printf("locked\n")
	//defer s.mu.Unlock()
	fmt.Printf("trying to Set %s in store, with key %s\n", url, key)
	if _, present := s.urls[key]; present {
		fmt.Printf("%s present\n", key)
		return false
	}
	fmt.Printf("inserting %s in store with key %s\n", url, key)
	s.urls[key] = url
	return true
}

func NewURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
}

func (s *URLStore) Count() int {
	//s.mu.RLock()
	//defer s.mu.RUnlock()
	return len(s.urls)
}

func (s *URLStore) Put(url string) string {
	for {
		key := genKey(s.Count())
		fmt.Printf("trying to Put %s in store, with key %s\n", url, key)
		if s.Set(key, url) {
			fmt.Printf("success\n")
			return key
		}
		fmt.Printf("failed\n")
	}
	return ""
}

func genKey(value int) string {
	return strconv.Itoa(value)
}
