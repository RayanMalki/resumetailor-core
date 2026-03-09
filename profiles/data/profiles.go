package data

type RawTermMeta struct {
	Bucket string
	Weight float64
}

type RawPromptHints struct {
	RoleFocus     []string
	EvidenceFocus []string
	ActionVerbs   []string
}

type RawProfile struct {
	CanonicalTerms    map[string]RawTermMeta
	Synonyms          map[string]string
	StopwordOverrides map[string]bool
	LowSignalTerms    map[string]bool
	Buckets           map[string][]string
	BucketWeights     map[string]float64
	PromptHints       RawPromptHints
}

const Version = "profiles-v1-2026-02-14"

func Profiles() map[string]RawProfile {
	return map[string]RawProfile{
		"it_software":          itSoftware(),
		"mechanical":           mechanical(),
		"electrical":           electrical(),
		"industrial_logistics": industrialLogistics(),
		"aerospace":            aerospace(),
	}
}

func commonLowSignal() map[string]bool {
	return map[string]bool{
		"experience": true, "required": true, "requirements": true, "preferred": true,
		"ability": true, "strong": true, "excellent": true, "role": true, "position": true,
		"team": true, "teams": true, "work": true, "working": true, "support": true,
		"development": true, "solutions": true, "process": true, "processes": true,
		"responsibility": true, "responsibilities": true, "candidate": true, "knowledge": true,
		"years": true, "year": true, "plus": true, "must": true,
	}
}

func itSoftware() RawProfile {
	return RawProfile{
		CanonicalTerms: map[string]RawTermMeta{
			"python": {Bucket: "languages", Weight: 1.1}, "java": {Bucket: "languages", Weight: 1.1},
			"javascript": {Bucket: "languages", Weight: 1.05}, "typescript": {Bucket: "languages", Weight: 1.1},
			"golang": {Bucket: "languages", Weight: 1.2}, "csharp": {Bucket: "languages", Weight: 1.1},
			"docker": {Bucket: "cloud_devops_db", Weight: 1.15}, "kubernetes": {Bucket: "cloud_devops_db", Weight: 1.2},
			"aws": {Bucket: "cloud_devops_db", Weight: 1.1}, "azure": {Bucket: "cloud_devops_db", Weight: 1.1},
			"gcp": {Bucket: "cloud_devops_db", Weight: 1.1}, "postgresql": {Bucket: "cloud_devops_db", Weight: 1.0},
			"mysql": {Bucket: "cloud_devops_db", Weight: 1.0}, "mongodb": {Bucket: "cloud_devops_db", Weight: 1.0},
			"ci": {Bucket: "practices", Weight: 1.05}, "cd": {Bucket: "practices", Weight: 1.05},
			"cicd": {Bucket: "practices", Weight: 1.1}, "agile": {Bucket: "practices", Weight: 1.0},
			"scrum": {Bucket: "practices", Weight: 1.0}, "microservice": {Bucket: "practices", Weight: 1.05},
			"api": {Bucket: "practices", Weight: 1.0}, "rest": {Bucket: "practices", Weight: 1.0},
			"communication": {Bucket: "soft_skills", Weight: 0.9}, "leadership": {Bucket: "soft_skills", Weight: 0.9},
		},
		Synonyms: map[string]string{
			"js": "javascript", "ts": "typescript", "go": "golang", "postgres": "postgresql",
			"k8s": "kubernetes", "node": "nodejs",
		},
		StopwordOverrides: map[string]bool{},
		LowSignalTerms:    commonLowSignal(),
		Buckets: map[string][]string{
			"languages":       {"python", "java", "javascript", "typescript", "golang", "csharp"},
			"cloud_devops_db": {"docker", "kubernetes", "aws", "azure", "gcp", "postgresql", "mysql", "mongodb"},
			"practices":       {"ci", "cd", "cicd", "agile", "scrum", "microservice", "api", "rest"},
			"soft_skills":     {"communication", "leadership"},
			"other":           {},
		},
		BucketWeights: map[string]float64{
			"languages": 1.05, "cloud_devops_db": 1.1, "practices": 1.0, "soft_skills": 0.85, "other": 0.8,
		},
		PromptHints: RawPromptHints{
			RoleFocus:     []string{"backend systems", "full-stack delivery", "software quality", "production reliability"},
			EvidenceFocus: []string{"technical stack fit", "delivery impact", "scalable systems", "automation"},
			ActionVerbs:   []string{"built", "implemented", "optimized", "deployed"},
		},
	}
}

