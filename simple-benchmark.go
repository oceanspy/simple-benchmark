package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

var forbiddenCmd = []string{"rm", "mv", "cp", "dd", "rmdir", "touch", "ln", "chmod", "chown"}

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
		iterationDuration := time.Since(iterationStart)

		fmt.Printf("¤")
		if i+1 > 1 && (i+1)%30 == 0 {
			fmt.Println()
		}

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

func ContainsForbiddenCmd(cmd string) bool {
	splitCmd := strings.Split(cmd, " ")

	for _, forbidden := range forbiddenCmd {
		for _, cmd := range splitCmd {
			if cmd == forbidden {
				return true
			}
		}
	}

	return false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print(Red)
		fmt.Println("Usage: simple-benchmark <bash command>")
		fmt.Print(Reset)
		os.Exit(1)
	}

	bashCmd := os.Args[1]

	forceCmd := false
	if len(os.Args) == 3 {
		if os.Args[2] == "force" {
			forceCmd = true
		}
	}

	if ContainsForbiddenCmd(bashCmd) && !forceCmd {
		fmt.Print(Red)
		fmt.Println("An harmful command has been detected in the bash command")
		fmt.Print(Reset)
		fmt.Println("Use 'simple-benchmark <bash command> force' to force the execution")
		os.Exit(0)
	}

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
