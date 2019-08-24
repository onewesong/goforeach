// version: 0.0.1
// author: wws
// Created Time: 2019-08-22
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/alecthomas/kingpin.v2"
)

var wg sync.WaitGroup

func run(cmd string, ch chan bool) {
	defer wg.Done()

	ch <- true
	args := strings.Fields(cmd)
	res, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		fmt.Println(err.Error())
		<-ch
		return
	}
	fmt.Print(string(res))
	<-ch
}

var (
	app = kingpin.New("foreach", "run the command which you want by goroutine.")

	rawCmd = app.Arg("execute", "the command to be executed").Required().String()
	cpuNum = app.Flag("cpu_num", "sets the maximum number of CPUs that can be executing simultaneously").Short('n').Default("2").Int()
	fork   = app.Flag("fork", "specify the number of concurrent goroutine").Short('f').Default("10").Int()
	show   = app.Flag("show", "show infos of localhost").Short('s').Bool()
)

func main() {
	app.Version("0.0.2")
	app.HelpFlag.Short('h')
	app.Parse(os.Args[1:])

	logicCPUNum := runtime.NumCPU()
	if *show == true {
		fmt.Printf("logic cpu num: %d\n", logicCPUNum)
		os.Exit(0)
	}

	var cpuNumber int
	if *cpuNum > logicCPUNum {
		cpuNumber = logicCPUNum
	} else {
		cpuNumber = *cpuNum
	}
	runtime.GOMAXPROCS(cpuNumber)

	ch := make(chan bool, *fork)

	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		cmd := strings.Replace(*rawCmd, "#1", input, -1)
		wg.Add(1)
		go run(cmd, ch)
	}

	wg.Wait()
	close(ch)
}
