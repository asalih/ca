package ca

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
)

//ScriptTemplate base script of the context
var ScriptTemplate string = `
	var severities = ["Best Practice", "Information", "Low", "Medium", "High", "Critical"]
	var BEST_PRACTICE = 0, INFORMATION = 1, LOW = 2, MEDIUM = 3, HIGH = 4, CRITICAL = 5;

	function Found(severity, title, additionalData){
		return {Title: title, Severity: severities[severity], AdditionalData: additionalData}
	}
`

//ScriptsManager Scripts file manager
type ScriptsManager struct {
	RootDir       string
	ActiveScripts []*ast.Program
	CrawlScripts  []*ast.Program
}

//NewScriptsManager Initiates the script files
func NewScriptsManager(rootDir string) *ScriptsManager {
	scriptsManager := &ScriptsManager{rootDir, nil, nil}

	scriptsManager.ActiveScripts = scriptsManager.readScriptFiles("active/")
	scriptsManager.CrawlScripts = scriptsManager.readScriptFiles("crawl_scripts/")

	return scriptsManager
}

func (scripts *ScriptsManager) readScriptFiles(dir string) []*ast.Program {
	var srcDir = scripts.RootDir + dir
	files, _ := ioutil.ReadDir(srcDir)

	var ref []*ast.Program
	for _, v := range files {
		if v.IsDir() {
			ref = append(ref, scripts.readScriptFiles(dir+v.Name()+"/")...)
		}

		if !strings.HasSuffix(v.Name(), ".js") {
			continue
		}
		program, err := parser.ParseFile(nil, srcDir+v.Name(), nil, 0)

		if err != nil {
			fmt.Println("Proglem with the file:")
			fmt.Println(err)
		}

		ref = append(ref, program)
	}

	return ref
}
