package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strconv"
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
	res, _ := cliExec(1000, "vpncmd", "localhost /CLIENT /CMD VersionGet");

	if (strings.Contains(res, "Error occurred")) {
		return "Offline"
	}

	return "Running"
}

func (a *App) ConnectVPN(host string, username string, password string) {
	res, err := cliExec(1000, "vpncmd", "localhost /CLIENT /CMD NicCreate VPN69");
	fmt.Print(res, err)
	res, err = cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountCreate blockguard /SERVER:localhost:433 /HUB:DEFAULT /USERNAME:admin /NICNAME:VPN69")
	fmt.Print(res, err)
	res, err = cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountPasswordSet blockguard /TYPE:\"standard\" /PASSWORD:\"" + password + "\"")
	fmt.Print(res, err)
	res, err = cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountUsernameSet blockguard /USERNAME:" + username)
	fmt.Print(res, err)
	res, err = cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountSet blockguard /HUB:DEFAULT /SERVER:\"" + host + "\"")
	fmt.Print(res, err)
	res, err = cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountConnect blockguard")
	fmt.Print(res, err)
}

func (a *App) DisconnectVPN() {
	res, err := cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountDisconnect blockguard")
	fmt.Print(res, err)
}

type Property struct {
	key string
	value string
}	

type VPNConnectionStatus struct {
	status string
	incomingBytes int
	outgoingBytes int
}

func (a *App) GetConnectionStatus() VPNConnectionStatus {
	res, err := cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountStatusGet blockguard")
	fmt.Print(res, err)
	
	if (strings.Contains(res, "Error code: 37")) {
		return VPNConnectionStatus{
			status: "Offline",
			incomingBytes: 0,
			outgoingBytes: 0,
		}
	}

	lines := strings.Split(res, "\n")
	properties := []Property{}

	for _, line := range lines {
		parts := strings.Split(line, "|")

		if (len(parts) != 2) {
			continue
		}

		key := strings.Trim(parts[0], "")
		value := strings.Trim(parts[1], "")

		if (key == "Item") {
			continue
		}

		properties = append(properties, Property{
			key: key,
			value: value,
		})
	}

	status, err := getProperty(properties, "Session Status")
	outgoing_raw, err := getProperty(properties, "Outgoing Data Size");
	incoming_raw, err := getProperty(properties, "Incoming Data Size")

	outgoing, err := strconv.ParseInt(strings.ReplaceAll(strings.Split(outgoing_raw, " ")[0], ",", ""), 10, 32)
	incoming, err := strconv.ParseInt(strings.ReplaceAll(strings.Split(incoming_raw, " ")[0], ",", ""), 10, 32)

	return VPNConnectionStatus{
		status: status,
		outgoingBytes: int(outgoing),
		incomingBytes: int(incoming),
	}
}

func converStatus(status string) string {
	if strings.Contains(status, "not connected") {
		return "Failed"
	}
	if (strings.Contains(status, "Connection to VPN Server Started") || strings.Contains(status, "Retrying")) {
		return "Connecting"
	}

	return "Failed"
}

func getProperty(properties []Property, key string) (string, error) {
	for _, prop := range properties {
		if (prop.key == key){
			return prop.value, nil;
		}
	}

	return "", fmt.Errorf("Key not found")
}