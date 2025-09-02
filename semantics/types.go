package semantics

// Field represents a single field in an H5P semantics definition
type Field struct {
	Name         string      `json:"name,omitempty"`
	Type         string      `json:"type"`
	Label        string      `json:"label,omitempty"`
	Description  string      `json:"description,omitempty"`
	Importance   string      `json:"importance,omitempty"`
	Optional     bool        `json:"optional,omitempty"`
	Default      interface{} `json:"default,omitempty"`
	Common       bool        `json:"common,omitempty"`
	Widget       string      `json:"widget,omitempty"`
	Placeholder  string      `json:"placeholder,omitempty"`
	
	// For group types
	Fields   []Field `json:"fields,omitempty"`
	Expanded bool    `json:"expanded,omitempty"`
	
	// For list types
	Entity     string `json:"entity,omitempty"`
	Min        int    `json:"min,omitempty"`
	Max        int    `json:"max,omitempty"`
	DefaultNum int    `json:"defaultNum,omitempty"`
	Field      *Field `json:"field,omitempty"`
	
	// For select types
	Options []SelectOption `json:"options,omitempty"`
	
	// For library types (options as string array)
	LibraryOptions []string `json:"options,omitempty"`
	
	// For number types
	MinValue int    `json:"minValue,omitempty"`
	MaxValue int    `json:"maxValue,omitempty"`
	Step     int    `json:"step,omitempty"`
	Unit     string `json:"unit,omitempty"`
	
	// For text types
	MaxLength int      `json:"maxLength,omitempty"`
	Tags      []string `json:"tags,omitempty"`
	
	// Widget-specific properties
	ShowWhen *ShowWhen `json:"showWhen,omitempty"`
}

// SelectOption represents an option in a select field
type SelectOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ShowWhen defines conditional visibility rules for fields
type ShowWhen struct {
	Rules []ShowRule `json:"rules"`
}

// ShowRule defines a single visibility rule
type ShowRule struct {
	Field  string      `json:"field"`
	Equals interface{} `json:"equals"`
}

// SemanticDefinition represents the top-level semantics array
type SemanticDefinition []Field