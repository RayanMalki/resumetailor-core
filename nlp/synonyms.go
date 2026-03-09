package nlp

import "strings"

// synonyms maps safe technology aliases to canonical forms so matching treats
// equivalent tokens as the same concept.
var synonyms = map[string]string{
	// Go / Golang
	"go":     "golang",
	"golang": "golang",
	// JavaScript / JS
	"js":         "javascript",
	"javascript": "javascript",
	// TypeScript / TS
	"ts":         "typescript",
	"typescript": "typescript",
	// PostgreSQL variants
	"postgres":   "postgresql",
	"postgresql": "postgresql",
	// Kubernetes / K8s
	"k8s":        "kubernetes",
	"kubernetes": "kubernetes",
	// C# / CSharp
	"csharp": "csharp",
	"dotnet": "dotnet",
	"aspnet": "aspnet",
	"cpp":    "cpp",
	"fsharp": "fsharp",
	// Continuous Integration / Continuous Deployment
	"ci":   "ci",
	"cd":   "cd",
	"cicd": "cicd",
	// ReactJS / React
	"reactjs": "react",
	"react":   "react",
	// NodeJS / Node
	"nodejs": "nodejs",
	"node":   "nodejs",
	// REST / RESTful
	"rest":    "rest",
	"restful": "rest",
	// Amazon Web Services
	"aws": "aws",
	// Google Cloud Platform
	"gcp": "gcp",
	// Infonuagique (French for cloud computing)
	"infonuagique": "cloud",
	"cloud":        "cloud",
	// Agile / Scrum
	"scrum": "scrum",
	"agile": "agile",
	// DevOps
	"devops": "devops",
	// Microservices
	"microservices": "microservice",
	"microservice":  "microservice",
	// API
	"api":  "api",
	"apis": "api",
	// Tests automatisés / automated tests
	"automatise":   "automatise",
	"automatises":  "automatise",
	"automatisee":  "automatise",
	"automatisees": "automatise",
	// Protect tech terms that naturally end in 's' from depluralization
	"jenkins": "jenkins",
	"redis":   "redis",
	"travis":  "travis",
	"atlas":   "atlas",
	"pandas":  "pandas",
	"keras":   "keras",
	"express": "express",
	// "postgres" already defined above
}

// Canonicalize maps a token to its canonical synonym form.
// It returns the canonical form if one exists, otherwise the original token.
func Canonicalize(token string) string {
	if canon, ok := synonyms[token]; ok {
		return canon
	}
	return token
}

// depluralize applies simple plural→singular normalization for English and French.
// Only strips trailing "-s". We never strip "-es" (2 chars) because French singulars
// keep the "e" (e.g. "techniques"→"technique", NOT "techniqu").
func depluralize(token string) string {
	n := len(token)
	if n < 4 {
		return token // too short to safely strip
	}

	// Don't touch known tech terms / synonyms — they're already canonical
	if _, ok := synonyms[token]; ok {
		return token
	}

	// Only strip trailing -s
	if !strings.HasSuffix(token, "s") {
		return token
	}
	// Don't strip if it ends in "ss" (e.g. "process", "class")
	if strings.HasSuffix(token, "ss") {
		return token
	}

	candidate := token[:n-1]
	if len(candidate) >= 3 {
		return candidate
	}
	return token
}
