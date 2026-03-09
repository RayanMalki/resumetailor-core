# resumetailor-core

A Go library and CLI tool for scoring resume-to-job-description compatibility using BM25 + corpus IDF. It identifies high-signal job terms, detects missing keywords, and reports category coverage — giving you actionable feedback to improve ATS pass rates.

## Features

- **BM25 scoring** with a pre-computed IDF table built from a large corpus of job descriptions
- **Discipline detection** — automatically classifies resumes and job postings into engineering disciplines (IT/Software, Mechanical, Electrical, Aerospace, Industrial/Logistics)
- **Category coverage** — reports keyword coverage broken down by bucket (languages, cloud/devops/db, practices, soft skills, etc.)
- **Phrase-aware matching** — multi-word phrases like "machine learning" are matched in addition to individual tokens
- **Low-signal filtering** — generic terms (e.g. "experience", "requirements") are excluded from scoring to reduce noise
- **Discipline-specific profiles** — each discipline has its own canonical terms, synonyms, stopwords, and bucket weights

## Installation

```bash
go install github.com/RayanMalki/resumetailor-core/cmd/resumetailor@latest
```

Or clone and build:

```bash
git clone https://github.com/RayanMalki/resumetailor-core.git
cd resumetailor-core
go build ./cmd/resumetailor
```

## CLI Usage

```
resumetailor analyze --resume <file> --job <file> [--discipline <name>] [--json]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--resume` | Path to your resume text file (required) |
| `--job` | Path to the job description text file (required) |
| `--discipline` | Override auto-detected discipline. Options: `it_software`, `mechanical`, `electrical`, `aerospace`, `industrial_logistics` |
| `--json` | Output full `Signals` struct as JSON instead of human-readable text |

**Example:**

```bash
resumetailor analyze --resume resume.txt --job job.txt
```

```
ATS Score:       72%  (raw coverage: 52%)
Discipline:      it_software  (auto-detected, confidence: 0.84)

Top Job Terms:
  kubernetes           3.41  [cloud_devops_db]
  terraform            2.87  [cloud_devops_db]
  python               2.64  [languages]
  ...

Missing Terms:
  helm                 2.11
  observability        1.93
  ...

Overlap Terms:
  api, aws, docker, postgresql, ...

Category Coverage:
  cloud_devops_db:          58%
  languages:                80%
  practices:                67%
```

## Library Usage

```go
import (
    "github.com/RayanMalki/resumetailor-core/scoring"
    "github.com/RayanMalki/resumetailor-core/profiles"
    "github.com/RayanMalki/resumetailor-core/scoring/classifier"
)

// Auto-detect discipline and score
signals, err := scoring.Compute(resumeText, jobText)

// Or detect discipline explicitly, then score with its profile
detection := classifier.Detect(resumeText, jobText)
profile := profiles.Get(detection.Discipline)
signals, err := scoring.ComputeWithProfile(resumeText, jobText, profile)

fmt.Printf("ATS Score: %.0f%%\n", signals.Score * 100)
fmt.Printf("Discipline: %s\n", signals.Discipline)
```

### `Signals` struct

| Field | Description |
|-------|-------------|
| `Score` | Normalized ATS compatibility score in [0, 1] (sqrt-scaled for intuitive range) |
| `RawCoverage` | Raw matched/total weighted keyword coverage ratio |
| `BM25Score` | Raw BM25 score (diagnostic only) |
| `TopJobTerms` | Top high-signal terms from the job description |
| `MissingJobTerms` | High-signal job terms absent from the resume |
| `OverlapTerms` | Terms present in both documents |
| `CategoryCoverage` | Per-bucket coverage ratios |
| `Discipline` | Detected or overridden discipline |
| `ProfileVersion` | Version of the discipline profile used |

## Supported Disciplines

| Discipline | Key Signal |
|------------|------------|
| `it_software` | Programming languages, cloud, devops, databases |
| `mechanical` | CAD, FEA, GD&T, materials, manufacturing |
| `electrical` | PCB, embedded, firmware, power systems |
| `aerospace` | FAA, DO-178, avionics, propulsion, stress analysis |
| `industrial_logistics` | Supply chain, ERP, lean, logistics, warehousing |

## Architecture

```
resumetailor-core/
├── cmd/resumetailor/   # CLI entrypoint
├── corpus/             # Pre-computed IDF table and phrase list
├── nlp/                # Tokenizer, normalizer, stopwords, synonyms
├── profiles/           # Discipline profiles (terms, weights, buckets)
└── scoring/
    ├── bm25.go         # Core BM25 + IDF scoring engine
    ├── coverage.go     # Score normalization
    ├── lowsignal.go    # Low-signal term filtering
    └── classifier/     # Discipline auto-detection
```

## How Scoring Works

1. **Tokenization** — Text is lowercased, accent-stripped, and tokenized using a discipline-aware tokenizer that applies synonyms (e.g. `golang` → `go`) and removes stopwords.
2. **Phrase extraction** — Multi-word phrases from a curated corpus list are matched in the normalized text.
3. **IDF weighting** — Each job term is weighted by its IDF from a pre-built corpus table rather than document-level IDF (which collapses to binary in a 2-document setup).
4. **Profile weighting** — Discipline profiles amplify terms in high-priority buckets and suppress low-signal terms.
5. **Coverage score** — `matched_weight / total_weight`, then sqrt-scaled into [0, 1].

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

## License

MIT
