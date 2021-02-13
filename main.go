package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: " + os.Args[0] + " [DNS Server:DNS Server] [Default Route] [Interface:IP Address]...")
		fmt.Println("Example: ")
		fmt.Println("  " + os.Args[0] + " 8.8.8.8 192.168.1.1 ens33:192.168.1.2")
		fmt.Println("  " + os.Args[0] + " 8.8.8.8:8.8.4.4 192.168.1.1 ens33:192.168.1.2 docker:172.17.0.1")
		exit(0)
	}

	if itemInArray("-b", os.Args) {
		args := os.Args[1:]
		for i := 0; i < len(args); i++ {
			if args[i] == "-b" {
				args[i] = ""
				break
			}
		}
		cmd := exec.Command(os.Args[0], args...)
		cmd.Start()
		os.Exit(0)
	}

	lg.setLevel("error")

	for {
		var ns1, ns2 string
		if strIn(":", os.Args[1]) {
			nss := strSplit(os.Args[1], ":")
			ns1 = nss[0]
			ns2 = nss[1]
			lg.trace("Name Server 1:", ns1)
			lg.trace("Name Server 2:", ns2)
			fd := open("/etc/resolv.conf", "w")
			fd.write("search mil" + "\n")
			fd.write("nameserver " + ns1 + "\n")
			fd.write("nameserver " + ns2 + "\n")
			fd.close()
		} else {
			ns1 = os.Args[1]
			lg.trace("Name Server:", ns1)
			fd := open("/etc/resolv.conf", "w")
			fd.write("search mil")
			fd.write("nameserver " + ns1)
			fd.close()
		}

		dfr := os.Args[2]
		lg.trace("Default route:", dfr)
		system("route add default gw " + dfr)

		for _, i := range os.Args[3:] {
			if len(strSplit(i, ":")) > 2 {
				lg.error("Too much columns:", i)
			} else if len(strSplit(i, ":")) == 2 {
				ii := strSplit(i, ":")
				interfac := ii[0]
				ip := ii[1]
				lg.trace("Interface:", interfac, ip)
				system("ifconfig " + interfac + " " + ip)
			} else {
				lg.error("Too less columns:", i)
			}
		}
		sleep(1)
	}
}
