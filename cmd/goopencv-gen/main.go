package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brainplusplus/go-opencv/internal/contract"
	"github.com/brainplusplus/go-opencv/internal/generator"
)

func main() {
	format := flag.String("format", "abi-md", "output format: abi-md or cpp")
	out := flag.String("out", "", "output path; stdout when empty")
	flag.Parse()

	var data string
	switch *format {
	case "abi-md":
		data = generator.RenderABIMarkdown(contract.Modules)
	case "cpp":
		data = generator.RenderCPPSkeleton(contract.Modules)
	default:
		fatalf("unknown format %q", *format)
	}

	if *out == "" {
		fmt.Print(data)
		return
	}
	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		fatalf("create output dir: %v", err)
	}
	if err := os.WriteFile(*out, []byte(data), 0o644); err != nil {
		fatalf("write %s: %v", *out, err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "goopencv-gen: "+format+"\n", args...)
	os.Exit(1)
}
