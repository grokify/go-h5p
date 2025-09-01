package goh5p

import (
	"testing"
)

func TestQuestionSetBuilder(t *testing.T) {
	builder := NewQuestionSetBuilder()
	
	answers := []Answer{
		CreateAnswer("Correct answer", true),
		CreateAnswer("Wrong answer", false),
	}
	
	qs, err := builder.
		SetTitle("Test Quiz").
		SetPassPercentage(70).
		AddMultipleChoiceQuestion("Test question?", answers).
		Build()
		
	if err != nil {
		t.Fatalf("Failed to build question set: %v", err)
	}
	
	if qs.Title != "Test Quiz" {
		t.Errorf("Expected title 'Test Quiz', got '%s'", qs.Title)
	}
	
	if qs.PassPercentage != 70 {
		t.Errorf("Expected pass percentage 70, got %d", qs.PassPercentage)
	}
	
	if len(qs.Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(qs.Questions))
	}
}

func TestQuestionSetValidation(t *testing.T) {
	qs := &QuestionSet{
		PassPercentage: 150,
		Questions:      []Question{},
	}
	
	err := qs.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid pass percentage and no questions")
	}
}

func TestJSONMarshalUnmarshal(t *testing.T) {
	builder := NewQuestionSetBuilder()
	
	answers := []Answer{
		CreateAnswer("Answer 1", true),
		CreateAnswer("Answer 2", false),
	}
	
	original, err := builder.
		SetTitle("JSON Test").
		AddMultipleChoiceQuestion("Test question?", answers).
		Build()
		
	if err != nil {
		t.Fatalf("Failed to build question set: %v", err)
	}
	
	jsonData, err := original.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}
	
	restored, err := FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}
	
	if restored.Title != original.Title {
		t.Errorf("Title mismatch after JSON round-trip: expected '%s', got '%s'", original.Title, restored.Title)
	}
	
	if len(restored.Questions) != len(original.Questions) {
		t.Errorf("Question count mismatch after JSON round-trip: expected %d, got %d", len(original.Questions), len(restored.Questions))
	}
}

func TestCreateAnswerFunctions(t *testing.T) {
	answer := CreateAnswer("Test answer", true)
	if answer.Text != "Test answer" || !answer.Correct {
		t.Error("CreateAnswer did not create answer correctly")
	}
	
	answerWithFeedback := CreateAnswerWithFeedback("Test answer", false, "Feedback text")
	if answerWithFeedback.Feedback != "Feedback text" {
		t.Error("CreateAnswerWithFeedback did not set feedback correctly")
	}
}

func TestFeedbackRange(t *testing.T) {
	fr := CreateFeedbackRange(80, 100, "Excellent!")
	if fr.From != 80 || fr.To != 100 || fr.Text != "Excellent!" {
		t.Error("CreateFeedbackRange did not create range correctly")
	}
}