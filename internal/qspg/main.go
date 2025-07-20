package qspg

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
)

type tmplData struct {
	ip      string
	quadlet string
	ports   []uint16
}

func Main(templates embed.FS) {
	ip := flag.String("ip", "10.0.0.1", "IP address where socket bind")
	quadlet := flag.String("quadlet", "", "Name of the the *.contianer or *.pod file that is read and parsed for port")
	flag.Parse()

	fmt.Println("verify parameters:")
	fmt.Println("- ip: " + *ip)
	fmt.Println("- container: " + *quadlet)

	found, err := verifyIP(*ip)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	if !found {
		fmt.Println("ERROR: defined ip address is not a loop back device")
		os.Exit(1)
	}

	ports, err := verifyContainer(*quadlet)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	data := tmplData{
		ip:      *ip,
		quadlet: *quadlet,
		ports:   ports,
	}
	fmt.Printf("creating socket and proxy files for ports: %+v\n", ports)

	err = generateSocketFile(templates, data)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	err = generateServiceFile(templates, data)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	postProcessing(data)
}

func postProcessing(data tmplData) {
	container := strings.Split(data.quadlet, ".")[0]

	fmt.Println("\nPost processing:")
	fmt.Println("================")
	fmt.Println("1. execute following commands to activate the generated data:")
	fmt.Println("   systemctl --user daemon-reload")

	fmt.Println("2. activate sockets")
	for _, port := range data.ports {
		fmt.Println("   be assume that [Unit] part contains the following in container files:")
		fmt.Printf("     %s -> BindsTo=%s-proxy-%d.service\n", data.quadlet, container, port)
		fmt.Println("     systemctl --user daemon-reload")
		fmt.Println("   execute command")
		fmt.Printf("     systemctl --user enable --now %s-proxy-%d.socket \n", container, port)
	}
}
