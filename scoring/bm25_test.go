package scoring

import (
	"testing"

	"github.com/RayanMalki/resumetailor-core/corpus"
	"github.com/RayanMalki/resumetailor-core/nlp"
	"github.com/RayanMalki/resumetailor-core/profiles"
)

func TestComputeEmptyInputs(t *testing.T) {
	got, err := Compute("", "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got.Score != 0 {
		t.Fatalf("expected score 0, got %v", got.Score)
	}
	if len(got.TopJobTerms) != 0 || len(got.MissingJobTerms) != 0 || len(got.OverlapTerms) != 0 {
		t.Fatalf("expected empty signals, got %+v", got)
	}
}

func TestComputeNoMissingTerms(t *testing.T) {
	resume := "Go developer with docker kubernetes"
	job := "Go developer docker kubernetes"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got.MissingJobTerms) != 0 {
		t.Fatalf("expected no missing terms, got %+v", got.MissingJobTerms)
	}
}

func TestComputeMissingTerms(t *testing.T) {
	resume := "Go developer"
	job := "Go developer kubernetes"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got.MissingJobTerms) == 0 {
		t.Fatalf("expected missing terms, got none")
	}
	if got.MissingJobTerms[0].Term != "kubernetes" {
		t.Fatalf("expected missing term kubernetes, got %+v", got.MissingJobTerms)
	}
}

func TestComputeDeterministicOrdering(t *testing.T) {
	resume := "alpha beta"
	job := "beta alpha"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(got.TopJobTerms) < 2 {
		t.Fatalf("expected at least 2 top terms, got %+v", got.TopJobTerms)
	}

	first := got.TopJobTerms[0].Term
	second := got.TopJobTerms[1].Term
	if first != "alpha" || second != "beta" {
		t.Fatalf("expected deterministic alpha/beta order, got %s/%s", first, second)
	}
}

// TestStaticIDFDifferentiatesTermImportance verifies that the static IDF table
// produces meaningfully different scores for rare vs common terms (the core fix
// for improvement #4).
func TestStaticIDFDifferentiatesTermImportance(t *testing.T) {
	resume := "python developer kubernetes docker machine learning tensorflow data science"
	job := "python developer kubernetes docker machine learning tensorflow data science"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.TopJobTerms) < 2 {
		t.Fatalf("expected multiple top terms, got %d", len(got.TopJobTerms))
	}

	// With static IDF, terms should have different scores. The old
	// single-document approach gave every present term the same IDF.
	first := got.TopJobTerms[0]
	last := got.TopJobTerms[len(got.TopJobTerms)-1]
	if first.Score == last.Score {
		t.Fatalf("static IDF should produce different scores for different terms, "+
			"but first (%s=%.4f) == last (%s=%.4f)",
			first.Term, first.Score, last.Term, last.Score)
	}

	// Rare terms (tensorflow, kubernetes) should score higher than common ones
	// (developer, data).
	t.Logf("top terms with static IDF:")
	for _, ts := range got.TopJobTerms {
		t.Logf("  %s: %.4f", ts.Term, ts.Score)
	}
}

func TestLookupIDFKnownTerm(t *testing.T) {
	val := lookupIDF("python")
	if val == corpus.DefaultCorpusIDF {
		t.Fatal("expected python to have a specific IDF, got default")
	}
	if val <= 0 {
		t.Fatalf("expected positive IDF for python, got %f", val)
	}
}

func TestLookupIDFUnknownTerm(t *testing.T) {
	val := lookupIDF("xyzzy_nonexistent_term_12345")
	if val != corpus.NoisyUnknownIDF {
		t.Fatalf("expected noisy unknown IDF %f, got %f", corpus.NoisyUnknownIDF, val)
	}
}

func TestWeightedCoveragePrioritizesImportantTerms(t *testing.T) {
	resume := "api api api"
	job := "api kubernetes"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// api IDF(2.5) vs kubernetes IDF(4.5) => raw coverage ≈ 2.5/7.0 ≈ 0.357
	// normalized via sqrt ≈ 0.597
	if got.RawCoverage >= 0.5 {
		t.Fatalf("expected raw coverage below 0.5, got %.4f", got.RawCoverage)
	}
	if got.RawCoverage <= 0.28 {
		t.Fatalf("expected raw coverage above 0.28, got %.4f", got.RawCoverage)
	}
	// Normalized score should be higher than raw
	if got.Score <= got.RawCoverage {
		t.Fatalf("expected normalized score > raw coverage, got score=%.4f raw=%.4f", got.Score, got.RawCoverage)
	}
	if got.Score >= 0.70 {
		t.Fatalf("expected normalized score below 0.70, got %.4f", got.Score)
	}
	if got.Score <= 0.50 {
		t.Fatalf("expected normalized score above 0.50, got %.4f", got.Score)
	}
}

