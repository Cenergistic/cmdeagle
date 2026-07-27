package bundle

import (
	"bytes"
	"embed"
	"runtime"

	"github.com/migsc/cmdeagle/file"
	"github.com/migsc/cmdeagle/types"

	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/afero"
)

//go:embed *
var PackageFS embed.FS

var StagingDirName = ".tmp"
var MainTemplateFileName = "main.template.go"
var MainTemplateReplacements = map[string][]byte{
	// The package switcharoo
	"package bundle": []byte("package main"),

	// Rename the template's intended main function
	"main_template": []byte("main"),

	// Set debug mode based on flag
	"var LOG_LEVEL = log.InfoLevel": []byte("var LOG_LEVEL = log.InfoLevel"), // default value
}

var packageSrcDirPath string
var bundleStagingDirPath string

func GetPackageSrcDirPath() string {
	log.Debug("Getting package src dir path:")
	if packageSrcDirPath != "" {
		log.Debug("Package src dir path already set:", "path", packageSrcDirPath)
		return packageSrcDirPath
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("failed to get current file path")
	}

	log.Debug("Package src dir path not set, setting it now:", "path", filepath.Dir(filename))

	packageSrcDirPath = filepath.Dir(filename)

	return packageSrcDirPath
}

func SetupStagingDir() (string, error) {
	// TODO possible to move into file package? Or is it reliant on executing from the bundle package?
	tempDirPath, err := afero.TempDir(afero.NewOsFs(), "", "cmdeagle")
	if err != nil {
		return "", err
	}

	log.Debug("Using bundle staging directory path for bundle of: %s\n", tempDirPath)

	bundleStagingDirPath = tempDirPath

	return tempDirPath, nil
}

func GetMainTemplateContent() ([]byte, error) {
	mainTemplateContent, err := PackageFS.ReadFile(MainTemplateFileName)
	if err != nil {
		return nil, fmt.Errorf("error reading template: %w", err)
	}
	return mainTemplateContent, nil
}

func InterpolateMainContent(content []byte) []byte {
	for old, new := range MainTemplateReplacements {
		content = bytes.ReplaceAll(content, []byte(old), new)
	}

	// TODO: We need to embed a filesystem of the current working directory (where main.go will be)?

	return content
}

func SetupMainFile(path string) (string, error) {
	mainTemplateContent, err := GetMainTemplateContent()
	if err != nil {
		return "", fmt.Errorf("error reading template: %w", err)
	}

	for old, new := range MainTemplateReplacements {
		mainTemplateContent = bytes.ReplaceAll(mainTemplateContent, []byte(old), new)
	}

	// Write processed template to main.go
	mainFilePath := filepath.Join(path, "main.go")
	if err := os.WriteFile(mainFilePath, mainTemplateContent, 0644); err != nil {
		return "", fmt.Errorf("error writing main.go: %w", err)
	}

	return mainFilePath, nil
}

func CopyIncludedFiles(config *types.CmdeagleConfig, command *types.CommandDefinition, namespace []string, targetDirPath string) error {
	log.Debug("Copying included files",
		"command", command.Name,
		"namespace", namespace,
		"targetDirPath", targetDirPath,
	)

	if len(command.Includes) == 0 {
		return nil
	}

	log.Info("processing includes",
		"command", command.Name,
		"includes", command.Includes,
	)

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, ns := range namespace {
		targetDirPath = filepath.Join(targetDirPath, ns)
	}

	for _, includePath := range command.Includes {
		log.Info("including bundle",
			"from", includePath,
			"to", targetDirPath,
		)
		if err := copyIncludedFile(filepath.Join(currentDir, includePath), targetDirPath); err != nil {
			return err
		}
	}

	return nil
}

func copyIncludedFile(includedFilePath string, targetDir string) error {
	log.Info("including bundle",
		"from", includedFilePath,
		"to", targetDir,
	)

	expandedPath, err := file.ExpandPath(includedFilePath)
	if err != nil {
		return fmt.Errorf("could not expand the path: %s\n%v", includedFilePath, err)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("could not create target directory %s: %w", targetDir, err)
	}

	// Copy the source into the target directory (targetDir/<base>), preserving
	// file modes. Uses a portable Go implementation rather than shelling out to
	// `cp`, which is unavailable on Windows.
	dst := filepath.Join(targetDir, filepath.Base(expandedPath))
	if err := file.Copy(expandedPath, dst); err != nil {
		return fmt.Errorf("could not copy %s to %s: %w", expandedPath, dst, err)
	}

	return nil
}
