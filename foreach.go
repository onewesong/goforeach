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

	"github.com/akamensky/argparse"
)

var wg sync.WaitGroup

func run(cmd string, ch chan bool) {
	defer wg.Done()

	ch <- true
	args := strings.Fields(cmd)
	res, err := exec.Command(args[0], args[1:]...).Output()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Print(string(res))
	<-ch
}

func main() {
	parser := argparse.NewParser("foreach", "run the command which you want by goroutine")
	rawCmd := parser.String("e", "execute", &argparse.Options{Required: true, Help: "the command to be executed"})
	cpuNum := parser.Int("n", "cpu_num", &argparse.Options{Required: false, Default: 2, Help: "sets the maximum number of CPUs that can be executing simultaneously"})
	fork := parser.Int("f", "fork", &argparse.Options{Required: false, Default: 10, Help: "specify the number of concurrent goroutine"})
	show := parser.Flag("s", "show", &argparse.Options{Required: false, Default: false, Help: "show infos of localhost"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

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
