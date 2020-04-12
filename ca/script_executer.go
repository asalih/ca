package ca

import (
	"fmt"
	"net/url"

	"github.com/robertkrimen/otto"
	"github.com/robertkrimen/otto/ast"
)

func executeCrawlScript(scr *ast.Program, u string) []*crawlData {
	URL, uerr := url.Parse(u)

	if uerr != nil {
		panic(uerr)
	}

	vm := otto.New()

	vm.Run(ScriptTemplate)
	vm.Run(scr)

	requests, _ := vm.Get("requests")
	requestsList, _ := requests.Export()
	requestsMap := requestsList.([]map[string]interface{})

	var crawlDatas []*crawlData
	for _, aObj := range requestsMap {
		URL.Path = aObj["url"].(string)

		cData := newCrawlData(URL, aObj["method"].(string))
		cData.Script = scr

		if aObj["parameters"] != nil {
			cData.Params = aObj["parameters"].([]string)
		}

		crawlDatas = append(crawlDatas, cData)
	}

	return crawlDatas
}

func executeAttackScript(scr *ast.Program, crwl *crawlData) []*Vulnerability {
	vm := otto.New()

	vm.Run(ScriptTemplate)
	vm.Run(scr)

	attacks, _ := vm.Get("attacks")

	attackHandler := func(attackStr string) *Vulnerability {
		handler := &httpRequestHandler{crwl, attackStr}
		response := handler.Do()

		method, _ := vm.Get("analyze")
		analyzer, anerr := method.Call(method, response)

		if anerr != nil {
			fmt.Println(anerr)
		}

		analyzeResult, _ := analyzer.Export()

		if analyzeResult == nil {
			return nil
		}

		vulnerabilityData := NewVulnerability(analyzeResult.(map[string]interface{}), crwl)

		if vulnerabilityData == nil {
			return nil
		}

		return vulnerabilityData
	}

	attacksList, _ := attacks.Export()

	var vulnerabilities []*Vulnerability
	if attacksList == nil {
		v := attackHandler("")
		if v != nil {
			vulnerabilities = append(vulnerabilities, v)
		}
		return vulnerabilities
	}

	attacksListArray := attacksList.([]map[string]interface{})
	for _, aObj := range attacksListArray {
		attackStr := aObj["attack"].(string)
		v := attackHandler(attackStr)
		if v != nil {
			vulnerabilities = append(vulnerabilities, v)
		}

	}

	return vulnerabilities
}
