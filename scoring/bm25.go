package scoring

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/RayanMalki/resumetailor-core/corpus"
	"github.com/RayanMalki/resumetailor-core/nlp"
	"github.com/RayanMalki/resumetailor-core/profiles"

	libbm25 "github.com/crawlab-team/bm25"
)

const (
	defaultK1   = 1.5
	defaultB    = 0.75
	defaultTopN = 10
)

// TermScore represents a term with its BM25-derived score.
type TermScore struct {
	Term     string  `json:"term"`
	Score    float64 `json:"score"`
	Category string  `json:"category,omitempty"`
}

// Signals exposes BM25-derived signals for reporting.
type Signals struct {
	TopJobTerms     []TermScore `json:"top_job_terms"`
	MissingJobTerms []TermScore `json:"missing_job_terms"`
	OverlapTerms    []string    `json:"overlap_terms"`
	// BucketedTopTerms groups high-signal terms into deterministic categories.
	BucketedTopTerms map[string][]TermScore `json:"bucketed_top_terms,omitempty"`
	// LowSignalTerms are filtered from top/missing displays to reduce noise.
	LowSignalTerms []TermScore `json:"low_signal_terms,omitempty"`
	// CategoryCoverage reports matched/total weighted coverage per bucket.
	CategoryCoverage map[string]float64 `json:"category_coverage,omitempty"`
	// Discipline is the selected discipline for this scoring pass.
	Discipline string `json:"discipline,omitempty"`
	// DisciplineEvidence are the matched terms that drove discipline selection.
	DisciplineEvidence []TermScore `json:"discipline_evidence,omitempty"`
	// ProfileVersion identifies the profile dictionary version used for scoring.
	ProfileVersion string `json:"profile_version,omitempty"`
	// Score is normalized IDF-weighted keyword coverage in [0,1].
	// Derived from RawCoverage via sqrt transformation for intuitive scaling.
	Score float64 `json:"score"`
	// RawCoverage is the un-normalized matched/total weight ratio for diagnostics.
	RawCoverage float64 `json:"raw_coverage"`
	// BM25Score is the raw 2-document BM25 score (resume scored against job query).
	// It is exposed for diagnostics only and is not used as ATS compatibility score.
	BM25Score float64 `json:"bm25_score"`
}

// Compute calculates BM25 signals for resume and job text matching.
//
// IDF values come from a pre-computed static table built from a large corpus
// of job descriptions, replacing the previous single-document IDF calculation
// which produced meaningless binary scores.
func Compute(resumeText, jobText string) (Signals, error) {
	return ComputeWithProfile(resumeText, jobText, profiles.Get(profiles.DefaultDiscipline()))
}