func mechanical() RawProfile {
	low := commonLowSignal()
	low["component"] = true
	low["components"] = true
	return RawProfile{
		CanonicalTerms: map[string]RawTermMeta{
			"solidworks": {Bucket: "design_tools", Weight: 1.25}, "catia": {Bucket: "design_tools", Weight: 1.25},
			"nx": {Bucket: "design_tools", Weight: 1.1}, "autocad": {Bucket: "design_tools", Weight: 1.15},
			"ansys": {Bucket: "simulation_analysis", Weight: 1.25}, "fea": {Bucket: "simulation_analysis", Weight: 1.3},
			"cfd": {Bucket: "simulation_analysis", Weight: 1.25}, "thermodynamics": {Bucket: "core_engineering", Weight: 1.2},
			"mechanics": {Bucket: "core_engineering", Weight: 1.1}, "materials": {Bucket: "core_engineering", Weight: 1.1},
			"manufacturing": {Bucket: "manufacturing_quality", Weight: 1.2}, "machining": {Bucket: "manufacturing_quality", Weight: 1.2},
			"cnc": {Bucket: "manufacturing_quality", Weight: 1.2}, "gdt": {Bucket: "manufacturing_quality", Weight: 1.2},
			"tolerance": {Bucket: "manufacturing_quality", Weight: 1.15}, "fmea": {Bucket: "manufacturing_quality", Weight: 1.2},
			"sixsigma": {Bucket: "manufacturing_quality", Weight: 1.15}, "plm": {Bucket: "design_tools", Weight: 1.1},
			"pneumatics": {Bucket: "core_engineering", Weight: 1.15}, "hydraulics": {Bucket: "core_engineering", Weight: 1.15},
			"communication": {Bucket: "soft_skills", Weight: 0.85}, "leadership": {Bucket: "soft_skills", Weight: 0.85},
		},
		Synonyms: map[string]string{
			"mecanique": "mechanics", "tolerances": "tolerance", "gd&t": "gdt",
		},
		StopwordOverrides: map[string]bool{},
		LowSignalTerms:    low,
		Buckets: map[string][]string{
			"design_tools":          {"solidworks", "catia", "nx", "autocad", "plm"},
			"simulation_analysis":   {"ansys", "fea", "cfd"},
			"core_engineering":      {"thermodynamics", "mechanics", "materials", "pneumatics", "hydraulics"},
			"manufacturing_quality": {"manufacturing", "machining", "cnc", "gdt", "tolerance", "fmea", "sixsigma"},
			"soft_skills":           {"communication", "leadership"},
			"other":                 {},
		},
		BucketWeights: map[string]float64{
			"design_tools": 1.1, "simulation_analysis": 1.2, "core_engineering": 1.15,
			"manufacturing_quality": 1.2, "soft_skills": 0.85, "other": 0.8,
		},
		PromptHints: RawPromptHints{
			RoleFocus:     []string{"mechanical design", "analysis and validation", "manufacturing readiness"},
			EvidenceFocus: []string{"simulation results", "tolerance/quality ownership", "lab and prototype delivery"},
			ActionVerbs:   []string{"designed", "validated", "modeled", "optimized"},
		},
	}
}

