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
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	wg  sync.WaitGroup
	app = kingpin.New("foreach", "run the command which you want by goroutine.")

	rawCmd = app.Arg("execute", "The command to be executed").Required().String()
	circle = app.Flag("circle", "The circle times of command run").Short('c').Default("1").Int()
	cpuNum = app.Flag("cpu_num", "Sets the maximum number of CPUs that can be executing simultaneously").Short('n').Default("2").Int()
	fork   = app.Flag("fork", "Specify the number of concurrent goroutine").Short('f').Default("10").Int()
	show   = app.Flag("show", "Show infos of localhost").Short('s').Bool()
)

func run(cmd string, ch chan bool) {
	defer wg.Done()

	ch <- true
	res, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
		<-ch
		return
	}
	fmt.Print(string(res))
	<-ch
}

func replaceCmd(cmd string, input string) string {
	count := 0
	for _, i := range strings.Fields(input) {
		count++
		cmd = strings.Replace(cmd, "#"+strconv.Itoa(count), i, -1)
	}
	return cmd
}

func isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdin.Fd()))
}

func main() {
	app.Version("0.0.7")
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

	if isTerminal() {
		for i := 0; i < *circle; i++ {
			wg.Add(1)
			go run(*rawCmd, ch)
		}
	} else {
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadString('\n')
			if err != nil && err == io.EOF {
				break
			}
			for i := 0; i < *circle; i++ {
				cmd := replaceCmd(*rawCmd, input)
				wg.Add(1)
				go run(cmd, ch)
			}
		}
	}
	wg.Wait()
	close(ch)
}
