package h5p

import (
	"github.com/grokify/h5p-go/schemas"
)

type QuestionSet struct {
	ProgressType       string           `json:"progressType,omitempty"`
	PassPercentage     int              `json:"passPercentage,omitempty"`
	BackgroundImage    *BackgroundImage `json:"backgroundImage,omitempty"`
	Questions          []Question       `json:"questions"`
	ShowIntroPage      bool             `json:"showIntroPage,omitempty"`
	StartButtonText    string           `json:"startButtonText,omitempty"`
	Introduction       string           `json:"introduction,omitempty"`
	Title              string           `json:"title,omitempty"`
	ShowResultPage     bool             `json:"showResultPage,omitempty"`
	Message            string           `json:"message,omitempty"`
	SolutionButtonText string           `json:"solutionButtonText,omitempty"`
	OverallFeedback    []FeedbackRange  `json:"overallFeedback,omitempty"`
}

type BackgroundImage struct {
	Path      string     `json:"path"`
	Mime      string     `json:"mime"`
	Copyright *Copyright `json:"copyright,omitempty"`
}

type Copyright struct {
	Title   string `json:"title,omitempty"`
	Author  string `json:"author,omitempty"`
	License string `json:"license,omitempty"`
	Version string `json:"version,omitempty"`
	Source  string `json:"source,omitempty"`
}

type Question struct {
	Library string      `json:"library"`
	Params  interface{} `json:"params"`
}

// MultiChoiceQuestion represents a typed H5P MultiChoice question
type MultiChoiceQuestion struct {
	Library string                     `json:"library"`
	Params  *schemas.MultiChoiceParams `json:"params"`
}

// ToQuestion converts a MultiChoiceQuestion to a generic Question
func (mcq *MultiChoiceQuestion) ToQuestion() *Question {
	return &Question{
		Library: mcq.Library,
		Params:  mcq.Params,
	}
}

// NewMultiChoiceQuestion creates a new typed MultiChoice question
func NewMultiChoiceQuestion(params *schemas.MultiChoiceParams) *MultiChoiceQuestion {
	return &MultiChoiceQuestion{
		Library: "H5P.MultiChoice 1.16",
		Params:  params,
	}
}

type FeedbackRange struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Text string `json:"text"`
}

// Legacy types - use schemas.MultiChoiceParams instead
// Kept for backward compatibility, will be deprecated
type Answer struct {
	Text     string `json:"text"`
	Correct  bool   `json:"correct"`
	Feedback string `json:"feedback,omitempty"`
}