func electrical() RawProfile {
	return RawProfile{
		CanonicalTerms: map[string]RawTermMeta{
			"circuit": {Bucket: "core_electrical", Weight: 1.2}, "pcb": {Bucket: "hardware_tools", Weight: 1.25},
			"altium": {Bucket: "hardware_tools", Weight: 1.25}, "kicad": {Bucket: "hardware_tools", Weight: 1.2},
			"fpga": {Bucket: "embedded_controls", Weight: 1.25}, "vhdl": {Bucket: "embedded_controls", Weight: 1.25},
			"verilog": {Bucket: "embedded_controls", Weight: 1.25}, "embedded": {Bucket: "embedded_controls", Weight: 1.2},
			"microcontroller": {Bucket: "embedded_controls", Weight: 1.2}, "plc": {Bucket: "industrial_power", Weight: 1.2},
			"scada": {Bucket: "industrial_power", Weight: 1.2}, "power": {Bucket: "industrial_power", Weight: 1.1},
			"inverter": {Bucket: "industrial_power", Weight: 1.15}, "signal": {Bucket: "core_electrical", Weight: 1.05},
			"electronics": {Bucket: "core_electrical", Weight: 1.15}, "instrumentation": {Bucket: "core_electrical", Weight: 1.1},
			"matlab": {Bucket: "analysis_testing", Weight: 1.1}, "simulink": {Bucket: "analysis_testing", Weight: 1.2},
			"oscilloscope": {Bucket: "analysis_testing", Weight: 1.15}, "communication": {Bucket: "soft_skills", Weight: 0.85},
		},
		Synonyms: map[string]string{
			"electrique": "electronics", "electrical": "electronics", "microcontroleur": "microcontroller",
		},
		StopwordOverrides: map[string]bool{},
		LowSignalTerms:    commonLowSignal(),
		Buckets: map[string][]string{
			"hardware_tools":    {"pcb", "altium", "kicad"},
			"embedded_controls": {"fpga", "vhdl", "verilog", "embedded", "microcontroller"},
			"industrial_power":  {"plc", "scada", "power", "inverter"},
			"core_electrical":   {"circuit", "signal", "electronics", "instrumentation"},
			"analysis_testing":  {"matlab", "simulink", "oscilloscope"},
			"soft_skills":       {"communication"},
			"other":             {},
		},
		BucketWeights: map[string]float64{
			"hardware_tools": 1.1, "embedded_controls": 1.2, "industrial_power": 1.15,
			"core_electrical": 1.1, "analysis_testing": 1.05, "soft_skills": 0.85, "other": 0.8,
		},
		PromptHints: RawPromptHints{
			RoleFocus:     []string{"electronic design", "embedded control", "validation and test"},
			EvidenceFocus: []string{"hardware bring-up", "signal/power integrity", "lab instrumentation"},
			ActionVerbs:   []string{"designed", "debugged", "validated", "integrated"},
		},
	}
}

func industrialLogistics() RawProfile {
	return RawProfile{
		CanonicalTerms: map[string]RawTermMeta{
			"lean": {Bucket: "process_excellence", Weight: 1.2}, "sixsigma": {Bucket: "process_excellence", Weight: 1.2},
			"kaizen": {Bucket: "process_excellence", Weight: 1.2}, "oee": {Bucket: "process_excellence", Weight: 1.15},
			"supplychain": {Bucket: "operations_logistics", Weight: 1.25}, "logistics": {Bucket: "operations_logistics", Weight: 1.2},
			"inventory": {Bucket: "operations_logistics", Weight: 1.1}, "warehouse": {Bucket: "operations_logistics", Weight: 1.1},
			"procurement": {Bucket: "operations_logistics", Weight: 1.15}, "forecasting": {Bucket: "planning_analytics", Weight: 1.1},
			"demand": {Bucket: "planning_analytics", Weight: 1.0}, "planning": {Bucket: "planning_analytics", Weight: 1.0},
			"optimization": {Bucket: "planning_analytics", Weight: 1.15}, "operations": {Bucket: "operations_logistics", Weight: 1.05},
			"throughput": {Bucket: "process_excellence", Weight: 1.15}, "erp": {Bucket: "systems_tools", Weight: 1.1},
			"sap": {Bucket: "systems_tools", Weight: 1.15}, "wms": {Bucket: "systems_tools", Weight: 1.15},
			"transport": {Bucket: "operations_logistics", Weight: 1.1}, "communication": {Bucket: "soft_skills", Weight: 0.85},
		},
		Synonyms: map[string]string{
			"gestion": "operations", "logistique": "logistics", "approvisionnement": "supplychain",
			"chaine": "supplychain", "entrepot": "warehouse",
		},
		StopwordOverrides: map[string]bool{},
		LowSignalTerms:    commonLowSignal(),
		Buckets: map[string][]string{
			"process_excellence":   {"lean", "sixsigma", "kaizen", "oee", "throughput"},
			"operations_logistics": {"supplychain", "logistics", "inventory", "warehouse", "procurement", "operations", "transport"},
			"planning_analytics":   {"forecasting", "demand", "planning", "optimization"},
			"systems_tools":        {"erp", "sap", "wms"},
			"soft_skills":          {"communication"},
			"other":                {},
		},
		BucketWeights: map[string]float64{
			"process_excellence": 1.2, "operations_logistics": 1.2, "planning_analytics": 1.1,
			"systems_tools": 1.05, "soft_skills": 0.85, "other": 0.8,
		},
		PromptHints: RawPromptHints{
			RoleFocus:     []string{"operations performance", "logistics flow", "continuous improvement"},
			EvidenceFocus: []string{"cycle time reduction", "cost/throughput impact", "planning accuracy"},
			ActionVerbs:   []string{"optimized", "streamlined", "coordinated", "improved"},
		},
	}
}

