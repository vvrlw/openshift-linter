package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gerald1248/httpscerts"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type PostStruct struct {
	Buffer string
}

func serve(certificate, key, hostname string, port int) {
	virtual_fs := &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo}

	//set up custom mux
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(virtual_fs))
	mux.HandleFunc("/openshift-linter/report", guiHandler)
	mux.HandleFunc("/openshift-linter", handler)

	// exception: no TLS req'd
	if port%2 == 0 {
		fmt.Print(listening(hostname, port, false))
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), mux))
		return
	}

	err := httpscerts.Check(certificate, key)
	if err != nil {
		cert, key, err := httpscerts.GenerateArrays(fmt.Sprintf("%s:%d", hostname, port))
		if err != nil {
			log.Fatal("Can't create https certs")
		}

		keyPair, err := tls.X509KeyPair(cert, key)
		if err != nil {
			log.Fatal("Can't create key pair")
		}

		//supply root certs in case code is run in from-scratch image
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(rootCerts)

		var certificates []tls.Certificate
		certificates = append(certificates, keyPair)

		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			PreferServerCipherSuites: true,
			Certificates:             certificates,
			RootCAs:                  pool,
		}

		s := &http.Server{
			Addr:           fmt.Sprintf("%s:%d", hostname, port),
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      cfg,
		}
		fmt.Print(listening(hostname, port, true))
		log.Fatal(s.ListenAndServeTLS("", ""))
	}
	fmt.Print(listening(hostname, port, false))
	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%d", hostname, port), certificate, key, mux))
}

func listening(hostname string, port int, selfCert bool) string {
	var selfCertMsg string
	if selfCert {
		selfCertMsg = " (self-certified)"
	}

	protocol := "https"
	if port%2 == 0 {
		protocol = "http"
	}

	// for Docker, mustn't specify localhost explicitly
	if hostname == "" {
		hostname = "0.0.0.0"
	}

	return fmt.Sprintf("Listening on port %d%s\n"+
		"POST JSON sources to %s://%s:%d/openshift-linter\n"+
		"Generate report at %s://%s:%d/openshift-linter/report\n", port, selfCertMsg, protocol, hostname, port, protocol, hostname, port)
}

func guiHandler(w http.ResponseWriter, r *http.Request) {
	//show GUI from static resources
	bytes, _ := Asset("static/index.html")
	fmt.Fprintf(w, "%s\n", string(bytes))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePost(&w, r)
	case "GET":
		handleGet(&w, r)
	}
}

func handleGet(w *http.ResponseWriter, r *http.Request) {
	//fetch configuration list
	fmt.Fprintf(*w, "GET request\nRequest struct = %v\n", r)
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	//process configuration list
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body: %s", err)
		return
	}

	any := ".*" //for POST requests, patterns are specified within config body
	combinedResultMap, err := processBytes(body, LinterParams{any, any, any, any, "", any, any, "json"})
	if err != nil {
		fmt.Fprintf(*w, "Can't process POST request body: %s\n%s\n", err, body)
		return
	}

	json, err := json.Marshal(combinedResultMap)
	if err != nil {
		fmt.Fprintf(*w, "Can't encode result: %s", err)
		return
	}

	fmt.Fprintf(*w, string(json))
}
