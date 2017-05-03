package store

import (
	"errors"
	"log"
	"sync"
)

// 用于在内存中处理实时日志通讯
// websocket channel set, delete
// tcp write

type MonitorList struct {
	chs map[string][](chan []byte)
	mu  sync.Mutex
}

func NewMonitorList() *MonitorList {
	chs := make(map[string][](chan []byte), 1000)
	return &MonitorList{chs: chs}
}

// 根据key设置通道
func (m *MonitorList) Set(key string, wsch chan []byte) error {
	log.Println("key: ", key)
	m.mu.Lock()
	_, ok := m.chs[key]
	if !ok {
		ch := make([](chan []byte), 0)
		ch = append(ch, wsch)
		m.chs[key] = ch
	} else {
		chs := m.chs[key]
		m.chs[key] = append(chs, wsch)
	}
	m.mu.Unlock()
	return nil
}

// 根据key删除对应channel
func (m *MonitorList) Delete(key string) error {
	_, ok := m.chs[key]
	if !ok {
		return errors.New("not found channels for " + key)
	}
	m.mu.Lock()
	delete(m.chs, key)
	//暂时不处理已关联的大量channel
	m.mu.Unlock()
	return nil
}

// 向对应channel写数据
// 同时清理已关闭的channel
func (m *MonitorList) Write(key string, content []byte) error {
	wsch, ok := m.chs[key]
	if !ok {
		log.Println("not found channels for " + key)
		return errors.New("not found channels for " + key)
	}

	chsNew := make([](chan []byte), 0)
	var wg sync.WaitGroup
	var chlock sync.Mutex

	for _, ch := range wsch {
		wg.Add(1)
		// 防止channel close write panic
		go func(ch chan []byte, d []byte) {
			defer func() {
				wg.Add(-1)
				if err := recover(); err == nil {
					chlock.Lock()
					chsNew = append(chsNew, ch)
					chlock.Unlock()
				} else {
					log.Println("write channel error!")
				}

			}()
			ch <- d
		}(ch, content)

	}
	wg.Wait()
	m.mu.Lock()
	m.chs[key] = chsNew
	m.mu.Unlock()

	return nil
}
