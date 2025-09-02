# H5P Go SDK

[![Build Status][build-status-svg]][build-status-url]
[![Lint Status][lint-status-svg]][lint-status-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![License][license-svg]][license-url]

A Go library for creating, manipulating, and validating H5P (HTML5 Package) content, with support for the official H5P file format and schemas.

## ✨ Features

- 📦 **Full H5P Package Support**: Create and extract `.h5p` ZIP files with proper structure
- 🔒 **Type-Safe Schema Implementation**: Official H5P MultiChoice schema with Go structs
- 🏗️ **Question Set Builder**: Fluent API for building interactive question sets
- ✅ **Validation**: Built-in validation for H5P compliance
- 🎯 **Multiple Question Types**: Support for single-answer and multi-answer questions
- 📋 **Official Schema Compliance**: Uses actual H5P semantics definitions
- 🔄 **JSON Serialization**: Full marshaling/unmarshaling support

## 📥 Installation

```bash
go get github.com/grokify/h5p-go
```

## 🚀 Quick Start

### Creating a Simple Question Set

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Create a question set using the builder pattern
    builder := goh5p.NewQuestionSetBuilder()
    
    answers := []goh5p.Answer{
        goh5p.CreateAnswer("Paris", true),
        goh5p.CreateAnswer("London", false),
        goh5p.CreateAnswer("Berlin", false),
        goh5p.CreateAnswer("Madrid", false),
    }
    
    questionSet, err := builder.
        SetTitle("Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(60).
        SetIntroduction("Welcome to our geography quiz!").
        AddMultipleChoiceQuestion("What is the capital of France?", answers).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Export to JSON
    jsonData, _ := questionSet.ToJSON()
    fmt.Printf("Generated H5P Question Set:\n%s\n", string(jsonData))
}
```

### Creating Typed MultiChoice Questions

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/grokify/h5p-go"
    "github.com/grokify/h5p-go/schemas"
)

func main() {
    // Create a typed MultiChoice question using official H5P schema
    params := &schemas.MultiChoiceParams{
        Question: "What is the capital of France?",
        Answers: []schemas.AnswerOption{
            {
                Text:    "Paris",
                Correct: true,
                TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                    ChosenFeedback: "Correct! Paris is the capital of France.",
                },
            },
            {
                Text:    "London",
                Correct: false,
                TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                    ChosenFeedback: "Incorrect. London is the capital of the UK.",
                },
            },
        },
        Behaviour: &schemas.Behaviour{
            Type:                  "single",
            EnableRetry:           true,
            EnableSolutionsButton: true,
            RandomAnswers:         true,
        },
    }
    
    // Validate the parameters
    if err := params.Validate(); err != nil {
        log.Fatal(err)
    }
    
    // Create the typed question
    question := goh5p.NewMultiChoiceQuestion(params)
    fmt.Printf("Created question: %s\n", question.Params.Question)
}
```

### Loading Question Sets from JSON

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Read JSON file
    jsonData, err := os.ReadFile("my-questionset.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Parse JSON into QuestionSet
    questionSet, err := goh5p.FromJSON(jsonData)
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate the loaded question set
    err = questionSet.Validate()
    if err != nil {
        log.Fatal(err)
    }
    
    // Use the question set
    fmt.Printf("Loaded question set: %s\n", questionSet.Title)
    fmt.Printf("Number of questions: %d\n", len(questionSet.Questions))
    fmt.Printf("Pass percentage: %d%%\n", questionSet.PassPercentage)
    
    // Access individual questions
    for i, question := range questionSet.Questions {
        fmt.Printf("Question %d Library: %s\n", i+1, question.Library)
        
        // For MultiChoice questions, you can access params
        if paramsMap, ok := question.Params.(map[string]interface{}); ok {
            if questionText, exists := paramsMap["question"]; exists {
                fmt.Printf("Question %d Text: %s\n", i+1, questionText)
            }
            
            if answers, exists := paramsMap["answers"]; exists {
                if answersList, ok := answers.([]interface{}); ok {
                    fmt.Printf("Question %d has %d answers\n", i+1, len(answersList))
                }
            }
        }
    }
}
```

### Creating H5P Packages

```go
package main

import (
    "log"
    
    "github.com/grokify/h5p-go"
    "github.com/grokify/h5p-go/schemas"
)

