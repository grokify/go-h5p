package goh5p

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type H5PPackage struct {
	PackageDefinition *PackageDefinition `json:"-"`
	Content           *Content           `json:"-"`
	Libraries         []*Library         `json:"-"`
}

type PackageDefinition struct {
	Title                 string              `json:"title"`
	Language              string              `json:"language"`
	MainLibrary           string              `json:"mainLibrary"`
	EmbedTypes            []string            `json:"embedTypes"`
	License               string              `json:"license,omitempty"`
	DefaultLanguage       string              `json:"defaultLanguage,omitempty"`
	Author                string              `json:"author,omitempty"`
	PreloadedDependencies []LibraryDependency `json:"preloadedDependencies"`
	EditorDependencies    []LibraryDependency `json:"editorDependencies,omitempty"`
}

type LibraryDependency struct {
	MachineName    string `json:"machineName"`
	MajorVersion   int    `json:"majorVersion"`
	MinorVersion   int    `json:"minorVersion"`
}

type Content struct {
	QuestionSet *QuestionSet `json:"questionSet,omitempty"`
	Params      interface{}  `json:",omitempty"`
}

type Library struct {
	Definition *LibraryDefinition `json:"-"`
	Semantics  interface{}        `json:"-"`
	MachineName string            `json:"-"`
	Files      map[string][]byte  `json:"-"`
}

type LibraryDefinition struct {
	Title         string              `json:"title"`
	MachineName   string              `json:"machineName"`
	MajorVersion  int                 `json:"majorVersion"`
	MinorVersion  int                 `json:"minorVersion"`
	PatchVersion  int                 `json:"patchVersion"`
	Runnable      bool                `json:"runnable"`
	Author        string              `json:"author,omitempty"`
	License       string              `json:"license,omitempty"`
	Description   string              `json:"description,omitempty"`
	PreloadedJs   []FileReference     `json:"preloadedJs,omitempty"`
	PreloadedCss  []FileReference     `json:"preloadedCss,omitempty"`
	DropLibraryCss []FileReference    `json:"dropLibraryCss,omitempty"`
	Dependencies  []LibraryDependency `json:"preloadedDependencies,omitempty"`
}

type FileReference struct {
	Path string `json:"path"`
}

func NewH5PPackage() *H5PPackage {
	return &H5PPackage{
		Libraries: make([]*Library, 0),
	}
}

func (pkg *H5PPackage) SetPackageDefinition(def *PackageDefinition) {
	pkg.PackageDefinition = def
}

func (pkg *H5PPackage) SetContent(content *Content) {
	pkg.Content = content
}

func (pkg *H5PPackage) AddLibrary(lib *Library) {
	pkg.Libraries = append(pkg.Libraries, lib)
}

func (pkg *H5PPackage) CreateZipFile(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	if err := pkg.writeToZip(zipWriter); err != nil {
		return fmt.Errorf("failed to write package to zip: %w", err)
	}

	return nil
}

func (pkg *H5PPackage) writeToZip(zipWriter *zip.Writer) error {
	if pkg.PackageDefinition != nil {
		h5pJSON, err := json.MarshalIndent(pkg.PackageDefinition, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal h5p.json: %w", err)
		}
		if err := writeFileToZip(zipWriter, "h5p.json", h5pJSON); err != nil {
			return err
		}
	}

	if pkg.Content != nil {
		contentJSON, err := json.MarshalIndent(pkg.Content, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal content.json: %w", err)
		}
		if err := writeFileToZip(zipWriter, "content/content.json", contentJSON); err != nil {
			return err
		}
	}

	for _, lib := range pkg.Libraries {
		if lib.Definition != nil {
			libJSON, err := json.MarshalIndent(lib.Definition, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal library.json: %w", err)
			}
			libPath := fmt.Sprintf("%s/library.json", lib.MachineName)
			if err := writeFileToZip(zipWriter, libPath, libJSON); err != nil {
				return err
			}
		}

		if lib.Semantics != nil {
			semJSON, err := json.MarshalIndent(lib.Semantics, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal semantics.json: %w", err)
			}
			semPath := fmt.Sprintf("%s/semantics.json", lib.MachineName)
			if err := writeFileToZip(zipWriter, semPath, semJSON); err != nil {
				return err
			}
		}

		for filePath, fileData := range lib.Files {
			fullPath := fmt.Sprintf("%s/%s", lib.MachineName, filePath)
			if err := writeFileToZip(zipWriter, fullPath, fileData); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeFileToZip(zipWriter *zip.Writer, filename string, data []byte) error {
	writer, err := zipWriter.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create zip entry for %s: %w", filename, err)
	}
	
	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %w", filename, err)
	}
	
	return nil
}

func LoadH5PPackage(filePath string) (*H5PPackage, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open H5P file: %w", err)
	}
	defer reader.Close()

	pkg := NewH5PPackage()
	
	for _, file := range reader.File {
		if err := pkg.processZipFile(file); err != nil {
			return nil, fmt.Errorf("failed to process file %s: %w", file.Name, err)
		}
	}

	return pkg, nil
}

func (pkg *H5PPackage) processZipFile(file *zip.File) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	switch {
	case file.Name == "h5p.json":
		var pkgDef PackageDefinition
		if err := json.Unmarshal(data, &pkgDef); err != nil {
			return err
		}
		pkg.PackageDefinition = &pkgDef

	case file.Name == "content/content.json":
		var content Content
		if err := json.Unmarshal(data, &content); err != nil {
			return err
		}
		pkg.Content = &content

	case strings.HasSuffix(file.Name, "/library.json"):
		libName := filepath.Dir(file.Name)
		lib := pkg.findOrCreateLibrary(libName)
		
		var libDef LibraryDefinition
		if err := json.Unmarshal(data, &libDef); err != nil {
			return err
		}
		lib.Definition = &libDef

	case strings.HasSuffix(file.Name, "/semantics.json"):
		libName := filepath.Dir(file.Name)
		lib := pkg.findOrCreateLibrary(libName)
		
		var semantics interface{}
		if err := json.Unmarshal(data, &semantics); err != nil {
			return err
		}
		lib.Semantics = semantics

	default:
		if strings.Contains(file.Name, "/") {
			libName := strings.Split(file.Name, "/")[0]
			if pkg.isLibraryDirectory(libName) {
				lib := pkg.findOrCreateLibrary(libName)
				if lib.Files == nil {
					lib.Files = make(map[string][]byte)
				}
				relativePath := strings.TrimPrefix(file.Name, libName+"/")
				lib.Files[relativePath] = data
			}
		}
	}

	return nil
}

func (pkg *H5PPackage) findOrCreateLibrary(machineName string) *Library {
	for _, lib := range pkg.Libraries {
		if lib.MachineName == machineName {
			return lib
		}
	}
	
	lib := &Library{
		MachineName: machineName,
		Files:       make(map[string][]byte),
	}
	pkg.Libraries = append(pkg.Libraries, lib)
	return lib
}

func (pkg *H5PPackage) isLibraryDirectory(name string) bool {
	return strings.HasPrefix(name, "H5P.")
}