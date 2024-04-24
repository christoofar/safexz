// xzspeed is a command line utility for testing the performance of the safexz compression library.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/christoofar/safexz"
)

func help(s string) error {
	if len(s) > 0 {
		println(s)
		println()
	}
	println("xzspeed is a command line utility for testing the performance of the safexz compression library.")
	println("Usage: speedtest -i <input file>")
	println("Options:")
	println("  -i, --input <input file>  The path to the input file")
	println("  -h, --help                Prints the help message")
	return nil
}

func main() {
	var inputPath string
	var helponly bool
	flag.StringVar(&inputPath, "input", "", "The path to the input file")
	flag.StringVar(&inputPath, "i", "", "The path to the input file")
	flag.BoolVar(&helponly, "help", false, "Prints the help message")
	flag.BoolVar(&helponly, "h", false, "Prints the help message")
	flag.Parse()

	if helponly {
		help("")
		return
	}

	if inputPath == "" {
		help("No input file specified.")
		return
	}

	println("Starting compression with CompressionSimpleFast...")
	timestamp := time.Now()
	times := make(map[string]time.Duration)
	bytes := make(map[string]int64)
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionSimpleFast)
	times["CompressionSimpleFast"] = time.Since(timestamp)
	bytes["CompressionSimpleFast"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionSimple...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionSimple)
	times["CompressionSimple"] = time.Since(timestamp)
	bytes["CompressionSimple"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionSimpleBetter...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionSimpleBetter)
	times["CompressionSimpleBetter"] = time.Since(timestamp)
	bytes["CompressionSimpleBetter"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionSimpleMax...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionSimpleMax)
	times["CompressionSimpleMax"] = time.Since(timestamp)
	bytes["CompressionSimpleMax"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Compression complete.  Moving on to CompressionMultiFast...")
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionMultiFast)
	times["CompressionMultiFast"] = time.Since(timestamp)
	bytes["CompressionMultiFast"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionMulti...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionMulti)
	times["CompressionMulti"] = time.Since(timestamp)
	bytes["CompressionMulti"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionMultiBetter...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionMultiBetter)
	times["CompressionMultiBetter"] = time.Since(timestamp)
	bytes["CompressionMultiBetter"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionMultiMax...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionMultiMax)
	times["CompressionMultiMax"] = time.Since(timestamp)
	bytes["CompressionMultiMax"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Compression complete.  Moving on to CompressionFullPowerFast...")
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionFullPowerFast)
	times["CompressionFullPowerFast"] = time.Since(timestamp)
	bytes["CompressionFullPowerFast"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionFullPower...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionFullPower)
	times["CompressionFullPower"] = time.Since(timestamp)
	bytes["CompressionFullPower"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionFullPowerBetter...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionFullPowerBetter)
	times["CompressionFullPowerBetter"] = time.Since(timestamp)
	bytes["CompressionFullPowerBetter"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Starting compression with CompressionFullPowerMax...")
	timestamp = time.Now()
	safexz.CompressFile(inputPath, inputPath+".xz", safexz.CompressionFullPowerMax)
	times["CompressionFullPowerMax"] = time.Since(timestamp)
	bytes["CompressionFullPowerMax"] = getFileSize(inputPath + ".xz")
	os.Remove(inputPath + ".xz")
	println("Compression complete.")

	println("Compression Results:")
	println(fmt.Sprintf("%-30s : %20s : %s", "Algorithm", "Time", "Size"))
	println(fmt.Sprintf("%-30s : %20s : %s", "---------", "----", "----"))
	for algo, duration := range times {
		println(fmt.Sprintf("%-30s : %20s : %d bytes", algo, duration.String(), bytes[algo]))
	}
}

func getFileSize(path string) int64 {
	info, _ := os.Stat(path)
	return info.Size()
}
