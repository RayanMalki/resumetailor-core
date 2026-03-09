package classifier

import (
	"testing"

	"github.com/RayanMalki/resumetailor-core/profiles"
)

func TestDetectPerDiscipline(t *testing.T) {
	tests := []struct {
		name string
		job  string
		want profiles.Discipline
	}{
		{
			name: "mechanical",
			job:  "Need SolidWorks, FEA, CFD, tolerance stack-up, GD&T, and manufacturing support for CNC prototypes.",
			want: profiles.DisciplineMechanical,
		},
		{
			name: "electrical",
			job:  "Design PCB in Altium, debug circuits with oscilloscope, and implement embedded firmware on microcontroller.",
			want: profiles.DisciplineElectrical,
		},
		{
			name: "industrial logistics",
			job:  "Drive supply chain optimization, demand planning, SAP ERP coordination, warehouse flow, and lean six sigma initiatives.",
			want: profiles.DisciplineIndustrialLogistics,
		},
		{
			name: "aerospace",
			job:  "Aerospace role focused on avionics integration, flight systems verification, ARINC and DO-178 compliance.",
			want: profiles.DisciplineAerospace,
		},
		{
			name: "it software",
			job:  "Backend engineer with Golang, Kubernetes, AWS, CI/CD and PostgreSQL required.",
			want: profiles.DisciplineITSoftware,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Detect("", tt.job)
			if got.Discipline != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got.Discipline)
			}
		})
	}
}

func TestResolveOverrideWins(t *testing.T) {
	override := profiles.DisciplineElectrical
	got := Resolve("", "Need Kubernetes and Go backend", &override)
	if got.Discipline != override {
		t.Fatalf("expected override discipline %s, got %s", override, got.Discipline)
	}
	if got.Source != profiles.SourceUserOverride {
		t.Fatalf("expected source user_override, got %s", got.Source)
	}
	if got.Confidence != 1 {
		t.Fatalf("expected confidence 1, got %.2f", got.Confidence)
	}
}

func TestLowConfidenceFlag(t *testing.T) {
	got := Detect("student", "team player motivated fast learner")
	if !got.LowConfidence {
		t.Fatalf("expected low confidence for weak signal input, got %+v", got)
	}
}
