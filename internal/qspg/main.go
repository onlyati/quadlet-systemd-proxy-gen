package qspg

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"strings"
)

type tmplData struct {
	ip          string
	port        uint
	quadlet     string
	ports       []uint
	quadletIP   string
	quadletPort uint
}

func Main(templates embed.FS) {
	ip := flag.String("ip", "10.0.0.1", "IP address where socket bind")
	port := flag.Uint("port", 0, "Port for socket address, default: same as in quadlet file")
	quadlet := flag.String("quadlet", "", "Name of the the *.contianer or *.pod file that is read and parsed for port")
	quadletIP := flag.String("quadlet-ip", "127.0.0.1", "IP address where socket bind")
	quadletPort := flag.Uint("quadlet-port", 0, "Port for socket file, if not defined then automatically discover")
	flag.Parse()

	fmt.Println("verify parameters:")
	fmt.Println("- ip: " + *ip)
	fmt.Printf("- port: %d\n", *port)
	fmt.Println("- container: " + *quadlet)
	fmt.Println("- quadletIP: " + *quadletIP)
	fmt.Printf("- quadletPort: %d\n", *quadletPort)

	// Validate the socket IP address
	found, err := verifyIP(*ip)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	if !found {
		fmt.Println("ERROR: defined ip address is not found: " + *ip)
		os.Exit(1)
	}

	// Valiudate the tareget IP address
	found, err = verifyIP(*quadletIP)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	if !found {
		fmt.Println("ERROR: defined ip address is not found: " + *quadletIP)
		os.Exit(1)
	}

	ports, err := verifyContainer(*quadlet)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}

	// If quadlet port specified, does not matter what has been discovered
	if *port != 0 {
		ports = []uint{*port}
	}

	data := tmplData{
		ip:          *ip,
		port:        *port,
		quadlet:     *quadlet,
		ports:       ports,
		quadletIP:   *quadletIP,
		quadletPort: *quadletPort,
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
