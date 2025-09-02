# H5P Go SDK

[![Build Status](https://github.com/grokify/h5p-go/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/grokify/h5p-go/actions/workflows/ci.yaml)
[![Lint Status](https://github.com/grokify/h5p-go/actions/workflows/lint.yaml/badge.svg?branch=main)](https://github.com/grokify/h5p-go/actions/workflows/lint.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/h5p-go)](https://goreportcard.com/report/github.com/grokify/h5p-go)
[![Docs](https://pkg.go.dev/badge/github.com/grokify/h5p-go)](https://pkg.go.dev/github.com/grokify/h5p-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/grokify/h5p-go/blob/master/LICENSE)

A Go library for creating, manipulating, and validating H5P (HTML5 Package) content with support for the official H5P file format and schemas.

## ✨ Features

- 📦 **Full H5P Package Support** - Create and extract `.h5p` ZIP files
- 🔒 **Type-Safe Schema Implementation** - Official H5P content type schemas  
- 🏗️ **Question Set Builder** - Fluent API for building interactive content
- ✅ **Validation** - Built-in H5P compliance validation
- 🎯 **Multiple Question Types** - Support for various H5P content types
- 🔄 **JSON Serialization** - Complete marshaling/unmarshaling support

## 🚀 Quick Start

### Installation

```bash
go get github.com/grokify/h5p-go
```

### Create Your First Quiz

```go
package main

import (
    "fmt"
    "log"
    "github.com/grokify/h5p-go"
)

func main() {
    // Create answers
    answers := []h5p.Answer{
        h5p.CreateAnswer("Paris", true),
        h5p.CreateAnswer("London", false),
        h5p.CreateAnswer("Berlin", false),
    }
    
    // Build question set
    questionSet, err := h5p.NewQuestionSetBuilder().
        SetTitle("Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(60).
        AddMultipleChoiceQuestion("What is the capital of France?", answers).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Export to JSON
    jsonData, _ := questionSet.ToJSON()
    fmt.Printf("Generated H5P content:\n%s\n", string(jsonData))
}
```

### Using Type-Safe Schemas

```go
import "github.com/grokify/h5p-go/schemas"

// Create strongly-typed content
params := &schemas.MultiChoiceParams{
    Question: "What is 2 + 2?",
    Answers: []schemas.AnswerOption{
        {Text: "4", Correct: true},
        {Text: "5", Correct: false},
    },
    Behaviour: &schemas.Behaviour{
        Type: "single",
        EnableRetry: true,
    },
}

question := h5p.NewMultiChoiceQuestion(params)
```

## 📚 Documentation

**[Full Documentation →](https://grokify.github.io/h5p-go/)**

### Quick Links

- **[Installation Guide](https://grokify.github.io/h5p-go/getting-started/installation/)** - Detailed setup instructions
- **[Quick Start Tutorial](https://grokify.github.io/h5p-go/getting-started/quick-start/)** - Step-by-step examples
- **[API Reference](https://grokify.github.io/h5p-go/api/core-types/)** - Complete API documentation
- **[Examples](https://grokify.github.io/h5p-go/examples/basic/)** - Real-world usage examples

### Key Topics

- [Question Sets](https://grokify.github.io/h5p-go/guide/question-sets/) - Building interactive question collections
- [Typed Questions](https://grokify.github.io/h5p-go/guide/typed-questions/) - Using type-safe H5P schemas
- [H5P Packages](https://grokify.github.io/h5p-go/guide/h5p-packages/) - Creating complete `.h5p` files
- [Semantics API](https://grokify.github.io/h5p-go/api/semantics/) - Working with H5P semantics format

## 🏗️ Architecture

```
h5p-go/
├── schemas/          # Official H5P content type schemas
├── semantics/        # Universal H5P semantics format  
├── builder.go        # Fluent API for content creation
├── questionset.go    # Core question set functionality
└── h5p_package.go    # Complete H5P package management
```

## 🧪 Testing

```bash
go test ./...
```

Run specific tests:
```bash
go test -v -run TestQuestionSet    # Question set tests
go test -v -run TestTyped          # Typed schema tests
go test -v -run TestH5PPackage     # Package management tests
```

## 📏 Standards Compliance

This library implements the official H5P specifications:

- [H5P File Format](https://h5p.org/documentation/developers/h5p-specification)
- [H5P Semantics](https://h5p.org/semantics) 
- [H5P Content Types](https://github.com/h5p/)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

See our [Contributing Guide](https://grokify.github.io/h5p-go/development/contributing/) for detailed information.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [H5P Group](https://h5p.org) for the H5P framework and specifications
- The Go community for excellent tooling and libraries