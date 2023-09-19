package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}


func (a *App) MTLSFetch(method string, path string, body string, csr string, privateKey string) any {
	cert, err := tls.X509KeyPair([]byte(csr), []byte(privateKey))

	if err != nil {
		return "Fail 1"
	}

	certs := []tls.Certificate{cert}

	tr := &http.Transport{
        TLSClientConfig: &tls.Config{Certificates: certs,
        InsecureSkipVerify: true},
    }

	client := &http.Client{Transport: tr}

	if err != nil {
		return "Fail 2"
	}

	req, err := http.NewRequest(method, path, bytes.NewReader([]byte(body)));

	if err != nil {
		return "Fail 3"
	}

	res, err := client.Do(req);

	if err != nil {
		return "Fail 4"
	}

	result, err := io.ReadAll(res.Body)

	if err != nil {
		return "Fail 5"
	}

	if (len(result) == 0) {
		return res.Status
	}

	return string(result[:])
}

func (a *App) SoftEtherStatus() string {
	res, _ := cliExec(1000, "vpncmd", "localhost", "/CLIENT", "/CMD NicGet", "VPN69");

	if (strings.Contains(res, "Error occurred. (Error code: 1)")) {
		return "Offline"
	}

	return "Running"
}

func (a *App) ConnectVPN(host string, username string, password string) {
	cliExec(1000, "vpncmd", "localhost", "/CLIENT", "/CMD NicGet", "VPN69");
	cliExec(1000, "vpncmd", "localhost", "/CLIENT", "/CMD AccountCreate", "blockguard", "/SERVER:localhost:433 /HUB:DEFAULT /USERNAME:admin /NICNAME:VPN69")

	cliExec(1000, "vpncmd", "localhost", "/CLIENT", "/CMD AccountPasswordSet", "blockguard", "/TYPE:\"standard\"", "/PASSWORD:" + password)
	cliExec(1000, "vpncmd", "localhost", "/CLIENT", "/CMD AccountUsernameSet", "blockguard", "/USERNAME", username)
	cliExec(1000, "vpncmd", "localhost", "/CLIENT", "/CMD AccountSet", "blockguard", "/HUB:DEFAULT", "/SERVER " + host)
}