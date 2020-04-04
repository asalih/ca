package ca

import (
	"container/list"
	"sync"
)

type Manager struct {
	Crawler  *Crawler
	Attacker *Attacker
	queue    *list.List
	mutex    sync.Mutex
}

func NewManager(url string, domain string) *Manager {
	return &Manager{NewCrawler(url, domain), NewAttacker(), list.New(), sync.Mutex{}}
}

func (m *Manager) CrawlAndAttack() {
	m.Crawler.Init()

	m.Crawler.OnRequest = func(c *CrawlData) {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		m.queue.PushBack(c)
	}

	m.Crawler.OnCrawlFinish = func() {
		m.Attacker.Finalize()
	}

	m.Attacker.Reader = func() *CrawlData {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		e := m.queue.Front()
		if e == nil {
			return nil
		}

		m.queue.Remove(e)

		if e.Value == nil {
			return nil
		}

		return e.Value.(*CrawlData)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go m.Crawler.Start(&wg)
	go m.Attacker.Start(&wg)

	wg.Wait()
}
