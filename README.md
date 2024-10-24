# simple-benchmark

## Caution

<span style="color:red">The command you're benchmarking will be executed in a shell. Make sure you trust the command you're benchmarking.</span>

## Description

Simple benchmarking tool in go for comparing the runtime of a bash command.

The command will be executed: 1, 10, 100, 1000 and 10.000 times.

## Usage

```bash
simple-benchmark "<bash command>"
```

## Installation

```bash
git clone https://github.com/oceanspy/simple-benchmark.git
cd simple-benchmark
go build -o simple-benchmark simple-benchmark.go
sudo ln -s {your_folder}/simple-benchmark/simple-benchmark /usr/local/bin/simple-benchmark
```
