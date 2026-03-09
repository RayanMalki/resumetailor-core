package profiles

type Discipline string

const (
	DisciplineMechanical          Discipline = "mechanical"
	DisciplineElectrical          Discipline = "electrical"
	DisciplineIndustrialLogistics Discipline = "industrial_logistics"
	DisciplineAerospace           Discipline = "aerospace"
	DisciplineITSoftware          Discipline = "it_software"
)

type Source string

const (
	SourceAuto         Source = "auto"
	SourceUserOverride Source = "user_override"
)

type TermMeta struct {
	Bucket string
	Weight float64
}

type PromptHints struct {
	RoleFocus     []string
	EvidenceFocus []string
	ActionVerbs   []string
}

type Profile struct {
	Discipline        Discipline
	Version           string
	CanonicalTerms    map[string]TermMeta
	Synonyms          map[string]string
	StopwordOverrides map[string]bool
	LowSignalTerms    map[string]bool
	Buckets           map[string][]string
	BucketWeights     map[string]float64
	PromptHints       PromptHints
}

func ValidDisciplines() []Discipline {
	return []Discipline{
		DisciplineMechanical,
		DisciplineElectrical,
		DisciplineIndustrialLogistics,
		DisciplineAerospace,
		DisciplineITSoftware,
	}
}

func ParseDiscipline(raw string) (Discipline, bool) {
	switch Discipline(raw) {
	case DisciplineMechanical,
		DisciplineElectrical,
		DisciplineIndustrialLogistics,
		DisciplineAerospace,
		DisciplineITSoftware:
		return Discipline(raw), true
	default:
		return "", false
	}
}
