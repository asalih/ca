package main

import (
	"bufio"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/asalih/ca/ca"
)

func main() {
	startingGs := runtime.NumGoroutine()
	go http.ListenAndServe("localhost:6060", nil)

	ca.NewManager("http://php.testsparker.com/process.php?file=Generics/index.nsp", "php.testsparker.com").CrawlAndAttack()
	endingGs := runtime.NumGoroutine()

	fmt.Println("Number of goroutines before:", startingGs)
	fmt.Println("Number of goroutines after :", endingGs)

	fmt.Println("Waiting an input for exit: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())
}
