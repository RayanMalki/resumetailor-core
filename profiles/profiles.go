package profiles

import (
	"strings"

	"github.com/RayanMalki/resumetailor-core/profiles/data"
)

var allProfiles = loadProfiles()

func Version() string {
	return data.Version
}

func DefaultDiscipline() Discipline {
	return DisciplineITSoftware
}

func Get(d Discipline) Profile {
	if p, ok := allProfiles[d]; ok {
		return p
	}
	return allProfiles[DefaultDiscipline()]
}

func GetAll() map[Discipline]Profile {
	out := make(map[Discipline]Profile, len(allProfiles))
	for k, v := range allProfiles {
		out[k] = cloneProfile(v)
	}
	return out
}

func BucketForTerm(profile Profile, term string) string {
	if meta, ok := profile.CanonicalTerms[term]; ok && meta.Bucket != "" {
		return meta.Bucket
	}
	for bucket, terms := range profile.Buckets {
		for _, candidate := range terms {
			if candidate == term {
				return bucket
			}
		}
	}
	return "other"
}

func Canonicalize(profile Profile, token string) string {
	token = strings.TrimSpace(token)
	if token == "" {
		return token
	}
	if canon, ok := profile.Synonyms[token]; ok && canon != "" {
		return canon
	}
	return token
}

func IsLowSignal(profile Profile, token string) bool {
	return profile.LowSignalTerms[token]
}

func BucketWeight(profile Profile, bucket string) float64 {
	if w, ok := profile.BucketWeights[bucket]; ok && w > 0 {
		return w
	}
	return 1.0
}

func loadProfiles() map[Discipline]Profile {
	raw := data.Profiles()
	out := make(map[Discipline]Profile, len(raw))
	for key, rp := range raw {
		discipline, ok := ParseDiscipline(key)
		if !ok {
			continue
		}
		p := Profile{
			Discipline:        discipline,
			Version:           data.Version,
			CanonicalTerms:    make(map[string]TermMeta, len(rp.CanonicalTerms)),
			Synonyms:          cloneStringMap(rp.Synonyms),
			StopwordOverrides: cloneBoolMap(rp.StopwordOverrides),
			LowSignalTerms:    cloneBoolMap(rp.LowSignalTerms),
			Buckets:           cloneBuckets(rp.Buckets),
			BucketWeights:     cloneFloatMap(rp.BucketWeights),
			PromptHints: PromptHints{
				RoleFocus:     cloneStringSlice(rp.PromptHints.RoleFocus),
				EvidenceFocus: cloneStringSlice(rp.PromptHints.EvidenceFocus),
				ActionVerbs:   cloneStringSlice(rp.PromptHints.ActionVerbs),
			},
		}
		for term, meta := range rp.CanonicalTerms {
			p.CanonicalTerms[term] = TermMeta{
				Bucket: meta.Bucket,
				Weight: meta.Weight,
			}
		}
		out[discipline] = p
	}
	return out
}

func cloneProfile(in Profile) Profile {
	out := Profile{
		Discipline:        in.Discipline,
		Version:           in.Version,
		CanonicalTerms:    make(map[string]TermMeta, len(in.CanonicalTerms)),
		Synonyms:          cloneStringMap(in.Synonyms),
		StopwordOverrides: cloneBoolMap(in.StopwordOverrides),
		LowSignalTerms:    cloneBoolMap(in.LowSignalTerms),
		Buckets:           cloneBuckets(in.Buckets),
		BucketWeights:     cloneFloatMap(in.BucketWeights),
		PromptHints: PromptHints{
			RoleFocus:     cloneStringSlice(in.PromptHints.RoleFocus),
			EvidenceFocus: cloneStringSlice(in.PromptHints.EvidenceFocus),
			ActionVerbs:   cloneStringSlice(in.PromptHints.ActionVerbs),
		},
	}
	for k, v := range in.CanonicalTerms {
		out.CanonicalTerms[k] = v
	}
	return out
}

func cloneStringMap(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneBoolMap(in map[string]bool) map[string]bool {
	out := make(map[string]bool, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneFloatMap(in map[string]float64) map[string]float64 {
	out := make(map[string]float64, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneBuckets(in map[string][]string) map[string][]string {
	out := make(map[string][]string, len(in))
	for k, v := range in {
		out[k] = cloneStringSlice(v)
	}
	return out
}

func cloneStringSlice(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}
