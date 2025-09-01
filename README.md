# go-h5p

A comprehensive Go library for creating, manipulating, and validating H5P (HTML5 Package) content, with full support for the official H5P file format and schemas.

[![Go Reference](https://pkg.go.dev/badge/github.com/grokify/go-h5p.svg)](https://pkg.go.dev/github.com/grokify/go-h5p)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/go-h5p)](https://goreportcard.com/report/github.com/grokify/go-h5p)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

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
go get github.com/grokify/go-h5p
```

## 🚀 Quick Start

### Creating a Simple Question Set

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/grokify/go-h5p"
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
    
    "github.com/grokify/go-h5p"
    "github.com/grokify/go-h5p/schemas"
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

### Creating H5P Packages

```go
package main

import (
    "log"
    
    "github.com/grokify/go-h5p"
    "github.com/grokify/go-h5p/schemas"
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

- `schemas/multichoice-semantics.json` - Official H5P MultiChoice schema
- `schemas/multichoice_types.go` - Go struct implementation

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
go-h5p/
├── schemas/                          # Official H5P schemas
│   ├── multichoice-semantics.json   # H5P MultiChoice schema
│   └── multichoice_types.go         # Go types for MultiChoice
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