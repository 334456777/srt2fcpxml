package main

import (
	"fmt"
	"os"
	"path/filepath"
	"srt2fcpxml/core"
	"strconv"
	"strings"

	"github.com/asticode/go-astisub"
)

var version = "dev" // Set via ldflags at build time

func main() {
	args := os.Args[1:] // Get all arguments except program name

	// 检查版本标志
	if len(args) == 1 && (args[0] == "--version" || args[0] == "-v") {
		fmt.Printf("srt2fcpxml version %s\n", version)
		os.Exit(0)
	}

	// 默认参数
	var srtFile string
	var frameDuration interface{} = 30
	width := 1920
	height := 1080

	// 根据参数数量解析不同的模式
	switch len(args) {
	case 0:
		// ./srt2fcpxml - Auto find SRT file in current directory with default parameters
		var err error
		srtFile, err = findSrtFileInCurrentDir()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			printUsage()
			os.Exit(1)
		}
		frameDuration = 30

	case 1:
		// ./srt2fcpxml 60 - Use specified frame rate
		frameRate, err := parseFrameRate(args[0])
		if err != nil {
			fmt.Printf("Error parsing frame rate: %v\n", err)
			printUsage()
			os.Exit(1)
		}
		frameDuration = frameRate

		srtFile, err = findSrtFileInCurrentDir()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			printUsage()
			os.Exit(1)
		}
	case 2:
		// ./srt2fcpxml 1920 1080 - Use specified resolution with default frame rate
		var err error
		width, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error parsing width: %v\n", err)
			printUsage()
			os.Exit(1)
		}

		height, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error parsing height: %v\n", err)
			printUsage()
			os.Exit(1)
		}
		frameDuration = 30

		srtFile, err = findSrtFileInCurrentDir()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			printUsage()
			os.Exit(1)
		}

	case 3:
		// ./srt2fcpxml 1920 1080 30 - Use specified resolution and frame rate
		var err error
		width, err = strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Error parsing width: %v\n", err)
			printUsage()
			os.Exit(1)
		}

		height, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Error parsing height: %v\n", err)
			printUsage()
			os.Exit(1)
		}

		frameDuration, err = parseFrameRate(args[2])
		if err != nil {
			fmt.Printf("Error parsing frame rate: %v\n", err)
			printUsage()
			os.Exit(1)
		}

		srtFile, err = findSrtFileInCurrentDir()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			printUsage()
			os.Exit(1)
		}

	default:
		fmt.Println("Error: Invalid number of arguments")
		printUsage()
		os.Exit(1)
	}

	// Open SRT file
	f, err := astisub.OpenFile(srtFile)
	if err != nil {
		fmt.Printf("Error opening SRT file %s: %v\n", srtFile, err)
		os.Exit(1)
	}

	// Generate XML output
	out := `<?xml version="1.0" encoding="UTF-8" ?>
	<!DOCTYPE fcpxml>
	
	`

	project, path := getPath(srtFile)
	result, err := core.Srt2FcpXmlExport(project, frameDuration, f, width, height)
	if err != nil {
		fmt.Printf("Error generating FCPXML: %v\n", err)
		os.Exit(1)
	}

	out += string(result)
	targetFile := fmt.Sprintf("%s/%s.fcpxml", path, project)

	fd, err := os.Create(targetFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer fd.Close()

	_, err = fd.Write([]byte(out))
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s --> %s\n", srtFile, targetFile)
	fmt.Printf("Settings: %dx%dp %v fps\n", width, height, frameDuration)
}

func getPath(filePath string) (projectName, targetPath string) {
	path, _ := filepath.Abs(filePath)
	parts := strings.Split(path, "/")
	projectName = func(file string) string {
		parts := strings.Split(file, ".")
		return strings.Join(parts[0:len(parts)-1], ".")
	}(parts[len(parts)-1])
	targetPath = func(parts []string) string {
		return strings.Join(parts, "/")
	}(parts[0 : len(parts)-1])
	return
}

// findSrtFileInCurrentDir finds SRT files in current directory
func findSrtFileInCurrentDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	files, err := os.ReadDir(currentDir)
	if err != nil {
		return "", fmt.Errorf("failed to read current directory: %v", err)
	}

	var srtFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".srt") {
			srtFiles = append(srtFiles, file.Name())
		}
	}

	if len(srtFiles) == 0 {
		return "", fmt.Errorf("no SRT files found in current directory")
	}

	if len(srtFiles) == 1 {
		return filepath.Join(currentDir, srtFiles[0]), nil
	}

	// If multiple SRT files exist, return the first one
	fmt.Printf("Found multiple SRT files, using: %s\n", srtFiles[0])
	return filepath.Join(currentDir, srtFiles[0]), nil
}

// parseFrameRate parses frame rate parameter
func parseFrameRate(frameRateStr string) (interface{}, error) {
	// Supported frame rates: 23.98, 24, 25, 29.97, 30, 50, 59.94, 60
	supportedRates := map[string]interface{}{
		"23.98": 23.98,
		"24":    24,
		"25":    25,
		"29.97": 29.97,
		"30":    30,
		"50":    50,
		"59.94": 59.94,
		"60":    60,
	}

	if rate, exists := supportedRates[frameRateStr]; exists {
		return rate, nil
	}

	// Try parsing as float
	if strings.Contains(frameRateStr, ".") {
		rate, err := strconv.ParseFloat(frameRateStr, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid frame rate format: %s", frameRateStr)
		}
		return rate, nil
	}

	// Try parsing as integer
	rate, err := strconv.Atoi(frameRateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid frame rate format: %s", frameRateStr)
	}

	return rate, nil
}

// printUsage prints usage instructions
func printUsage() {
	fmt.Printf("srt2fcpxml version %s\n", version)
	fmt.Println("Convert SRT subtitle files to Final Cut Pro XML format")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ./srt2fcpxml                    - Auto find SRT file in current directory, use 1920x1080@30fps")
	fmt.Println("  ./srt2fcpxml <framerate>        - Auto find SRT file, use specified framerate with 1920x1080")
	fmt.Println("  ./srt2fcpxml <width> <height> - Auto find SRT file, use specified resolution with 30fps")
	fmt.Println("  ./srt2fcpxml <width> <height> <framerate> - Auto find SRT file, use specified resolution and framerate")
	fmt.Println("  ./srt2fcpxml --version | -v     - Show version information")
	fmt.Println("")
	fmt.Println("Supported frame rates: 23.98, 24, 25, 29.97, 30, 50, 59.94, 60")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  ./srt2fcpxml                    - Convert with 1920x1080@30fps")
	fmt.Println("  ./srt2fcpxml 60                 - Convert with 1920x1080@60fps")
	fmt.Println("  ./srt2fcpxml 1920 1080          - Convert with 1920x1080@30fps")
	fmt.Println("  ./srt2fcpxml 1920 1080 29.97    - Convert with 1920x1080@29.97fps")
}
