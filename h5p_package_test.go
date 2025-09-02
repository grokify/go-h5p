package h5p

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestH5PPackageCreationAndExtraction(t *testing.T) {
	// Load test data files
	semanticsData, err := os.ReadFile("schemas/multichoice-semantics.json")
	if err != nil {
		t.Fatalf("Failed to read multichoice-semantics.json: %v", err)
	}

	contentData, err := os.ReadFile("testdata/content.json")
	if err != nil {
		t.Fatalf("Failed to read content.json: %v", err)
	}

	h5pData, err := os.ReadFile("testdata/h5p.json")
	if err != nil {
		t.Fatalf("Failed to read h5p.json: %v", err)
	}

	libraryData, err := os.ReadFile("testdata/library.json")
	if err != nil {
		t.Fatalf("Failed to read library.json: %v", err)
	}

	// Parse the files
	var semantics interface{}
	if err := json.Unmarshal(semanticsData, &semantics); err != nil {
		t.Fatalf("Failed to parse semantics.json: %v", err)
	}

	var content Content
	if err := json.Unmarshal(contentData, &content); err != nil {
		t.Fatalf("Failed to parse content.json: %v", err)
	}

	var packageDef PackageDefinition
	if err := json.Unmarshal(h5pData, &packageDef); err != nil {
		t.Fatalf("Failed to parse h5p.json: %v", err)
	}

	var libraryDef LibraryDefinition
	if err := json.Unmarshal(libraryData, &libraryDef); err != nil {
		t.Fatalf("Failed to parse library.json: %v", err)
	}

	// Create H5P package
	pkg := NewH5PPackage()
	pkg.SetPackageDefinition(&packageDef)
	pkg.SetContent(&content)

	// Create library with semantics and definition
	lib := &Library{
		MachineName: "H5P.MultiChoice-1.16",
		Definition:  &libraryDef,
		Semantics:   semantics,
		Files:       make(map[string][]byte),
	}
	// Add some dummy JS/CSS files
	lib.Files["js/multichoice.js"] = []byte("// MultiChoice JavaScript code")
	lib.Files["css/multichoice.css"] = []byte("/* MultiChoice CSS styles */")

	pkg.AddLibrary(lib)

	// Create temporary H5P file
	tempFile := filepath.Join(os.TempDir(), "test_package.h5p")
	defer os.Remove(tempFile)

	err = pkg.CreateZipFile(tempFile)
	if err != nil {
		t.Fatalf("Failed to create H5P package: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("H5P file was not created")
	}

	// Load the package back
	loadedPkg, err := LoadH5PPackage(tempFile)
	if err != nil {
		t.Fatalf("Failed to load H5P package: %v", err)
	}

	// Validate package definition
	if loadedPkg.PackageDefinition == nil {
		t.Fatal("Package definition is nil")
	}
	if loadedPkg.PackageDefinition.Title != "Geography Quiz - Capital of France" {
		t.Errorf("Expected title 'Geography Quiz - Capital of France', got '%s'", loadedPkg.PackageDefinition.Title)
	}
	if loadedPkg.PackageDefinition.MainLibrary != "H5P.MultiChoice" {
		t.Errorf("Expected main library 'H5P.MultiChoice', got '%s'", loadedPkg.PackageDefinition.MainLibrary)
	}

	// Validate content
	if loadedPkg.Content == nil {
		t.Fatal("Content is nil")
	}

	// Validate libraries
	if len(loadedPkg.Libraries) != 1 {
		t.Fatalf("Expected 1 library, got %d", len(loadedPkg.Libraries))
	}

	loadedLib := loadedPkg.Libraries[0]
	if loadedLib.Definition == nil {
		t.Fatal("Library definition is nil")
	}
	if loadedLib.Definition.MachineName != "H5P.MultiChoice" {
		t.Errorf("Expected machine name 'H5P.MultiChoice', got '%s'", loadedLib.Definition.MachineName)
	}
	if loadedLib.Semantics == nil {
		t.Fatal("Library semantics is nil")
	}

	// Validate library files
	if len(loadedLib.Files) != 2 {
		t.Errorf("Expected 2 library files, got %d", len(loadedLib.Files))
	}
	if _, exists := loadedLib.Files["js/multichoice.js"]; !exists {
		t.Error("JavaScript file not found in library")
	}
	if _, exists := loadedLib.Files["css/multichoice.css"]; !exists {
		t.Error("CSS file not found in library")
	}

	t.Log("H5P package creation and extraction test completed successfully")
}

func TestSemanticsValidation(t *testing.T) {
	semanticsData, err := os.ReadFile("schemas/multichoice-semantics.json")
	if err != nil {
		t.Fatalf("Failed to read multichoice-semantics.json: %v", err)
	}

	var semantics []interface{}
	if err := json.Unmarshal(semanticsData, &semantics); err != nil {
		t.Fatalf("Failed to parse semantics.json: %v", err)
	}

	// Validate semantics structure
	if len(semantics) == 0 {
		t.Fatal("Semantics array is empty")
	}

	// Check for required fields in semantics
	foundQuestion := false
	foundAnswers := false
	foundBehaviour := false

	for _, field := range semantics {
		if fieldMap, ok := field.(map[string]interface{}); ok {
			if name, exists := fieldMap["name"]; exists {
				switch name {
				case "question":
					foundQuestion = true
				case "answers":
					foundAnswers = true
				case "behaviour":
					foundBehaviour = true
				}
			}
		}
	}

	if !foundQuestion {
		t.Error("Question field not found in semantics")
	}
	if !foundAnswers {
		t.Error("Answers field not found in semantics")
	}
	if !foundBehaviour {
		t.Error("Behaviour field not found in semantics")
	}

	t.Log("Semantics validation completed successfully")
}

func TestContentValidation(t *testing.T) {
	contentData, err := os.ReadFile("testdata/content.json")
	if err != nil {
		t.Fatalf("Failed to read content.json: %v", err)
	}

	var content map[string]interface{}
	if err := json.Unmarshal(contentData, &content); err != nil {
		t.Fatalf("Failed to parse content.json: %v", err)
	}

	// Validate required content fields
	if _, exists := content["question"]; !exists {
		t.Error("Question field not found in content")
	}
	if _, exists := content["answers"]; !exists {
		t.Error("Answers field not found in content")
	}
	if _, exists := content["behaviour"]; !exists {
		t.Error("Behaviour field not found in content")
	}

	// Validate answers structure
	if answers, ok := content["answers"].([]interface{}); ok {
		if len(answers) != 4 {
			t.Errorf("Expected 4 answers, got %d", len(answers))
		}

		correctAnswers := 0
		for _, answer := range answers {
			if answerMap, ok := answer.(map[string]interface{}); ok {
				if correct, exists := answerMap["correct"]; exists && correct.(bool) {
					correctAnswers++
				}
			}
		}

		if correctAnswers != 1 {
			t.Errorf("Expected 1 correct answer for single choice, got %d", correctAnswers)
		}
	} else {
		t.Error("Answers is not an array")
	}

	// Validate behaviour
	if behaviour, ok := content["behaviour"].(map[string]interface{}); ok {
		if questionType, exists := behaviour["type"]; exists {
			if questionType != "single" {
				t.Errorf("Expected question type 'single', got '%v'", questionType)
			}
		} else {
			t.Error("Question type not specified in behaviour")
		}
	} else {
		t.Error("Behaviour is not an object")
	}

	t.Log("Content validation completed successfully")
}