func main() {
    // Create H5P package
    pkg := goh5p.NewH5PPackage()
    
    // Set package definition
    packageDef := &goh5p.PackageDefinition{
        Title:       "My Geography Quiz",
        Language:    "en",
        MainLibrary: "H5P.MultiChoice",
        EmbedTypes:  []string{"div"},
        PreloadedDependencies: []goh5p.LibraryDependency{
            {
                MachineName:  "H5P.MultiChoice",
                MajorVersion: 1,
                MinorVersion: 16,
            },
        },
    }
    pkg.SetPackageDefinition(packageDef)
    
    // Create content
    params := &schemas.MultiChoiceParams{
        Question: "What is 2 + 2?",
        Answers: []schemas.AnswerOption{
            {Text: "3", Correct: false},
            {Text: "4", Correct: true},
            {Text: "5", Correct: false},
        },
    }
    
    content := &goh5p.Content{
        Params: params,
    }
    pkg.SetContent(content)
    
    // Add library
    lib := &goh5p.Library{
        MachineName: "H5P.MultiChoice-1.16",
        Definition: &goh5p.LibraryDefinition{
            Title:        "Multiple Choice",
            MachineName:  "H5P.MultiChoice",
            MajorVersion: 1,
            MinorVersion: 16,
            PatchVersion: 4,
            Runnable:     true,
        },
        Files: map[string][]byte{
            "js/multichoice.js":  []byte("// JavaScript code"),
            "css/multichoice.css": []byte("/* CSS styles */"),
        },
    }
    pkg.AddLibrary(lib)
    
    // Create H5P file
    err := pkg.CreateZipFile("my-quiz.h5p")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Created my-quiz.h5p successfully!")
}
```

## 📚 API Reference

### Question Set Builder

The `QuestionSetBuilder` provides a fluent interface for creating question sets:

```go
builder := goh5p.NewQuestionSetBuilder()

questionSet, err := builder.
    SetTitle("My Quiz").
    SetProgressType("textual").
    SetPassPercentage(80).
    SetIntroduction("Welcome to the quiz!").
    SetStartButtonText("Begin").
    AddMultipleChoiceQuestion("Question text", answers).
    AddOverallFeedback(feedbackRanges).
    Build()
```

### Typed Questions

Create type-safe questions using official H5P schemas:

```go
// MultiChoice question with full type safety
question := goh5p.NewMultiChoiceQuestion(params)

// Convert to generic question for question sets
genericQuestion := question.ToQuestion()
```

### H5P Package Management

Work with complete H5P packages:

```go
// Create new package
pkg := goh5p.NewH5PPackage()

// Load existing package
pkg, err := goh5p.LoadH5PPackage("path/to/file.h5p")

// Export as ZIP file
err = pkg.CreateZipFile("output.h5p")
```

### Validation

Built-in validation ensures H5P compliance:

```go
// Validate question parameters
err := params.Validate()

