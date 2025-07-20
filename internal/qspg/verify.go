package qspg

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

// Read the container from the `~/.config/containers/systemd/{{name}}` file
// and looking for `PublishPort=` lines. parse the line and get the exposed port
// if it is exposed via 127.0.0.1 interface.
func verifyContainer(name string) ([]uint16, error) {
	if name == "" {
		return nil, errors.New("container not defined")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filePath := path.Join(homeDir, ".config", "containers", "systemd", name)
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return nil, errors.New("file does not exists: " + filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var ports []uint16
	for scanner.Scan() {
		if s, found := strings.CutPrefix(scanner.Text(), "PublishPort="); found {
			tmp := strings.Split(s, ":")
			if len(tmp) >= 2 {
				if tmp[0] != "127.0.0.1" {
					fmt.Println("WARNING: port not exposed on 127.0.0.1, ignored: " + scanner.Text())
					continue
				}

				p, err := strconv.ParseUint(tmp[1], 10, 16)
				if err != nil {
					return nil, err
				}
				ports = append(ports, uint16(p))
			}
		}
	}

	return ports, nil
}

// Verify that IP address is a valid IP address then check if this address
// belongs to the loopback device of the system.
func verifyIP(ipAddress string) (bool, error) {
	targetIP := net.ParseIP(ipAddress)
	if targetIP == nil {
		return false, errors.New("invalid ip address " + ipAddress)
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}

	for _, iface := range interfaces {
		if (iface.Flags & net.FlagLoopback) == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return false, err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.Equal(targetIP) {
				return true, nil
			}
		}
	}

	return false, nil
}
