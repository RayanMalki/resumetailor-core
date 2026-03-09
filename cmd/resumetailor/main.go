package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/RayanMalki/resumetailor-core/profiles"
	"github.com/RayanMalki/resumetailor-core/scoring"
	"github.com/RayanMalki/resumetailor-core/scoring/classifier"
)

func main() {
	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	resumeFile := analyzeCmd.String("resume", "", "Path to resume text file (required)")
	jobFile := analyzeCmd.String("job", "", "Path to job description text file (required)")
	disciplineFlag := analyzeCmd.String("discipline", "", "Override discipline (e.g. it_software, mechanical, electrical, aerospace, industrial_logistics)")
	jsonOutput := analyzeCmd.Bool("json", false, "Output raw Signals struct as JSON")

	if len(os.Args) < 2 || os.Args[1] != "analyze" {
		fmt.Fprintf(os.Stderr, "Usage: resumetailor analyze --resume <file> --job <file> [--discipline <name>] [--json]\n")
		os.Exit(1)
	}

	if err := analyzeCmd.Parse(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *resumeFile == "" || *jobFile == "" {
		fmt.Fprintln(os.Stderr, "Error: --resume and --job are required")
		analyzeCmd.Usage()
		os.Exit(1)
	}

	resumeBytes, err := os.ReadFile(*resumeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading resume: %v\n", err)
		os.Exit(1)
	}
	jobBytes, err := os.ReadFile(*jobFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading job description: %v\n", err)
		os.Exit(1)
	}

	resumeText := string(resumeBytes)
	jobText := string(jobBytes)

	// Detect or override discipline
	var overrideDiscipline *profiles.Discipline
	if *disciplineFlag != "" {
		d, ok := profiles.ParseDiscipline(*disciplineFlag)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown discipline %q. Valid options: it_software, mechanical, electrical, aerospace, industrial_logistics\n", *disciplineFlag)
			os.Exit(1)
		}
		overrideDiscipline = &d
	}

	detection := classifier.Resolve(resumeText, jobText, overrideDiscipline)
	profile := profiles.Get(detection.Discipline)

	signals, err := scoring.ComputeWithProfile(resumeText, jobText, profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error computing score: %v\n", err)
		os.Exit(1)
	}

	if *jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(signals); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Human-readable output
	fmt.Printf("ATS Score:       %.0f%%  (raw coverage: %.0f%%)\n", signals.Score*100, signals.RawCoverage*100)
	fmt.Printf("Discipline:      %s", signals.Discipline)
	if detection.Source == profiles.SourceAuto {
		fmt.Printf("  (auto-detected, confidence: %.2f)", detection.Confidence)
	} else {
		fmt.Printf("  (user override)")
	}
	fmt.Println()
	fmt.Println()

	if len(signals.TopJobTerms) > 0 {
		fmt.Println("Top Job Terms:")
		for _, ts := range signals.TopJobTerms {
			cat := ""
			if ts.Category != "" {
				cat = fmt.Sprintf("  [%s]", ts.Category)
			}
			fmt.Printf("  %-20s %.2f%s\n", ts.Term, ts.Score, cat)
		}
		fmt.Println()
	}

	if len(signals.MissingJobTerms) > 0 {
		fmt.Println("Missing Terms:")
		for _, ts := range signals.MissingJobTerms {
			fmt.Printf("  %-20s %.2f\n", ts.Term, ts.Score)
		}
		fmt.Println()
	}

	if len(signals.OverlapTerms) > 0 {
		fmt.Printf("Overlap Terms:\n  %s\n\n", joinStrings(signals.OverlapTerms, ", "))
	}

	if len(signals.CategoryCoverage) > 0 {
		fmt.Println("Category Coverage:")
		for bucket, cov := range signals.CategoryCoverage {
			fmt.Printf("  %-25s %.0f%%\n", bucket+":", cov*100)
		}
	}
}

func joinStrings(ss []string, sep string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
