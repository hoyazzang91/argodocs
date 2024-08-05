package mdgen

import (
	"strconv"
	"strings"

	"github.com/hoyazzang91/argodocs/markdown"
	"github.com/hoyazzang91/argodocs/workflow"
)

// GetMdDoc transforms a workflow.TemplateFile into an opinionated markdown file
func GetMdDoc(templateFile *workflow.TemplateFile) (*markdown.Doc, error) {
	templateTypes := map[workflow.TemplateType]string{
		workflow.CONTAINER_SET_TEMPLATE: "Container Set",
		workflow.CONTAINER_TEMPLATE:     "Container",
		workflow.DAG_TEMPLATE:           "DAG",
		workflow.DATA_TEMPLATE:          "DATA",
		workflow.HTTP_TEMPLATE:          "HTTP",
		workflow.PLUGIN_TEMPLATE:        "Plugin",
		workflow.SCRIPT_TEMPLATE:        "Script",
	}
	md := markdown.NewDoc()

	md.Write("---\n")
	md.Write("title: " + templateFile.Name + " \n")
	md.Write("tags: " + strings.Join(templateFile.Tags, ",") + " \n")
	md.Write("version: " + templateFile.Version + " \n")
	md.Write("icon: " + templateFile.Icon + " \n")
	md.Writeln("layout: page")
	md.Write("descriptoin: " + strings.Replace(templateFile.Description, "\n", " ", -1) + "\n")
	md.Write("---\n")
	md.WriteHeader(templateFile.Name, 1)
	md.Write("\n")
	md.Write("\n")
	table := markdown.NewTable(1, 4)
	table.SetTableTitle(0, "Kind")
	table.SetTableTitle(1, "Version")
	table.SetTableTitle(2, "Entrypoint Template")
	table.SetTableTitle(3, "Last Updated At")
	table.SetTableContent(0, 0, templateFile.Kind)
	table.SetTableContent(0, 1, templateFile.Version)
	table.SetTableContent(0, 2, templateFile.EntrypointTemplate)
	table.SetTableContent(0, 3, templateFile.LastUpdatedAt)
	md.Writeln("{: .table .table-striped .table-hover}")
	md.WriteTable(table)
	md.Write("\n")
	md.Writeln(templateFile.Description)
	md.WriteHeader("Templates", 2)
	md.Write("\n")
	md.Writeln("A list of all the templates present in the Workflow Template")
	md.Write("\n")

	table = markdown.NewTable(len(templateFile.Templates), 3)
	table.SetTableTitle(0, "Name")
	table.SetTableTitle(1, "Type")
	table.SetTableTitle(2, "Description")

	for i, template := range templateFile.Templates {
		table.SetTableContent(i, 0, markdown.GetLink(template.Name, "#"+template.Name))
		table.SetTableContent(i, 1, markdown.GetMonospaceCode(templateTypes[template.Type]))
		table.SetTableContent(i, 2, strconv.Itoa(template.LineNumber))
	}
	md.Writeln("{: .table .table-striped .table-hover}")
	md.WriteTable(table)

	for _, template := range templateFile.Templates {
		md.Writeln("---")
		md.Write("\n")
		md.WriteHeader(template.Name, 3)
		md.Write("\n")
		md.Writeln("Type: " + markdown.GetMonospaceCode(templateTypes[template.Type]))
		md.Write("\n")
		md.Writeln("Description: " + template.Description + "\n")

		if template.Type == workflow.CONTAINER_TEMPLATE || template.Type == workflow.SCRIPT_TEMPLATE {
			md.Writeln("\nImage: " + markdown.GetMonospaceCode(template.ContainerImageTag) + "\n")
		}

		// inputs
		var inputs, inputParamList, inputArtifactList markdown.ListNode

		if template.Inputs != nil {
			inputs.Value = "Inputs"
			inputs.NodeType = markdown.ListTypeUnordered

			if len(template.Inputs.Parameters) > 0 {
				inputParamList.Value = "Parameters"
				inputParamList.NodeType = markdown.ListTypeUnordered
			}

			for _, param := range template.Inputs.Parameters {
				var inputParamChild markdown.ListNode
				trimmedDescription := strings.Trim(param.Description, "\n")
				if trimmedDescription != "" {
					if param.Required {
						inputParamChild.Value = markdown.GetMonospaceCode(param.Name) + " - (required) " + trimmedDescription
					} else {
						inputParamChild.Value = markdown.GetMonospaceCode(param.Name) + " - " + trimmedDescription
					}
				} else {
					if param.Required {
						inputParamChild.Value = markdown.GetMonospaceCode(param.Name) + " (required)"
					} else {
						inputParamChild.Value = markdown.GetMonospaceCode(param.Name)
					}
				}
				inputParamChild.NodeType = markdown.ListTypeUnordered
				inputParamList.Children = append(inputParamList.Children, &inputParamChild)
			}

			if len(template.Inputs.Artifacts) > 0 {
				inputArtifactList.Value = "Artifacts"
				inputArtifactList.NodeType = markdown.ListTypeUnordered
			}

			for _, artifact := range template.Inputs.Artifacts {
				var inputArtifactChild markdown.ListNode
				trimmedDescription := strings.Trim(artifact.Description, "\n")
				if trimmedDescription != "" {
					if artifact.Required {
						inputArtifactChild.Value = markdown.GetMonospaceCode(artifact.Name) + " - (required) " + trimmedDescription
					} else {
						inputArtifactChild.Value = markdown.GetMonospaceCode(artifact.Name) + " - " + trimmedDescription
					}
				} else {
					if artifact.Required {
						inputArtifactChild.Value = markdown.GetMonospaceCode(artifact.Name) + " (required)"
					} else {
						inputArtifactChild.Value = markdown.GetMonospaceCode(artifact.Name)
					}
				}
				inputArtifactChild.NodeType = markdown.ListTypeUnordered
				inputArtifactList.Children = append(inputArtifactList.Children, &inputArtifactChild)
			}
		}

		inputs.Children = append(inputs.Children, &inputParamList)
		inputs.Children = append(inputs.Children, &inputArtifactList)

		// outputs
		var outputs, outputParamList, outputArtifactList markdown.ListNode

		if template.Outputs != nil {
			outputs.Value = "Outputs"
			outputs.NodeType = markdown.ListTypeUnordered

			if len(template.Outputs.Parameters) > 0 {
				outputParamList.Value = "Parameters"
				outputParamList.NodeType = markdown.ListTypeUnordered
			}

			for _, param := range template.Outputs.Parameters {
				var outputParamChild markdown.ListNode
				trimmedDescription := strings.Trim(param.Description, "\n")
				if trimmedDescription != "" {
					outputParamChild.Value = markdown.GetMonospaceCode(param.Name) + " - " + trimmedDescription
				} else {
					outputParamChild.Value = markdown.GetMonospaceCode(param.Name)
				}
				outputParamChild.NodeType = markdown.ListTypeUnordered
				outputParamList.Children = append(outputParamList.Children, &outputParamChild)
			}

			if len(template.Outputs.Artifacts) > 0 {
				outputArtifactList.Value = "Artifacts"
				outputArtifactList.NodeType = markdown.ListTypeUnordered
			}

			for _, artifact := range template.Outputs.Artifacts {
				var outputArtifactChild markdown.ListNode
				trimmedDescription := strings.Trim(artifact.Description, "\n")
				if trimmedDescription != "" {
					outputArtifactChild.Value = markdown.GetMonospaceCode(artifact.Name) + " - " + trimmedDescription
				} else {
					outputArtifactChild.Value = markdown.GetMonospaceCode(artifact.Name)
				}
				outputArtifactChild.NodeType = markdown.ListTypeUnordered
				outputArtifactList.Children = append(outputArtifactList.Children, &outputArtifactChild)
			}
		}

		outputs.Children = append(outputs.Children, &outputParamList)
		outputs.Children = append(outputs.Children, &outputArtifactList)

		var tasks markdown.ListNode

		if len(template.Tasks) > 0 {
			tasks.Value = "Tasks"
			tasks.NodeType = markdown.ListTypeUnordered
		}

		for _, task := range template.Tasks {
			var name, description, taskTemplate markdown.ListNode
			name.Value = markdown.GetMonospaceCode(task.Name)
			name.NodeType = markdown.ListTypeUnordered

			description.Value = strings.Trim(task.Description, "\n")
			description.NodeType = markdown.ListTypeUnordered
			name.Children = append(name.Children, &description)

			if strings.Contains(task.Template, "::") {
				taskTemplate.Value = "Template: " + task.Template
			} else {
				taskTemplate.Value = "Template: " + markdown.GetLink(task.Template, "#"+task.Template)
			}
			taskTemplate.NodeType = markdown.ListTypeUnordered
			name.Children = append(name.Children, &taskTemplate)

			tasks.Children = append(tasks.Children, &name)
		}

		var parent markdown.ListNode

		parent.Value = ""
		parent.NodeType = markdown.ListTypeUnordered
		parent.Children = append(parent.Children, &inputs)
		parent.Children = append(parent.Children, &outputs)
		parent.Children = append(parent.Children, &tasks)

		md.WriteList(&parent)
	}
	return md, nil
}
