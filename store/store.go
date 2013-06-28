package store

import (
	"strconv"
	"os"
	"io"
	"encoding/gob"
	"sync"
	"log"
)

type URLStore struct {
	sync.Mutex
	urls map[string]string
	file *os.File
	  
}

func (s *URLStore) Get(key string) string {
	s.Lock()
	defer s.Unlock()
	url := s.urls[key]
	return url
}
func (s *URLStore) Set(key string, url string) bool {
	s.Lock()
	defer s.Unlock()
	return s.set(key, url)

}
func (s *URLStore) set(key string, url string) bool {
	
	//log.Printf("trying to Set %s in store, with key %s\n", url, key)
	if _, present := s.urls[key]; present {
		//log.Print("%s present\n", key)
		return false
	}
	//log.Printf("inserting %s in store with key %s\n", url, key)
	s.urls[key] = url
	return true
}

func NewURLStore(filename string) *URLStore {
	s := &URLStore{urls: make(map[string]string)}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		//log.Fatal("URLStore:", err)
	}
	s.file = f
	if err := s.load(); err != nil {
		//log.Printf("Error loading data in URLStore:", err)
	}
	return s
}

func (s *URLStore) Count() int {
	return len(s.urls)
}

func (s *URLStore) Put(url string) string {
	s.Lock()
	defer s.Unlock()
	for {
		key := genKey(s.Count())
		////log.Printf("trying to Put %s in store, with key %s\n", url, key)
		if s.set(key, url) {
			if err := s.save(key,url); err != nil {
				log.Printf("Error saving data in URLStore:", err)
			}
			//log.Printf("success\n")
			return key
		}
		//log.Printf("failed\n")
	}
	return ""
}

func genKey(value int) string {
	return strconv.Itoa(value)
}

type record struct {
	Key, URL string
}

func (s *URLStore) save(key string, url string) error {
	e := gob.NewEncoder(s.file)
	return e.Encode(record{key, url})
}

func (s *URLStore) load() error {
	if _, err := s.file.Seek(0,0); err != nil {
		return err
	}
	d := gob.NewDecoder(s.file)
	var err error
	for err == nil {
		var r record
		if err = d.Decode(&r); err == nil {
			s.Set(r.Key, r.URL)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}
