package ca

import (
	"fmt"
	"sync"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/ast"
)

type Attacker struct {
	Reader               func() *CrawlData
	ScriptsManager       *ScriptsManager
	finished             bool
	ticker               *time.Ticker
	tickerChannel        chan bool
	refWg                *sync.WaitGroup
	FoundVulnerabilities []*Vulnerability
}

var template string = `
	var severities = ["Best Practice", "Information", "Low", "Medium", "High", "Critical"]
	var BEST_PRACTICE = 0, INFORMATION = 1, LOW = 2, MEDIUM = 3, HIGH = 4, CRITICAL = 5;

	function Found(severity, title, additionalData){
		return {Title: title, Severity: severities[severity], AdditionalData: additionalData}
	}
`

func NewAttacker(scriptsManager *ScriptsManager) *Attacker {
	return &Attacker{nil, scriptsManager, false, nil, nil, nil, nil}
}

func (attacker *Attacker) attack() {
	if attacker.Reader == nil {
		return
	}

	crawlData := attacker.Reader()

	if crawlData == nil {
		return
	}

	//loop
	defer attacker.attack()

	fmt.Println("Attacking " + crawlData.URL.RequestURI())

	var wg sync.WaitGroup
	wg.Add(len(attacker.ScriptsManager.ActiveScripts))

	for _, p := range attacker.ScriptsManager.ActiveScripts {
		go func(scr *ast.Program, crwl *CrawlData, wg *sync.WaitGroup) {
			defer wg.Done()
			vm := otto.New()

			vm.Run(template)
			vm.Run(scr)

			attacks, _ := vm.Get("attacks")

			attackHandler := func(attackStr string) {
				handler := &httpRequestHandler{crwl, attackStr}
				response := handler.Do()

				method, _ := vm.Get("analyze")
				analyzer, anerr := method.Call(method, response)

				if anerr != nil {
					fmt.Println(anerr)
				}

				analyzeResult, _ := analyzer.Export()

				if analyzeResult == nil {
					return
				}

				vulnerabilityData := NewVulnerability(analyzeResult.(map[string]interface{}), crawlData)

				if vulnerabilityData == nil {
					return
				}

				attacker.FoundVulnerabilities = append(attacker.FoundVulnerabilities, vulnerabilityData)
			}

			attacksList, _ := attacks.Export()

			if attacksList == nil {
				attackHandler("")
			} else {
				attacksListArray := attacksList.([]map[string]interface{})
				for _, aObj := range attacksListArray {
					attackStr := aObj["attack"].(string)
					attackHandler(attackStr)
				}
			}

		}(p, crawlData, &wg)
	}

	wg.Wait()
}

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

func (attacker *Attacker) Finalize() {
	attacker.tickerChannel <- true
	attacker.ticker.Stop()
	close(attacker.tickerChannel)

	attacker.attack()
	attacker.refWg.Done()
	attacker.finished = true
}