func ComputeWithProfile(resumeText, jobText string, profile profiles.Profile) (Signals, error) {
	if profile.Discipline == "" {
		profile = profiles.Get(profiles.DefaultDiscipline())
	}
	tokenizer := func(text string) []string {
		return nlp.TokenizeWithProfile(text, profile)
	}

	resumeTokens := tokenizer(resumeText)
	jobTokens := tokenizer(jobText)

	if len(resumeTokens) == 0 || len(jobTokens) == 0 {
		return Signals{
			Discipline:     string(profile.Discipline),
			ProfileVersion: profile.Version,
		}, nil
	}

	// We still use the library for the overall document score, but now pass
	// both resume and job as separate corpus documents so the library's
	// internal IDF has at least two documents to work with.
	bm25Instance, err := libbm25.NewBM25Okapi(
		[]string{resumeText, jobText},
		tokenizer, defaultK1, defaultB, nil,
	)
	if err != nil {
		return Signals{}, fmt.Errorf("bm25 init: %w", err)
	}

	scores, err := bm25Instance.GetScores(jobTokens)
	if err != nil {
		return Signals{}, fmt.Errorf("bm25 score: %w", err)
	}

	// The resume is document 0 in the corpus.
	bm25Score := 0.0
	if len(scores) > 0 {
		bm25Score = scores[0]
	}

	jobFreq := nlp.TermFreq(jobTokens)
	resumeFreq := nlp.TermFreq(resumeTokens)

	normalizedJob := nlp.NormalizeForTokenization(jobText)
	normalizedResume := nlp.NormalizeForTokenization(resumeText)
	for phrase, count := range extractPhraseCounts(normalizedJob) {
		jobFreq[phrase] += count
	}
	for phrase, count := range extractPhraseCounts(normalizedResume) {
		resumeFreq[phrase] += count
	}

	overlapTerms := make([]string, 0)
	missingTerms := make([]TermScore, 0)
	topTerms := make([]TermScore, 0, len(jobFreq))
	lowSignalList := make([]TermScore, 0)
	buckets := make(map[string][]TermScore, 5)
	bucketTotals := make(map[string]float64, len(profile.Buckets))
	bucketMatched := make(map[string]float64, len(profile.Buckets))
	totalWeight := 0.0
	matchedWeight := 0.0

	for term, qtf := range jobFreq {
		tf := resumeFreq[term]

		// Job-side importance is based on static corpus IDF and query frequency.
		termIDF := lookupIDF(term)
		importance := termIDF * float64(qtf)
		if meta, ok := profile.CanonicalTerms[term]; ok && meta.Weight > 0 {
			importance *= meta.Weight
		}
		category := bucketForTerm(profile, term)
		weightedImportance := importance * profiles.BucketWeight(profile, category)
		scored := TermScore{
			Term:     term,
			Score:    weightedImportance,
			Category: category,
		}

		// Filter generic/noisy terms out of top/missing and coverage math.
		if isLowSignalTerm(term, category, termIDF, profile) {
			lowSignalList = append(lowSignalList, scored)
			continue
		}

		totalWeight += weightedImportance
		bucketTotals[category] += weightedImportance

		if tf > 0 {
			overlapTerms = append(overlapTerms, term)
			matchedWeight += weightedImportance
			bucketMatched[category] += weightedImportance
		} else {
			missingTerms = append(missingTerms, scored)
		}

		topTerms = append(topTerms, scored)
		buckets[category] = append(buckets[category], scored)
	}

	sort.Strings(overlapTerms)
	sortTermScores(topTerms)
	sortTermScores(missingTerms)
	sortTermScores(lowSignalList)
	for category, terms := range buckets {
		sortTermScores(terms)
		if len(terms) > defaultTopN {
			terms = terms[:defaultTopN]
		}
		buckets[category] = terms
	}

	if len(topTerms) > defaultTopN {
		topTerms = topTerms[:defaultTopN]
	}

	rawCoverage := 0.0
	coverageScore := 0.0
	if totalWeight > 0 {
		rawCoverage = matchedWeight / totalWeight
		coverageScore = NormalizeCoverageScore(rawCoverage)
	}
	categoryCoverage := make(map[string]float64, len(bucketTotals))
	for bucket, total := range bucketTotals {
		if total <= 0 {
			continue
		}
		categoryCoverage[bucket] = bucketMatched[bucket] / total
	}

	return Signals{
		TopJobTerms:      topTerms,
		MissingJobTerms:  missingTerms,
		OverlapTerms:     overlapTerms,
		BucketedTopTerms: buckets,
		LowSignalTerms:   lowSignalList,
		CategoryCoverage: categoryCoverage,
		Discipline:       string(profile.Discipline),
		ProfileVersion:   profile.Version,
		Score:            coverageScore,
		RawCoverage:      rawCoverage,
		BM25Score:        bm25Score,
	}, nil
}

const (
	categoryLanguages   = "languages"
	categoryCloudDevOps = "cloud_devops_db"
	categoryPractices   = "practices"
	categorySoftSkills  = "soft_skills"
	categoryOther       = "other"
)

