package ca

import (
	"container/list"
	"net/url"
	"sync"
)

type Manager struct {
	Crawler  *Crawler
	Attacker *Attacker
	queue    *list.List
	mutex    sync.Mutex
}

func NewManager(uriStr string) *Manager {
	uriParsed, _ := url.Parse(uriStr)

	return &Manager{NewCrawler(uriParsed.String(), uriParsed.Host), NewAttacker(), list.New(), sync.Mutex{}}
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

func (m *Manager) IsFinished() bool {
	return m.Crawler.finished && m.Attacker.finished
}
