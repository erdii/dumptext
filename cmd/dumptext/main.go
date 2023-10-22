package main

import (
	"debug/elf"
	"fmt"
	"os"
)

type formatType string

const (
	envFormat                  = "FORMAT"
	envFormatRaw    formatType = "raw" // default
	envFormatEscape formatType = "escape"
	envFormatDump   formatType = "dump"

	envValidate = "VALIDATE"
	envToggleOn = "1"

	textSection = ".text"
)

func main() {
	isValidateMode := os.Getenv(envValidate) == envToggleOn
	format := formatType(os.Getenv(envFormat))
	if format == "" || (format != envFormatRaw &&
		format != envFormatEscape &&
		format != envFormatDump) {
		format = envFormatEscape
	}

	// Read ELF binary from either stdin or given file path.
	var (
		f   *elf.File
		err error
	)
	if stdinIsPipe() {
		f, err = elf.NewFile(NewBufferingReaderAt(os.Stdin))
	} else if len(os.Args) == 2 && os.Args[1] != "" {
		f, err = elf.Open(os.Args[1])
	} else {
		printUsageAndExit()
	}
	must(err)

	section := f.Section(textSection)
	if section == nil {
		panic(".text section not found")
	}

	data, err := section.Data()
	must(err)

	if isValidateMode {
		violations := []int{}
		for i, v := range data {
			if v == 0 {
				violations = append(violations, i)
			}
		}
		if len(violations) != 0 {
			fmt.Fprintf(os.Stderr, "Null bytes found at indices: %d\n", violations)
			os.Exit(1)
		}
	}

	var wt writerTo
	switch format {
	case envFormatRaw:
		wt = nativeEndianWriterTo{}
	case envFormatEscape:
		wt = escapedHexBytesWriterTo{}
	case envFormatDump:
		wt = hexdumpWriterTo{}
	}
	must(wt.WriteTo(os.Stdout, data))
}

func printUsageAndExit() {
	exe, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve process executable: %s\n", err)
		fmt.Println("Usage: $exe path/to/elf/binary")
		fmt.Println("Usage: cat path/to/elf/binary | $exe")
	} else {
		fmt.Printf("Usage: %s path/to/elf/binary\n", exe)
		fmt.Printf("Usage: cat path/to/elf/binary | %s\n", exe)
	}
	fmt.Println("Description: Reads the .text section from an ELF binary and dumps (optionally formatted) bytes to stdout.")
	fmt.Println("Flags:")
	fmt.Printf("\tEnvvar: %s=%s(default)|%s|%s - specifies output format.\n", envFormat, envFormatEscape, envFormatDump, envFormatRaw)
	fmt.Printf("\tEnvvar: %s=%s - validates that there are no null bytes in the data.\n", envValidate, envToggleOn)
	os.Exit(1)
}

func stdinIsPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	return stat.Mode()&os.ModeCharDevice == 0
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