var (
	languageTerms = map[string]struct{}{
		"python": {}, "java": {}, "javascript": {}, "typescript": {}, "golang": {}, "csharp": {},
		"cpp": {}, "ruby": {}, "php": {}, "rust": {}, "kotlin": {}, "swift": {}, "scala": {},
		"sql": {}, "bash": {}, "powershell": {}, "r": {}, "perl": {}, "matlab": {},
	}
	cloudDevOpsDBTerms = map[string]struct{}{
		"aws": {}, "azure": {}, "gcp": {}, "cloud": {}, "devops": {}, "docker": {}, "kubernetes": {},
		"terraform": {}, "ansible": {}, "jenkins": {}, "helm": {}, "linux": {}, "nginx": {},
		"postgresql": {}, "mysql": {}, "mongodb": {}, "redis": {}, "dynamodb": {}, "snowflake": {},
		"bigquery": {}, "databricks": {}, "kafka": {}, "rabbitmq": {}, "prometheus": {}, "grafana": {},
		"nodejs": {}, "nextjs": {},
	}
	practiceTerms = map[string]struct{}{
		"agile": {}, "scrum": {}, "kanban": {}, "tdd": {}, "ddd": {}, "sre": {}, "ci": {},
		"cd": {}, "cicd": {}, "microservice": {}, "api": {}, "rest": {}, "graphql": {},
		"testing": {}, "test": {}, "automation": {}, "architecture": {}, "observability": {},
		"monitoring": {}, "reliability": {}, "security": {}, "performance": {},
	}
	softSkillTerms = map[string]struct{}{
		"communication": {}, "leadership": {}, "mentoring": {}, "collaboration": {}, "stakeholder": {},
		"ownership": {}, "initiative": {}, "teamwork": {}, "presentation": {}, "adaptability": {},
		"problem": {}, "problemsolving": {}, "problem-solving": {}, "coaching": {},
	}
	reSQLFamily = regexp.MustCompile(`.+sql$`)
)

func classifyTerm(term string) string {
	if _, ok := languageTerms[term]; ok {
		return categoryLanguages
	}
	if _, ok := cloudDevOpsDBTerms[term]; ok || reSQLFamily.MatchString(term) {
		return categoryCloudDevOps
	}
	if _, ok := practiceTerms[term]; ok {
		return categoryPractices
	}
	if _, ok := softSkillTerms[term]; ok {
		return categorySoftSkills
	}
	return categoryOther
}

// extractPhraseCounts scans normalizedText for all known phrases and returns
// their occurrence counts. normalizedText must already be lowercased/accent-stripped.
func extractPhraseCounts(normalizedText string) map[string]int {
	counts := make(map[string]int, len(corpus.PhraseIDF))
	for phrase := range corpus.PhraseIDF {
		if n := strings.Count(normalizedText, phrase); n > 0 {
			counts[phrase] = n
		}
	}
	return counts
}

func bucketForTerm(profile profiles.Profile, term string) string {
	if strings.Contains(term, " ") {
		if cat, ok := corpus.PhraseCategories[term]; ok {
			return cat
		}
		return categoryOther
	}
	bucket := profiles.BucketForTerm(profile, term)
	if bucket != "" && bucket != categoryOther {
		return bucket
	}
	if profile.Discipline == profiles.DisciplineITSoftware {
		return classifyTerm(term)
	}
	return categoryOther
}

// lookupIDF returns the IDF for a term, also trying the canonical synonym form.
func lookupIDF(term string) float64 {
	if v, ok := corpus.CorpusIDF[term]; ok {
		return v
	}
	if v, ok := corpus.PhraseIDF[term]; ok {
		return v
	}
	// Try the canonical form (e.g. "go" → "golang")
	if canon := nlp.Canonicalize(term); canon != term {
		if v, ok := corpus.CorpusIDF[canon]; ok {
			return v
		}
	}
	if corpus.IsLikelyNoisyToken(term) {
		return corpus.NoisyUnknownIDF
	}
	if len(term) <= 3 {
		return corpus.ShortUnknownIDF
	}
	return corpus.DefaultCorpusIDF
}

func sortTermScores(items []TermScore) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].Term < items[j].Term
		}
		return items[i].Score > items[j].Score
	})
}

