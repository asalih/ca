package ca

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
)

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
