package main

import (
	"embed"

	"github.com/onlyati/quadlet-systemd-proxy-gen/internal/qspg"
)

//go:embed templates
var TemplateFS embed.FS

func main() {
	qspg.Main(TemplateFS)
}
