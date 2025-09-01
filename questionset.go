package goh5p

import (
	"encoding/json"
)

type QuestionSet struct {
	ProgressType     string           `json:"progressType,omitempty"`
	PassPercentage   int              `json:"passPercentage,omitempty"`
	BackgroundImage  *BackgroundImage `json:"backgroundImage,omitempty"`
	Questions        []Question       `json:"questions"`
	ShowIntroPage    bool             `json:"showIntroPage,omitempty"`
	StartButtonText  string           `json:"startButtonText,omitempty"`
	Introduction     string           `json:"introduction,omitempty"`
	Title            string           `json:"title,omitempty"`
	ShowResultPage   bool             `json:"showResultPage,omitempty"`
	Message          string           `json:"message,omitempty"`
	SolutionButtonText string         `json:"solutionButtonText,omitempty"`
	OverallFeedback  []FeedbackRange  `json:"overallFeedback,omitempty"`
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
	Library string          `json:"library"`
	Params  json.RawMessage `json:"params"`
}

type FeedbackRange struct {
	From int    `json:"from"`
	To   int    `json:"to"`
	Text string `json:"text"`
}

type MultipleChoiceParams struct {
	Question         string           `json:"question"`
	Answers          []Answer         `json:"answers"`
	UI               *UISettings      `json:"UI,omitempty"`
	Behaviour        *Behaviour       `json:"behaviour,omitempty"`
	OverallFeedback  []FeedbackRange  `json:"overallFeedback,omitempty"`
}

type Answer struct {
	Text     string `json:"text"`
	Correct  bool   `json:"correct"`
	Tipsandanswers bool `json:"tipsAndFeedback,omitempty"`
	Feedback string `json:"feedback,omitempty"`
}

type UISettings struct {
	CheckAnswerButton   string `json:"checkAnswerButton,omitempty"`
	RetryButton         string `json:"retryButton,omitempty"`
	ShowSolutionButton  string `json:"showSolutionButton,omitempty"`
}

type Behaviour struct {
	EnableRetry          bool   `json:"enableRetry,omitempty"`
	EnableSolutionsButton bool   `json:"enableSolutionsButton,omitempty"`
	EnableCheckButton    bool   `json:"enableCheckButton,omitempty"`
	Type                 string `json:"type,omitempty"`
	SinglePoint          bool   `json:"singlePoint,omitempty"`
	RandomAnswers        bool   `json:"randomAnswers,omitempty"`
	PassPercentage       int    `json:"passPercentage,omitempty"`
	RequireAnswer        bool   `json:"requireAnswer,omitempty"`
	ConfirmCheckDialog   bool   `json:"confirmCheckDialog,omitempty"`
	ConfirmRetryDialog   bool   `json:"confirmRetryDialog,omitempty"`
	AutoCheck            bool   `json:"autoCheck,omitempty"`
}