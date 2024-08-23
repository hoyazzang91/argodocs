package workflow

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ParseFiles accepts a glob pattern for argo workflow template files, parses them and returns them as TemplateFile structs
func ParseFiles(root string) ([]*TemplateFile, error) {
	var result []*TemplateFile

	//dir, err := os.ReadDir(path)
	//if err != nil {
	//	panic(err)
	//}

	var paths []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileInfo, _ := os.Stat(path)
		if !info.IsDir() && fileInfo.Size() > 0 {
			paths = append(paths, path)
		}
		return nil
	})

	for _, f := range paths {
		var templateFile *TemplateFile
		if strings.Contains(filepath.Ext(f), "yaml") {
			templateFile, err = parseFile(root, f)
			if err != nil {
				continue
			}
			result = append(result, templateFile)
		}
	}

	return result, nil
}

// parseFile parses a single argo workflow template file and returns that as a TemplateFile object
func parseFile(path string, name string) (*TemplateFile, error) {
	yamlNode := yaml.Node{}
	filePath := filepath.Dir(name)
	fileName := name
	url := "https://raw.githubusercontent.com/argoproj-labs/argo-workflows-catalog/master/" + fileName

	yamlFileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFileContent, &yamlNode)
	if err != nil {
		return nil, err
	}

	var templateFile *TemplateFile
	templateFile, err = parseTemplateFile(&yamlNode)
	if err != nil {
		return nil, err
	}

	templateFile.Data = yamlFileContent
	templateFile.Name = strings.Split(filePath, string(filepath.Separator))[1]
	templateFile.FilePath = path + "/" + templateFile.Name
	templateFile.LastUpdatedAt = time.Now().Format(time.RFC850)
	templateFile.Command = "kubectl apply -f " + url

	_, err = os.Open(filePath + "/icon.png")
	if err == nil {
		templateFile.Icon = strings.Replace(filePath, "_", "", 1) + "/icon.png"
	}

	if templateFile.EntrypointTemplate == "" {
		templateFile.EntrypointTemplate = "nil"
	}

	return templateFile, nil
}