func aerospace() RawProfile {
	return RawProfile{
		CanonicalTerms: map[string]RawTermMeta{
			"aerodynamics": {Bucket: "core_aero", Weight: 1.25}, "propulsion": {Bucket: "core_aero", Weight: 1.2},
			"avionics": {Bucket: "systems_avionics", Weight: 1.25}, "flight": {Bucket: "systems_avionics", Weight: 1.15},
			"aircraft": {Bucket: "core_aero", Weight: 1.2}, "spacecraft": {Bucket: "core_aero", Weight: 1.25},
			"satellite": {Bucket: "systems_avionics", Weight: 1.2}, "orbital": {Bucket: "core_aero", Weight: 1.2},
			"cfd": {Bucket: "analysis_validation", Weight: 1.2}, "fea": {Bucket: "analysis_validation", Weight: 1.2},
			"ansys": {Bucket: "analysis_validation", Weight: 1.2}, "matlab": {Bucket: "analysis_validation", Weight: 1.1},
			"simulink": {Bucket: "analysis_validation", Weight: 1.15}, "arinc": {Bucket: "standards_safety", Weight: 1.25},
			"do178": {Bucket: "standards_safety", Weight: 1.3}, "do254": {Bucket: "standards_safety", Weight: 1.3},
			"verification": {Bucket: "analysis_validation", Weight: 1.1}, "validation": {Bucket: "analysis_validation", Weight: 1.1},
			"safety": {Bucket: "standards_safety", Weight: 1.2}, "communication": {Bucket: "soft_skills", Weight: 0.85},
		},
		Synonyms: map[string]string{
			"aerospatial": "aerospace", "aeronautique": "aerodynamics", "aerospatiale": "aerospace",
			"satcom": "satellite",
		},
		StopwordOverrides: map[string]bool{},
		LowSignalTerms:    commonLowSignal(),
		Buckets: map[string][]string{
			"core_aero":           {"aerodynamics", "propulsion", "aircraft", "spacecraft", "orbital"},
			"systems_avionics":    {"avionics", "flight", "satellite"},
			"analysis_validation": {"cfd", "fea", "ansys", "matlab", "simulink", "verification", "validation"},
			"standards_safety":    {"arinc", "do178", "do254", "safety"},
			"soft_skills":         {"communication"},
			"other":               {},
		},
		BucketWeights: map[string]float64{
			"core_aero": 1.15, "systems_avionics": 1.2, "analysis_validation": 1.15,
			"standards_safety": 1.25, "soft_skills": 0.85, "other": 0.8,
		},
		PromptHints: RawPromptHints{
			RoleFocus:     []string{"aircraft/space systems", "verification and validation", "safety and standards compliance"},
			EvidenceFocus: []string{"test evidence", "modeling fidelity", "certification-aligned deliverables"},
			ActionVerbs:   []string{"validated", "verified", "analyzed", "designed"},
		},
	}
}
