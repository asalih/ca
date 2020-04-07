package ca

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
)

type Attacker struct {
	Reader               func() *CrawlData
	ticker               *time.Ticker
	tickerChannel        chan bool
	refWg                *sync.WaitGroup
	scripts              []*ast.Program
	FoundVulnerabilities []*Vulnerability
}

var template string = `
	var BEST_PRACTICE = 0, INFORMATION = 1, LOW = 2, MEDIUM = 3, HIGH = 4, CRITICAL = 5;

	function Found(severity, title, additionalInfo){
		return {Title: title, Severity: severity, AdditionalInfo: additionalInfo}
	}
`

func NewAttacker() *Attacker {
	return &Attacker{nil, nil, nil, nil, nil, nil}
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
	wg.Add(len(attacker.scripts))

	for _, p := range attacker.scripts {
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

func (attacker *Attacker) readScriptFiles(subDir string) {
	var srcDir = "./scripts/" + subDir
	files, _ := ioutil.ReadDir(srcDir)
	for _, v := range files {
		if v.IsDir() {
			attacker.readScriptFiles(v.Name() + "/")
		}

		if !strings.HasSuffix(v.Name(), ".js") {
			continue
		}
		program, err := parser.ParseFile(nil, srcDir+v.Name(), nil, 0)

		if err != nil {
			fmt.Println("Proglem with the file:")
			fmt.Println(err)
		}

		attacker.scripts = append(attacker.scripts, program)
	}
}

func (attacker *Attacker) Start(wg *sync.WaitGroup) {
	attacker.ticker = time.NewTicker(500 * time.Millisecond)
	attacker.tickerChannel = make(chan bool)
	attacker.refWg = wg

	attacker.readScriptFiles("")

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
}
