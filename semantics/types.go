package semantics

// Field represents a single field in an H5P semantics definition
type Field struct {
	Name        string `json:"name,omitempty"`
	Type        string `json:"type"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	Importance  string `json:"importance,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
	Default     any    `json:"default,omitempty"`
	Common      bool   `json:"common,omitempty"`
	Widget      string `json:"widget,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`

	// For group types
	Fields   []Field `json:"fields,omitempty"`
	Expanded bool    `json:"expanded,omitempty"`

	// For list types
	Entity     string `json:"entity,omitempty"`
	Min        int    `json:"min,omitempty"`
	Max        int    `json:"max,omitempty"`
	DefaultNum int    `json:"defaultNum,omitempty"`
	Field      *Field `json:"field,omitempty"`

	// For select and library types - polymorphic options field
	Options interface{} `json:"options,omitempty"`

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

// GetLibraryOptions returns the options as a slice of library strings.
// This is used for library type fields where options contains library names/versions.
// Returns nil if the options are not in the expected format.
func (f *Field) GetLibraryOptions() []string {
	if f.Options == nil {
		return nil
	}

	// Handle the case where JSON unmarshaling creates []interface{}
	if optionsSlice, ok := f.Options.([]interface{}); ok {
		result := make([]string, 0, len(optionsSlice))
		for _, opt := range optionsSlice {
			if str, ok := opt.(string); ok {
				result = append(result, str)
			} else {
				return nil // Mixed types, not a library options array
			}
		}
		return result
	}

	// Handle direct []string assignment
	if optionsSlice, ok := f.Options.([]string); ok {
		return optionsSlice
	}

	return nil
}

// GetSelectOptions returns the options as a slice of SelectOption structs.
// This is used for select type fields where options contains value/label pairs.
// Returns nil if the options are not in the expected format.
func (f *Field) GetSelectOptions() []SelectOption {
	if f.Options == nil {
		return nil
	}

	// Handle the case where JSON unmarshaling creates []interface{}
	if optionsSlice, ok := f.Options.([]interface{}); ok {
		result := make([]SelectOption, 0, len(optionsSlice))
		for _, opt := range optionsSlice {
			if optMap, ok := opt.(map[string]interface{}); ok {
				selectOpt := SelectOption{}
				if value, exists := optMap["value"]; exists {
					if str, ok := value.(string); ok {
						selectOpt.Value = str
					} else {
						return nil // Invalid value type
					}
				}
				if label, exists := optMap["label"]; exists {
					if str, ok := label.(string); ok {
						selectOpt.Label = str
					} else {
						return nil // Invalid label type
					}
				}
				result = append(result, selectOpt)
			} else {
				return nil // Not an object, likely library options
			}
		}
		return result
	}

	// Handle direct []SelectOption assignment
	if optionsSlice, ok := f.Options.([]SelectOption); ok {
		return optionsSlice
	}

	return nil
}

// SetLibraryOptions sets the options field with library strings.
func (f *Field) SetLibraryOptions(options []string) {
	f.Options = options
}

// SetSelectOptions sets the options field with SelectOption structs.
func (f *Field) SetSelectOptions(options []SelectOption) {
	f.Options = options
}
