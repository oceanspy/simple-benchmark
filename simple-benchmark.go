package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func BenchmarkBashCmd(bashCmd string, iterations int) time.Duration {
	start := time.Now()
	var maxRSS int64

	fmt.Println()
	fmt.Print(Yellow)
	fmt.Printf("Starting %d iteration", iterations)
	fmt.Print(Reset)
	fmt.Println()

	for i := 0; i < iterations; i++ {
		cmd := exec.Command("bash", "-c", bashCmd)
		err := cmd.Run()
		if err != nil {
			fmt.Print(Red)
			fmt.Printf("Error running bash command %s: %v\n", bashCmd, err)
			fmt.Print(Reset)
			return 0
		}
		fmt.Printf("¤")
		if i+1 > 1 && (i+1)%30 == 0 {
			fmt.Println()
		}

		maxRSS += cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss
	}

	avgMaxRSS := maxRSS / int64(iterations)
	fmt.Println()
	fmt.Print(Cyan, "Avg. MaxRSS: ", Reset)
	fmt.Printf("%d\n", avgMaxRSS)
	fmt.Print(Reset)

	return time.Since(start)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(Red)
		fmt.Println("Usage: simple-benchmark <bash command>")
		fmt.Print(Reset)
		os.Exit(1)
	}

	bashCmd := os.Args[1]

	iterations := []int{1, 10, 100, 1000, 10000}

	for _, iter := range iterations {
		duration := BenchmarkBashCmd(bashCmd, iter)
		fmt.Print(Magenta, "Time: ", Reset)
		fmt.Printf("%v\n", duration)
		fmt.Print(Reset)
		fmt.Print(Green, "Avg.: ", Reset)
		fmt.Printf("%v\n", duration/time.Duration(iter))
		fmt.Println()
	}
}