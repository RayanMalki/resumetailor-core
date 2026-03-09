package classifier

import (
	"sort"
	"strings"
	"unicode"

	"github.com/RayanMalki/resumetailor-core/nlp"
	"github.com/RayanMalki/resumetailor-core/profiles"
)

const defaultLowConfidenceThreshold = 0.18

type EvidenceTerm struct {
	Term  string  `json:"term"`
	Score float64 `json:"score"`
}

type Detection struct {
	Discipline    profiles.Discipline `json:"discipline"`
	Confidence    float64             `json:"confidence"`
	LowConfidence bool                `json:"lowConfidence"`
	Evidence      []EvidenceTerm      `json:"evidence"`
	Scores        map[string]float64  `json:"scores,omitempty"`
	Source        profiles.Source     `json:"source"`
}

func Detect(resumeText, jobText string) Detection {
	return detectWithThreshold(resumeText, jobText, defaultLowConfidenceThreshold)
}

func Resolve(resumeText, jobText string, override *profiles.Discipline) Detection {
	base := Detect(resumeText, jobText)
	if override == nil {
		base.Source = profiles.SourceAuto
		return base
	}
	base.Discipline = *override
	base.Confidence = 1.0
	base.LowConfidence = false
	base.Source = profiles.SourceUserOverride
	return base
}

func detectWithThreshold(resumeText, jobText string, threshold float64) Detection {
	text := normalizeText(resumeText + "\n" + jobText)
	jobLead := normalizeText(prefix(jobText, 900))
	resumeLead := normalizeText(prefix(resumeText, 650))
	tokenFreq := freq(tokenize(text))

	allProfiles := profiles.GetAll()
	if len(allProfiles) == 0 {
		return Detection{
			Discipline:    profiles.DefaultDiscipline(),
			Confidence:    0,
			LowConfidence: true,
			Evidence:      nil,
			Scores:        nil,
			Source:        profiles.SourceAuto,
		}
	}

	type scored struct {
		discipline profiles.Discipline
		score      float64
		evidence   []EvidenceTerm
	}
	results := make([]scored, 0, len(allProfiles))
	scores := make(map[string]float64, len(allProfiles))
	for discipline, profile := range allProfiles {
		total, evidence := scoreProfile(profile, tokenFreq, jobLead, resumeLead)
		results = append(results, scored{
			discipline: discipline,
			score:      total,
			evidence:   evidence,
		})
		scores[string(discipline)] = total
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].score == results[j].score {
			return results[i].discipline < results[j].discipline
		}
		return results[i].score > results[j].score
	})

	top := results[0]
	second := scored{score: 0}
	if len(results) > 1 {
		second = results[1]
	}

	confidence := 0.0
	if top.score > 0 {
		confidence = (top.score - second.score) / top.score
		if confidence < 0 {
			confidence = 0
		}
		if confidence > 1 {
			confidence = 1
		}
	}
	lowConfidence := top.score == 0 || confidence < threshold

	return Detection{
		Discipline:    top.discipline,
		Confidence:    confidence,
		LowConfidence: lowConfidence,
		Evidence:      top.evidence,
		Scores:        scores,
		Source:        profiles.SourceAuto,
	}
}

func scoreProfile(profile profiles.Profile, tokenFreq map[string]int, jobLead, resumeLead string) (float64, []EvidenceTerm) {
	total := 0.0
	evidence := make([]EvidenceTerm, 0, 12)

	for term, meta := range profile.CanonicalTerms {
		count := tokenFreq[term]
		if count == 0 && strings.Contains(jobLead, term) {
			count = 1
		}
		if count == 0 && strings.Contains(resumeLead, term) {
			count = 1
		}
		if count == 0 {
			continue
		}

		weight := meta.Weight
		if weight <= 0 {
			weight = 1.0
		}
		// Signals near role headers are stronger than body text.
		if strings.Contains(jobLead, term) {
			weight *= 1.35
		}
		if strings.Contains(resumeLead, term) {
			weight *= 1.15
		}
		score := float64(count) * weight
		total += score
		evidence = append(evidence, EvidenceTerm{Term: term, Score: score})
	}

	sort.Slice(evidence, func(i, j int) bool {
		if evidence[i].Score == evidence[j].Score {
			return evidence[i].Term < evidence[j].Term
		}
		return evidence[i].Score > evidence[j].Score
	})
	if len(evidence) > 8 {
		evidence = evidence[:8]
	}

	return total, evidence
}

func prefix(s string, n int) string {
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[:n]
}

func freq(tokens []string) map[string]int {
	out := make(map[string]int, len(tokens))
	for _, token := range tokens {
		out[token]++
	}
	return out
}

func tokenize(text string) []string {
	if text == "" {
		return nil
	}
	var out []string
	var b strings.Builder
	flush := func() {
		if b.Len() < 2 {
			b.Reset()
			return
		}
		token := b.String()
		out = append(out, token)
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
	return out
}

func normalizeText(text string) string {
	text = strings.ToLower(text)
	text = nlp.StripAccents(text)
	replacer := strings.NewReplacer(
		"ci/cd", " cicd ",
		"ci-cd", " cicd ",
		"gd&t", " gdt ",
		"g.d.t", " gdt ",
		"supply chain", " supplychain ",
		"chaine d approvisionnement", " supplychain ",
	)
	return replacer.Replace(text)
}
