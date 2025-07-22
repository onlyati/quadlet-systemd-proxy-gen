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
	containerRoot := splitQuadlet[0]
	if splitQuadlet[1] == "pod" {
		containerRoot = fmt.Sprintf("%s-pod", splitQuadlet[0])
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	tmplData := map[string]any{
		"ContainerRoot": containerRoot,
		"QuadletIP":     data.quadletIP,
	}

	for _, port := range data.ports {
		fileName := fmt.Sprintf("%s-proxy-%d.service", containerRoot, port)
		filePath := path.Join(homeDir, ".config", "systemd", "user", fileName)
		fmt.Println("generate file: " + filePath)

		tmplData["Port"] = port

		if data.quadletPort != 0 {
			tmplData["QuadletPort"] = data.quadletPort
		} else {
			tmplData["QuadletPort"] = port
		}
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
	containerRoot := splitQuadlet[0]
	if splitQuadlet[1] == "pod" {
		containerRoot = fmt.Sprintf("%s-pod", splitQuadlet[0])
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	tmplData := map[string]any{
		"Container": data.quadlet,
		"DevIP":     data.ip,
	}

	for _, port := range data.ports {
		fileName := fmt.Sprintf("%s-proxy-%d.socket", containerRoot, port)
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
