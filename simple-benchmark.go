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

type BenchmarkResult struct {
	Duration    time.Duration
	MinDuration time.Duration
	MaxDuration time.Duration
	TotalMaxRSS int64
}

func BenchmarkBashCmd(bashCmd string, iterations int) BenchmarkResult {
	var benchmarkResult BenchmarkResult

	fmt.Println()
	fmt.Print(Yellow)
	fmt.Printf("Starting %d iteration", iterations)
	fmt.Print(Reset)
	fmt.Println()

	for i := 0; i < iterations; i++ {
		iterationStart := time.Now()
		cmd := exec.Command("bash", "-c", bashCmd)
		err := cmd.Run()
		if err != nil {
			fmt.Print(Red)
			fmt.Printf("Error running bash command %s: %v\n", bashCmd, err)
			fmt.Print(Reset)
			os.Exit(1)
		}
		fmt.Printf("¤")
		if i+1 > 1 && (i+1)%30 == 0 {
			fmt.Println()
		}
		iterationDuration := time.Since(iterationStart)

		if iterationDuration < benchmarkResult.MinDuration || benchmarkResult.MinDuration == 0 {
			benchmarkResult.MinDuration = iterationDuration
		}

		if iterationDuration > benchmarkResult.MaxDuration {
			benchmarkResult.MaxDuration = iterationDuration
		}

		benchmarkResult.Duration += iterationDuration
		benchmarkResult.TotalMaxRSS += cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss
	}

	return benchmarkResult
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
		benchmarkResult := BenchmarkBashCmd(bashCmd, iter)
		fmt.Println()
		fmt.Print(Blue)
		fmt.Print("Avg. MaxRSS     ", Reset)
		fmt.Printf("%d\n", benchmarkResult.TotalMaxRSS/int64(iter))
		fmt.Print(Magenta)
		fmt.Print("Time total      ", Reset)
		fmt.Printf("%v\n", benchmarkResult.Duration)
		fmt.Print(Red)
		fmt.Print("Time Max        ", Reset)
		fmt.Printf("%v\n", benchmarkResult.MaxDuration)
		fmt.Print(Cyan)
		fmt.Print("Time Avg.       ", Reset)
		fmt.Printf("%v\n", benchmarkResult.Duration/time.Duration(iter))
		fmt.Print(Green)
		fmt.Print("Time Min        ", Reset)
		fmt.Printf("%v\n", benchmarkResult.MinDuration)
		fmt.Println()
	}
}