// Validate complete question set
err = questionSet.Validate()
```

## 🔧 Schema Support

### Official H5P Schemas

The library includes official H5P schemas in the `schemas/` directory:

- `schemas/multichoice_semantics.json` - Official H5P MultiChoice schema
- `schemas/essay_semantics.json` - Official H5P Essay schema
- `schemas/truefalse_semantics.json` - Official H5P True/False schema
- `schemas/multichoice_types.go` - Go struct implementation

### H5P Semantics Format

The library now includes a standardized Go implementation of the H5P semantics format in the `semantics/` directory:

- `semantics/types.go` - Universal Go structs for H5P semantics definition

The H5P semantics format is a standardized JSON schema used across all H5P content types to define form fields, validation rules, and UI behaviors. Our Go implementation supports all semantic field types including:

- **Field Types**: `group`, `text`, `boolean`, `select`, `list`, `number`, `library`
- **Common Attributes**: `name`, `type`, `label`, `importance`, `optional`, `default`
- **Advanced Features**: Conditional visibility (`showWhen`), nested structures, validation rules
- **UI Configuration**: Widget types, placeholder text, field descriptions

These are defined by H5P here:

1. True-False
    1. [H5P `github.com/h5p/h5p-true-false` project](https://github.com/h5p/h5p-true-false)
    1. [H5P true-false `semantics.json` on GitHub](https://github.com/h5p/h5p-true-false/blob/master/semantics.json)
1. Multiple-Choice
    1. [H5P `github.com/h5p/h5p-multi-choice` project](https://github.com/h5p/h5p-multi-choice)
    1. [H5P multichoice `semantics.json` on GitHub](https://github.com/h5p/h5p-multi-choice/blob/master/semantics.json)
1. Essay
    1. [H5P `github.com/otacke/h5p-essay` project](https://github.com/otacke/h5p-essay)
    1. [H5P essay `semantics.json` on GitHub](https://github.com/otacke/h5p-essay/blob/master/semantics.json)

### Question Types

Currently supported question types:

- **MultiChoice**: Single and multiple answer questions with rich feedback
- **Question Sets**: Collections of questions with overall scoring

### MultiChoice Features

- Single-answer (radio buttons) and multi-answer (checkboxes) questions
- Individual answer feedback and tips
- Overall feedback based on score ranges
- Customizable UI text and behavior settings
- Media support (images, videos)
- Randomizable answer order

## 📁 File Structure

```
h5p-go/
├── schemas/                          # Official H5P schemas
│   ├── multichoice_semantics.json   # H5P MultiChoice schema
│   ├── essay_semantics.json         # H5P Essay schema
│   ├── truefalse_semantics.json     # H5P True/False schema
│   └── multichoice_types.go         # Go types for MultiChoice
├── semantics/                       # H5P semantics format
│   └── types.go                     # Universal semantics Go structs
├── testdata/                        # Test data files
│   ├── content.json                 # Sample content
│   ├── h5p.json                     # Sample package definition
│   └── library.json                 # Sample library definition
├── builder.go                       # Question set builder
├── h5p_package.go                   # H5P package management
├── questionset.go                   # Core question set types
└── README.md
```

## 📄 H5P File Format

The library supports the complete H5P file format:

```
example.h5p (ZIP file)
├── h5p.json                    # Package metadata
├── content/
│   └── content.json           # Content parameters
└── H5P.LibraryName-1.0/      # Library folders
    ├── library.json          # Library definition
    ├── semantics.json        # Content schema
    ├── js/                   # JavaScript files
    └── css/                  # Stylesheet files
```

## 💡 Examples

See the `*_test.go` files for comprehensive examples:

- `questionset_test.go` - Basic question set usage
- `typed_question_test.go` - Typed question examples  
- `h5p_package_test.go` - Complete H5P package creation

## 🧪 Testing

Run the test suite:

```bash
go test -v
```

Run specific test categories:

```bash
go test -v -run TestH5PPackage        # H5P package tests
go test -v -run TestTyped             # Typed question tests
go test -v -run TestSemanticsValidation # Schema validation tests
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📏 Standards Compliance

This library follows the official H5P specifications:

- [H5P File Format Specification](https://h5p.org/documentation/developers/h5p-specification)
- [H5P Semantics Definition](https://h5p.org/semantics)
- [H5P MultiChoice Content Type](https://github.com/h5p/h5p-multi-choice)

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [H5P Group](https://h5p.org) for the H5P framework and specifications
- Official H5P schemas and content type definitions
- The Go community for excellent tooling and libraries

## 🗺️ Roadmap

- [ ] Support for additional H5P content types (Essay, True/False, etc.)
- [ ] Advanced validation and error reporting
- [ ] H5P content migration tools
- [ ] Integration with H5P hosting platforms
- [ ] Performance optimizations for large content packages

 [build-status-svg]: https://github.com/grokify/h5p-go/actions/workflows/ci.yaml/badge.svg?branch=main
 [build-status-url]: https://github.com/grokify/h5p-go/actions/workflows/ci.yaml
 [lint-status-svg]: https://github.com/grokify/h5p-go/actions/workflows/lint.yaml/badge.svg?branch=main
 [lint-status-url]: https://github.com/grokify/h5p-go/actions/workflows/lint.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/h5p-go
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/h5p-go
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/h5p-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/h5p-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/h5p-go/blob/master/LICENSE
 [used-by-svg]: https://sourcegraph.com/github.com/grokify/h5p-go/-/badge.svg
 [used-by-url]: https://sourcegraph.com/github.com/grokify/h5p-go?badge
