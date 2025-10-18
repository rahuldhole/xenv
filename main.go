package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rahuldhole/xenv/app/processor"
)

// Version is set during build via ldflags
var Version = "dev"

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	for _, a := range os.Args[1:] {
		switch a {
		case "-h", "--help":
			printHelp()
			return
		case "-v", "--version":
			fmt.Println(Version)
			return
		}
	}

	for _, a := range os.Args[1:] {
		if a == "-h" || a == "--help" {
			printHelp()
			return
		}
		if a == "-v" || a == "--version" {
			fmt.Println(Version)
			return
		}
	}

	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <form-file> [-o|--output <file>] [-d|--defaults] [-r|--run-scripts] [-m|--merge] [-f|--force] [-h|--help]\n", os.Args[0])
		os.Exit(1)
	}

	formFile := os.Args[1]
	outputFile := ""
	defaultsMode := false
	allScripts := false
	preMerge := false
	forceOverwrite := false

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-o", "--output":
			if i+1 < len(os.Args) {
				outputFile = os.Args[i+1]
				i++
			} else {
				fmt.Println("Error: -o/--output requires a value")
				os.Exit(1)
			}
		case "-d", "--defaults":
			defaultsMode = true
		case "-r", "--run-scripts":
			allScripts = true
		case "-m", "--merge":
			preMerge = true
		case "-f", "--force":
			forceOverwrite = true
		case "-h", "--help":
			printHelp()
			return
		default:
			fmt.Printf("Unknown flag: %s\n\n", arg)
			printHelp()
			os.Exit(1)
		}
	}

	if outputFile == "" {
		outputFile = processor.DetermineOutputFile(formFile)
	}

	if _, err := os.Stat(formFile); os.IsNotExist(err) {
		fmt.Printf("Error: Form file '%s' not found.\n", formFile)
		os.Exit(1)
	}

	if preMerge && forceOverwrite {
		fmt.Println("Error: Cannot use both --merge (-m) and --force (-f) together.")
		os.Exit(1)
	}

	mergeMode := false
	
	if _, err := os.Stat(outputFile); err == nil {
		if forceOverwrite {
			fmt.Println("Overwrite (force) selected.")
		} else if preMerge {
			mergeMode = true
			fmt.Println("Merge (flag) selected. Existing values will be used as 'current' and conflicts shown.")
		} else {
			fmt.Printf("Output file '%s' already exists. Overwrite, merge, or cancel? [y/N/m (merge)]: ", outputFile)
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			switch response {
			case "y", "yes":
				fmt.Println("Overwrite selected.")
			case "m", "merge":
				mergeMode = true
				fmt.Println("Merge selected. Existing values will be used as 'current' and conflicts shown.")
			default:
				fmt.Println("Operation cancelled.")
				os.Exit(0)
			}
		}
	}

	hasDSL := processor.CheckForDSL(formFile)
	fmt.Printf("Interactive configuration for %s\n", filepath.Base(formFile))
	fmt.Println(strings.Repeat("-", 50))

	outputLines, err := processor.ProcessFormFile(formFile, outputFile, hasDSL, mergeMode, defaultsMode, allScripts)
	if err != nil {
		fmt.Printf("Error processing form file: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, []byte(strings.Join(outputLines, "\n")+"\n"), 0644)
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("Configuration saved to %s\n", outputFile)
}

func printHelp() {
	fmt.Printf("xenv %s - interactive/automated environment file generator\n", Version)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  xenv <form-file> [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -o, --output <file>   Write to specific output file (default: dot-prefixed template name)")
	fmt.Println("  -d, --defaults        Use existing / template defaults without interactive prompts")
	fmt.Println("  -r, --run-scripts     Run all inline scripts automatically (no confirmations)")
	fmt.Println("  -m, --merge           Merge with existing output (preserve unknown keys, show conflicts)")
	fmt.Println("  -f, --force           Overwrite existing output file without prompting")
	fmt.Println("  -v, --version         Show version and exit")
	fmt.Println("  -h, --help            Show this help and exit")
	fmt.Println()
	fmt.Println("Rules:")
	fmt.Println("  * --merge and --force are mutually exclusive.")
	fmt.Println("  * If neither --merge nor --force is provided and the output file exists, you will be prompted.")
	fmt.Println("  * Combine --defaults with --run-scripts for a fully automated generation including scripts.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  xenv config.xenv # You can use any file extension (.xenv, .template, .example, etc.)")
	fmt.Println("  xenv config.xenv -d -f")
	fmt.Println("  xenv config.xenv -m -r")
	fmt.Println("  xenv config.xenv -o .env.production -d -r")
	fmt.Println()
	fmt.Println("Inline Scripts:")
	fmt.Println("  Add script=\"...\" or script=`...` to a directive (e.g. @text, @button).")
	fmt.Println("  Use -r / --run-scripts to auto-run all scripts (including in --defaults mode).")
	fmt.Println()
}
