package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/noedaka/go-config-parser/internal/config"
	"github.com/noedaka/go-config-parser/internal/parser"
	"github.com/noedaka/go-config-parser/internal/service"
)

func main() {
	cfg := config.Init()

	var input []byte
	var err error

	if cfg.IsStdin {
		input, err = io.ReadAll(os.Stdin)
	} else {
		if flag.NArg() == 0 {
			log.Fatalf("File path not specified")
		}

		input, err = os.ReadFile(flag.Arg(0))
	}

	if err != nil {
		log.Fatalf("Reading error: %v\n", err)
	}

	data, err := parser.ParseConfig(input)
	if err != nil {
		log.Fatalf("Parsing error: %v\n", err)
	}

	rules := []service.Rule{
		service.DebugLogRule{},
		service.PlaintextPasswordRule{},
		service.ZeroHostRule{},
		service.TLSDisabledRule{},
		service.NewWeakAlgorithmRule(),
	}

	var issues []service.Issue
	for _, rule := range rules {
		issues = append(issues, rule.Check(data)...)
	}

	for _, issue := range issues {
		fmt.Printf("[%s] %s\nРекомендация: %s\n\n", issue.Severity, issue.Message, issue.Recommendation)
	}

	if len(issues) > 0 && !cfg.IsSilent {
		return
	}
}