func TestTopJobTermsIncludeMissingImportantKeywords(t *testing.T) {
	resume := "python engineer"
	job := "python kubernetes"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.TopJobTerms) < 2 {
		t.Fatalf("expected two top job terms, got %+v", got.TopJobTerms)
	}

	// kubernetes should outrank python by static IDF (4.5 > 4.2),
	// even though it's missing from resume.
	if got.TopJobTerms[0].Term != "kubernetes" {
		t.Fatalf("expected top term kubernetes, got %+v", got.TopJobTerms)
	}
	if got.MissingJobTerms[0].Term != "kubernetes" {
		t.Fatalf("expected missing term kubernetes, got %+v", got.MissingJobTerms)
	}
}

func TestTokenizeTechPunctuationVariants(t *testing.T) {
	resume := "Built APIs in C# and C++ on .NET with Node.js"
	job := "csharp cpp dotnet nodejs"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.MissingJobTerms) != 0 {
		t.Fatalf("expected no missing tech terms, got %+v", got.MissingJobTerms)
	}
}

func TestAmazonIsNotForcedToAWS(t *testing.T) {
	resume := "Improved delivery operations at Amazon"
	job := "aws"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.MissingJobTerms) == 0 || got.MissingJobTerms[0].Term != "aws" {
		t.Fatalf("expected aws to remain missing when only amazon is present, got %+v", got.MissingJobTerms)
	}
}

func TestLookupIDFShortUnknownAcronym(t *testing.T) {
	val := lookupIDF("gke")
	if val != corpus.ShortUnknownIDF {
		t.Fatalf("expected short unknown IDF %f, got %f", corpus.ShortUnknownIDF, val)
	}
}

func TestLookupIDFDefaultUnknownWord(t *testing.T) {
	val := lookupIDF("platformization")
	if val != corpus.DefaultCorpusIDF {
		t.Fatalf("expected default unknown IDF %f, got %f", corpus.DefaultCorpusIDF, val)
	}
}

func TestLowSignalTermsAreFilteredFromMissing(t *testing.T) {
	resume := "python kubernetes docker"
	job := "python kubernetes docker experience responsibilities years"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, missing := range got.MissingJobTerms {
		if missing.Term == "experience" || missing.Term == "responsibilities" || missing.Term == "years" {
			t.Fatalf("expected low-signal terms to be filtered, got %+v", got.MissingJobTerms)
		}
	}
}

func TestBucketedTopTermsClassifiesCoreSkills(t *testing.T) {
	resume := "python aws docker scrum communication"
	job := "python aws docker scrum communication"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got.BucketedTopTerms[categoryLanguages]) == 0 {
		t.Fatalf("expected language bucket to be populated, got %+v", got.BucketedTopTerms)
	}
	if len(got.BucketedTopTerms[categoryCloudDevOps]) == 0 {
		t.Fatalf("expected cloud/devops/db bucket to be populated, got %+v", got.BucketedTopTerms)
	}
	if len(got.BucketedTopTerms[categoryPractices]) == 0 {
		t.Fatalf("expected practices bucket to be populated, got %+v", got.BucketedTopTerms)
	}
	if len(got.BucketedTopTerms[categorySoftSkills]) == 0 {
		t.Fatalf("expected soft skills bucket to be populated, got %+v", got.BucketedTopTerms)
	}
}

