package nlp

import (
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

var techTokenReplacer = strings.NewReplacer(
	"asp.net", " aspnet ",
	"next.js", " nextjs ",
	"node.js", " nodejs ",
	"nuxt.js", " nuxtjs ",
	"c++", " cpp ",
	"c#", " csharp ",
	"f#", " fsharp ",
	".net", " dotnet ",
	"ci/cd", " cicd ",
	"ci-cd", " cicd ",
)

// NormalizeForTokenization lowercases, strips accents, and normalizes
// punctuation-heavy tech forms (C#, C++, .NET, etc.) before tokenizing.
func NormalizeForTokenization(text string) string {
	text = strings.ToLower(StripAccents(text))
	return techTokenReplacer.Replace(text)
}

// StripAccents removes diacritics/accents from text using Unicode NFD
// decomposition (e.g. "développement" → "developpement", "résumé" → "resume").
// This allows French and other accented text to match English IDF table terms.
// It also strips modifier letter accents (ˊ U+02CA, ˋ U+02CB, etc.) which
// some PDF extractors produce instead of proper combining accents.
func StripAccents(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range norm.NFD.String(s) {
		// After NFD decomposition, accents become separate combining marks.
		if unicode.Is(unicode.Mn, r) { // Mn = Mark, Nonspacing (combining accents)
			continue
		}
		// Also strip modifier letter accents (U+02B0–U+02FF) which some
		// PDF extractors emit instead of proper combining marks.
		// Includes ˊ (U+02CA), ˋ (U+02CB), ˆ (U+02C6), etc.
		if r >= 0x02B0 && r <= 0x02FF {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
