package ca

import (
	"fmt"
	"sync"
	"time"

	"github.com/robertkrimen/otto/ast"
)

//Attacker Attacker model
type Attacker struct {
	CrawlDataReader      func() *crawlData
	ScriptsManager       *ScriptsManager
	finished             bool
	ticker               *time.Ticker
	tickerChannel        chan bool
	refWg                *sync.WaitGroup
	FoundVulnerabilities []*Vulnerability
}

//NewAttacker Attacker initiator
func NewAttacker(scriptsManager *ScriptsManager) *Attacker {
	return &Attacker{nil, scriptsManager, false, nil, nil, nil, nil}
}

func (attacker *Attacker) attack() {
	if attacker.CrawlDataReader == nil {
		return
	}

	cData := attacker.CrawlDataReader()

	if cData == nil {
		return
	}

	//loop
	defer attacker.attack()

	fmt.Println("Attacking " + cData.URL.RequestURI())

	var wg sync.WaitGroup
	wg.Add(len(attacker.ScriptsManager.ActiveScripts))

	for _, p := range attacker.ScriptsManager.ActiveScripts {
		go func(scr *ast.Program) {
			defer wg.Done()

			result := executeAttackScript(scr, cData)
			if result == nil {
				return
			}

			attacker.FoundVulnerabilities = append(attacker.FoundVulnerabilities, result...)
		}(p)
	}

	wg.Wait()

	if cData.Script != nil {
		result := executeAttackScript(cData.Script, cData)
		if result == nil {
			return
		}

		attacker.FoundVulnerabilities = append(attacker.FoundVulnerabilities, result...)
	}
}

//Start Starts attacking phase
func (attacker *Attacker) Start(wg *sync.WaitGroup) {
	attacker.ticker = time.NewTicker(500 * time.Millisecond)
	attacker.tickerChannel = make(chan bool)
	attacker.refWg = wg

	go func() {
		for {
			select {
			case <-attacker.tickerChannel:
				return
			case <-attacker.ticker.C:
				attacker.attack()
			}
		}
	}()
}

//Finalize Cleanup func for attacking phase
func (attacker *Attacker) Finalize() {
	attacker.tickerChannel <- true
	attacker.ticker.Stop()
	close(attacker.tickerChannel)

	attacker.attack()
	attacker.refWg.Done()
	attacker.finished = true
}