func TestComputeWithProfileUsesDisciplineBuckets(t *testing.T) {
	profile := profiles.Get(profiles.DisciplineMechanical)
	resume := "solidworks fea cfd gdt"
	job := "solidworks fea cfd gdt tolerance"

	got, err := ComputeWithProfile(resume, job, profile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Discipline != string(profiles.DisciplineMechanical) {
		t.Fatalf("expected mechanical discipline, got %s", got.Discipline)
	}
	if len(got.BucketedTopTerms["design_tools"]) == 0 {
		t.Fatalf("expected design_tools bucket terms, got %+v", got.BucketedTopTerms)
	}
	if len(got.BucketedTopTerms["simulation_analysis"]) == 0 {
		t.Fatalf("expected simulation_analysis bucket terms, got %+v", got.BucketedTopTerms)
	}
}

func TestComputeWithProfileHidesLowSignalOtherTerms(t *testing.T) {
	profile := profiles.Get(profiles.DisciplineAerospace)
	resume := "student team player"
	job := "student team player do178 verification safety"

	got, err := ComputeWithProfile(resume, job, profile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, missing := range got.MissingJobTerms {
		if missing.Term == "student" || missing.Term == "team" || missing.Term == "player" {
			t.Fatalf("expected low-signal terms to be filtered, got %+v", got.MissingJobTerms)
		}
	}
	foundCritical := false
	for _, missing := range got.MissingJobTerms {
		if missing.Term == "do178" {
			foundCritical = true
			break
		}
	}
	if !foundCritical {
		t.Fatalf("expected high-signal term do178 to remain visible, got %+v", got.MissingJobTerms)
	}
}

func TestNormalizeCoverageScore(t *testing.T) {
	tests := []struct {
		name string
		raw  float64
		want float64
	}{
		{"zero", 0.0, 0.0},
		{"one", 1.0, 1.0},
		{"negative", -0.5, 0.0},
		{"above_one", 1.5, 1.0},
		{"sqrt_0.25", 0.25, 0.5},
		{"sqrt_0.64", 0.64, 0.8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeCoverageScore(tt.raw)
			if diff := got - tt.want; diff > 1e-9 || diff < -1e-9 {
				t.Fatalf("NormalizeCoverageScore(%f) = %f, want %f", tt.raw, got, tt.want)
			}
		})
	}
}

func TestAllCanonicalTermsHaveIDF(t *testing.T) {
	allProfiles := profiles.GetAll()
	if len(allProfiles) == 0 {
		t.Fatal("expected at least one profile from GetAll()")
	}

	var missing []string
	for discipline, p := range allProfiles {
		for term := range p.CanonicalTerms {
			if !hasExplicitIDF(term) {
				missing = append(missing, string(discipline)+":"+term)
			}
		}
	}
	if len(missing) > 0 {
		t.Fatalf("canonical terms with no IDF entry (falling to defaults): %v", missing)
	}
}

// hasExplicitIDF checks whether a term resolves to an explicit corpusIDF entry
// (directly or via canonicalize), as opposed to falling back to a default.
func hasExplicitIDF(term string) bool {
	if _, ok := corpus.CorpusIDF[term]; ok {
		return true
	}
	if canon := nlp.Canonicalize(term); canon != term {
		if _, ok := corpus.CorpusIDF[canon]; ok {
			return true
		}
	}
	return false
}

func TestExtractPhraseCounts(t *testing.T) {
	text := nlp.NormalizeForTokenization("experience with machine learning and deep learning systems")
	counts := extractPhraseCounts(text)
	if counts["machine learning"] != 1 {
		t.Fatalf("expected 'machine learning' count=1, got %d", counts["machine learning"])
	}
	if counts["deep learning"] != 1 {
		t.Fatalf("expected 'deep learning' count=1, got %d", counts["deep learning"])
	}
	if counts["computer vision"] != 0 {
		t.Fatalf("expected 'computer vision' count=0, got %d", counts["computer vision"])
	}
}

func TestPhraseMatchedWhenPresent(t *testing.T) {
	resume := "I have experience with machine learning pipelines"
	job := "proficiency in machine learning is required"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, term := range got.OverlapTerms {
		if term == "machine learning" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected 'machine learning' in OverlapTerms, got %v", got.OverlapTerms)
	}
}

func TestPhraseMissingWhenAbsent(t *testing.T) {
	resume := "python developer with statistics background"
	job := "experience with machine learning required"

	got, err := Compute(resume, job)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, ts := range got.MissingJobTerms {
		if ts.Term == "machine learning" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected 'machine learning' in MissingJobTerms, got %v", got.MissingJobTerms)
	}
}

func TestPhraseIDFExceedsConstituentUnigrams(t *testing.T) {
	phraseVal := lookupIDF("machine learning")
	machineVal := lookupIDF("machine")
	learningVal := lookupIDF("learning")

	if phraseVal <= machineVal {
		t.Fatalf("phrase IDF %.2f should exceed 'machine' IDF %.2f", phraseVal, machineVal)
	}
	if phraseVal <= learningVal {
		t.Fatalf("phrase IDF %.2f should exceed 'learning' IDF %.2f", phraseVal, learningVal)
	}
}

func TestCuratedPhraseNotSuppressedInOtherCategory(t *testing.T) {
	profile := profiles.Get(profiles.DisciplineITSoftware)
	resume := "python developer"
	job := "python developer with finite element analysis experience"

	got, err := ComputeWithProfile(resume, job, profile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, ts := range got.MissingJobTerms {
		if ts.Term == "finite element" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected curated phrase 'finite element' visible in MissingJobTerms, got %v", got.MissingJobTerms)
	}
}

func TestPhraseBucketAssignment(t *testing.T) {
	bucket := bucketForTerm(profiles.Get(profiles.DisciplineITSoftware), "cloud native")
	if bucket != categoryCloudDevOps {
		t.Fatalf("expected 'cloud native' in %s, got %s", categoryCloudDevOps, bucket)
	}
}
