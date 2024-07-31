package workflow

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ParseFiles accepts a glob pattern for argo workflow template files, parses them and returns them as TemplateFile structs
func ParseFiles(path string) ([]*TemplateFile, error) {
	var result []*TemplateFile

	dir, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, f := range dir {
		var templateFile *TemplateFile
		templateFile, err = parseFile(path, f.Name())
		if err != nil {
			continue
		}
		result = append(result, templateFile)
	}

	return result, nil
}

// parseFile parses a single argo workflow template file and returns that as a TemplateFile object
func parseFile(path string, name string) (*TemplateFile, error) {
	yamlNode := yaml.Node{}
	filePath := path + "/" + name
	fileName := filePath + "/manifests.yaml"
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

	templateFile.Name = name
	templateFile.FilePath = filePath
	templateFile.LastUpdatedAt = time.Now().Format(time.RFC850)
	templateFile.Command = "kubectl apply -f " + url

	_, err = os.Open(filePath + "/icon.png")
	if err == nil {
		templateFile.Icon = filePath + "/icon.png"
	}

	if templateFile.EntrypointTemplate == "" {
		templateFile.EntrypointTemplate = "nil"
	}

	return templateFile, nil
}
