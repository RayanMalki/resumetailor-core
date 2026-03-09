package scoring

import (
	"strings"
	"unicode"

	"github.com/RayanMalki/resumetailor-core/corpus"
	"github.com/RayanMalki/resumetailor-core/profiles"
)

const highIDFOtherThreshold = 6.2

var genericLowSignalTerms = map[string]struct{}{
	"experience": {}, "years": {}, "year": {}, "preferred": {}, "required": {}, "requirements": {},
	"ability": {}, "strong": {}, "excellent": {}, "good": {}, "role": {}, "position": {},
	"responsibility": {}, "responsibilities": {}, "candidate": {}, "skills": {}, "skill": {},
	"knowledge": {}, "understanding": {}, "familiarity": {}, "using": {}, "plus": {}, "must": {},
	"nice": {}, "seeking": {}, "join": {}, "company": {}, "business": {}, "customer": {},
	"customers": {}, "team": {}, "teams": {}, "support": {}, "supporting": {}, "work": {},
	"working": {}, "develop": {}, "development": {}, "design": {}, "implement": {},
	"implementation": {}, "maintain": {}, "maintenance": {}, "build": {}, "building": {},
	"solutions": {}, "solution": {}, "environments": {}, "environment": {},
}

func isLowSignalTerm(term, category string, idf float64, profile profiles.Profile) bool {
	if profiles.IsLowSignal(profile, term) {
		return true
	}
	if isDigitsOnly(term) {
		return true
	}
	if _, ok := genericLowSignalTerms[term]; ok {
		return true
	}
	// Curated phrases are always surfaced regardless of category or IDF.
	if _, ok := corpus.PhraseIDF[term]; ok {
		return false
	}
	// "Other" terms are shown only if they look highly distinctive.
	if category == "other" && idf < highIDFOtherThreshold {
		if len(term) <= 3 {
			return true
		}
		if strings.HasSuffix(term, "ing") || strings.HasSuffix(term, "tion") || strings.HasSuffix(term, "ment") {
			return true
		}
		if _, ok := genericLowSignalTerms[term]; ok {
			return true
		}
	}
	// Non-standard tokens are often parser noise.
	if strings.Contains(term, "_") {
		return true
	}
	return false
}

func isDigitsOnly(term string) bool {
	if term == "" {
		return false
	}
	for _, r := range term {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
