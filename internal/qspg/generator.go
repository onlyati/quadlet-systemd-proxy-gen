package qspg

import (
	"embed"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

// Generate {{container}}-proxy-{{port}}.service systemd unit in the user directory
func generateServiceFile(templates embed.FS, data tmplData) error {
	tmpl, err := template.ParseFS(templates, "templates/service.tmpl")
	if err != nil {
		return err
	}

	splitQuadlet := strings.Split(data.quadlet, ".")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	tmplData := map[string]any{
		"ContainerRoot": splitQuadlet[0],
	}

	for _, port := range data.ports {
		fileName := fmt.Sprintf("%s-proxy-%d.service", splitQuadlet[0], port)
		filePath := path.Join(homeDir, ".config", "systemd", "user", fileName)
		fmt.Println("generate file: " + filePath)

		tmplData["Port"] = port
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = tmpl.Execute(file, tmplData)
		if err != nil {
			return err
		}
	}

	return nil
}

// Generate {{container}}-proxy-{{port}}.socket systemd unit in the user directory
func generateSocketFile(templates embed.FS, data tmplData) error {
	tmpl, err := template.ParseFS(templates, "templates/socket.tmpl")
	if err != nil {
		return err
	}

	splitQuadlet := strings.Split(data.quadlet, ".")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	tmplData := map[string]any{
		"Container": data.quadlet,
		"DevIP":     data.ip,
	}

	for _, port := range data.ports {
		fileName := fmt.Sprintf("%s-proxy-%d.socket", splitQuadlet[0], port)
		filePath := path.Join(homeDir, ".config", "systemd", "user", fileName)
		fmt.Println("generate file: " + filePath)

		tmplData["Port"] = port
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = tmpl.Execute(file, tmplData)
		if err != nil {
			return err
		}
	}

	return nil
}
