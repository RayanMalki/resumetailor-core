package nlp

import (
	"strings"
	"unicode"

	"github.com/RayanMalki/resumetailor-core/profiles"
)

const minTokenLength = 2

// Tokenize tokenizes text using the default IT/Software discipline profile.
func Tokenize(text string) []string {
	return tokenizeWithProfile(text, profiles.Get(profiles.DefaultDiscipline()))
}

// TokenizeWithProfile tokenizes text using a specific discipline profile.
func TokenizeWithProfile(text string, profile profiles.Profile) []string {
	return tokenizeWithProfile(text, profile)
}

func tokenizeWithProfile(text string, profile profiles.Profile) []string {
	if text == "" {
		return nil
	}

	// Normalize accents and punctuation-heavy tech forms before tokenizing.
	text = NormalizeForTokenization(text)

	var tokens []string
	var b strings.Builder

	flush := func() {
		if b.Len() < minTokenLength {
			b.Reset()
			return
		}
		token := b.String()
		if isStopword(token) && !profile.StopwordOverrides[token] {
			b.Reset()
			return
		}
		// Normalize plural → singular only when token is not explicitly canonical.
		// This avoids breaking tools like "solidworks" / "ansys".
		if _, isCanonical := profile.CanonicalTerms[token]; !isCanonical {
			token = depluralize(token)
		}
		if isStopword(token) && !profile.StopwordOverrides[token] {
			b.Reset()
			return
		}
		tokens = append(tokens, canonicalizeWithProfile(token, profile))
		b.Reset()
	}

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		flush()
	}
	flush()

	return tokens
}

func canonicalizeWithProfile(token string, profile profiles.Profile) string {
	token = profiles.Canonicalize(profile, token)
	token = Canonicalize(token)
	token = profiles.Canonicalize(profile, token)
	return token
}

// TermFreq builds a term frequency map from a token slice.
func TermFreq(tokens []string) map[string]int {
	freq := make(map[string]int, len(tokens))
	for _, t := range tokens {
		freq[t]++
	}
	return freq
}
