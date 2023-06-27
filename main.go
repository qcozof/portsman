/**
 * @Description Provide a simple service to check ports are opening.
 * @Author qcozof
 * @Date 2023.06
 **/
package main

import (
	"crypto/tls"
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
var btSslCert []byte

//go:embed keys/sample.com.key
var btSslKey []byte

const (
	defaultPauseSeconds = 30
)

var err error
var (
	hostUtils    utils.HostUtils
	commandUtils utils.CommandUtils
	fileUtils    utils.FileUtils
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

	flag.StringVar(&domain, "domain", "", "--domain sample.com")
	flag.StringVar(&ports, "ports", "9090", "--ports 9090,9091")
	flag.StringVar(&webDir, "webDir", "./", "--webDir /path/to")
	flag.BoolVar(&enableSsl, "enableSsl", false, "--enableSsl true")
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

	if len(certFile) > 0 {
		enableSsl = true
		btSslCert, err = os.ReadFile(certFile)
		if err != nil {
			msg := fmt.Sprintf(`os.ReadFile [%s] . %v `, certFile, err)
			commandUtils.PauseThenExit(defaultPauseSeconds, msg)
		}

		btSslKey, err = os.ReadFile(keyFile)
		if err != nil {
			msg := fmt.Sprintf(`os.ReadFile [%s] . %v `, keyFile, err)
			commandUtils.PauseThenExit(defaultPauseSeconds, msg)
		}
	}

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.HandleFunc("/portsman", portsman)
	http.HandleFunc("/keys/", keys)

	if !fileUtils.Exists(webDir) {
		msg := fmt.Sprintf("webDir:%s does not exist", webDir)
		commandUtils.PauseThenExit(defaultPauseSeconds, msg)
	}

	appName := filepath.Base(os.Args[0])
	fmt.Printf(`Run "%s --help" to show help.`, appName)
	fmt.Printf("\nWeb directory:%s\n\n", webDir)

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

		if enableSsl {
			protocol = "https://"
			go serveTLS(port)
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

func serveTLS(port int) {
	addr := fmt.Sprintf(":%d", port)
	cert, err := tls.X509KeyPair(btSslCert, btSslKey)
	if err != nil {
		panic(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	server := http.Server{
		TLSConfig: tlsConfig,
		Addr:      addr,
	}

	fmt.Printf("run at: %s%s%s\n", protocol, ipOrDomain, addr)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		panic(err)
	}
}

func keys(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("forbidden"))
}

func portsman(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("forbidden"))
}
