/**
 * @Description Provide a simple service to check ports are opening.
 * @Author qcozof
 * @Date 2023.06
 **/
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/qcozof/portsman/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed build/versionData.txt
var versionData string

//go:embed keys/fullchain.cer
var sslCertStr []byte

//go:embed keys/sample.com.key
var sslKeyStr []byte

const (
	defaultPauseSeconds = 30
	defaultCertFilePath = "keys/fullchain.cer"
	defaultKeyFilePath  = "keys/sample.com.key"
)

var err error
var (
	hostUtils      utils.HostUtils
	commandUtils   utils.CommandUtils
	fileUtils      utils.FileUtils
	directoryUtils utils.DirectoryUtils
)

var (
	ipOrDomain string
	protocol   string
)

func main() {
	var ch = make(chan bool)
	var (
		domain, ports, webDir string
		enableSsl             bool
		certFile, keyFile     string
	)
	fmt.Printf(versionData)

	checkAndGenerateSlFile()
	flag.StringVar(&domain, "domain", "", "--domain sample.com")
	flag.StringVar(&ports, "ports", "9090", "--ports 9090,9091")
	flag.StringVar(&webDir, "webDir", "./", "--webDir /path/to")
	flag.StringVar(&certFile, "certFile", "", "--certFile /path/to/fullchain.cer")
	flag.StringVar(&keyFile, "keyFile", "", "--keyFile /path/to/sample.com.key")
	flag.Parse()

	if domain == "" {
		ipOrDomain, err = hostUtils.GetHostIP()
		if err != nil {
			fmt.Println("hostUtils.GetHostIP:", err)
		}
	} else {
		ipOrDomain = domain
	}

	if certFile != "" {
		enableSsl = true
	}

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/portsman", portsman)
	http.HandleFunc("/keys/", keys)

	if !fileUtils.Exists(webDir) {
		msg := fmt.Sprintf("%s does not exist", webDir)
		commandUtils.PauseThenExit(defaultPauseSeconds, msg)
	}

	appName := filepath.Base(os.Args[0])
	for _, p := range strings.Split(ports, ",") {
		port, err := strconv.Atoi(p)
		if err != nil {
			msg := fmt.Sprintf(`Port [%s] is not an integer. `, p)
			commandUtils.PauseThenExit(defaultPauseSeconds, msg)
		}

		if port > 65535 {
			msg := fmt.Sprintf("Ports must between 1 and 65535")
			commandUtils.PauseThenExit(defaultPauseSeconds, msg)
		}

		fmt.Printf(`Run "%s --help" to show help.`, appName)
		fmt.Printf("\nWeb directory:%s\n\n", webDir)

		if enableSsl {
			protocol = "https://"
			go serveTLS(port, certFile, keyFile)
		} else {
			protocol = "http://"
			go serve(port)
		}
	}

	<-ch
}

func serve(port int) {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("run at: %s%s%s\n", protocol, ipOrDomain, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func serveTLS(port int, certFile, keyFile string) {
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("run at: %s%s%s\n", protocol, ipOrDomain, addr)
	if err := http.ListenAndServeTLS(addr, certFile, keyFile, nil); err != nil {
		panic(err)
	}
}

func keys(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("forbidden"))
}

func portsman(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("forbidden"))
}

func checkAndGenerateSlFile() {
	const keysDir = "keys/"
	const perm = 0600
	if exists, _ := directoryUtils.PathExists(keysDir); !exists {
		if err := directoryUtils.CreateDir(keysDir); err != nil {
			fmt.Println("CreateDir", keysDir, err)
			return
		}

		if err := fileUtils.Write(defaultCertFilePath, sslCertStr, perm); err != nil {
			fmt.Println("Write", defaultCertFilePath, err)
			return
		}

		if err := fileUtils.Write(defaultKeyFilePath, sslKeyStr, perm); err != nil {
			fmt.Println("Write", defaultKeyFilePath, err)
			return
		}
	}
}
