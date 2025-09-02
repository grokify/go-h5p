package h5p

import (
	"fmt"
	"log"
)

func ExampleQuestionSetBuilder() {
	builder := NewQuestionSetBuilder()

	answers := []Answer{
		CreateAnswer("Paris", true),
		CreateAnswer("London", false),
		CreateAnswer("Berlin", false),
		CreateAnswer("Madrid", false),
	}

	feedbackRanges := []FeedbackRange{
		CreateFeedbackRange(0, 50, "You need more practice!"),
		CreateFeedbackRange(51, 80, "Good job!"),
		CreateFeedbackRange(81, 100, "Excellent work!"),
	}

	questionSet, err := builder.
		SetTitle("Geography Quiz").
		SetProgressType("textual").
		SetPassPercentage(60).
		SetIntroduction("Welcome to our geography quiz!").
		SetStartButtonText("Start Quiz").
		AddMultipleChoiceQuestion("What is the capital of France?", answers).
		AddOverallFeedback(feedbackRanges).
		Build()

	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := questionSet.ToJSON()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Generated H5P Question Set:\n%s\n", string(jsonData))

	loadedQuestionSet, err := FromJSON(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	err = loadedQuestionSet.Validate()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Question set validated successfully!\n")
	fmt.Printf("Number of questions: %d\n", len(loadedQuestionSet.Questions))
	fmt.Printf("Pass percentage: %d%%\n", loadedQuestionSet.PassPercentage)
}
