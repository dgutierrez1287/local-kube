package template

import (
	"bytes"
	"embed"
	"io/fs"
	"text/template"

	"github.com/dgutierrez1287/local-kube/logger"
)

//go:embed provision/*
var provisionTemplatesFS  embed.FS

//go:embed vagrantfiles/*
var vagrantfileTemplatesFS embed.FS

func GetProvisionFS() fs.FS {
  return provisionTemplatesFS
}

func GetVagrantFS() fs.FS {
  return vagrantfileTemplatesFS
}

func RenderProvisionTemplate(name string, data interface{}) (string, error) {
  templateName := "provision/" + name + ".tmpl"

  content, err := provisionTemplatesFS.ReadFile(templateName)
  if err != nil {
    logger.LogError("Error reading provision template", "template", templateName)
    return "", err
  }

  tmpl, err := template.New(templateName).Parse(string(content))
  if err != nil {
    logger.LogError("Error parsing provision template", "template", templateName)
    return "", err
  }

  var result bytes.Buffer
  if err := tmpl.Execute(&result, data); err != nil {
    logger.LogError("Error rendering provision template", "template", templateName)
    return "", err
  }

  return result.String(), nil
}

func RenderVagrantfileTemplate(providerType string, clusterType string, data interface{}) (string, error) {
  templateName := "vagrantfiles/" + providerType + "-" + clusterType + ".tmpl"

  content, err := vagrantfileTemplatesFS.ReadFile(templateName)
  if err != nil {
    logger.LogError("Error reading vagrantfile template", "template", templateName)
    return "", err
  }

  tmpl, err := template.New(templateName).Parse(string(content))
  if err != nil {
    logger.LogError("Error parsing vagrantfile template", "template", templateName)
    return "", err
  }

  var result bytes.Buffer
  if err := tmpl.Execute(&result, data); err != nil {
    logger.LogError("Error rendering vagrantfile template", "template", templateName)
    return "", err
  }

  return result.String(), nil
}

