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

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

type MTLSFetchResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

func (a *App) MTLSFetch(method string, path string, body string, csr string, privateKey string) MTLSFetchResponse {
	runtime.LogDebug(a.ctx, "MTLSFetch - "+path)
	cert, _ := tls.X509KeyPair([]byte(csr), []byte(privateKey))
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert},
			InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	req, _ := http.NewRequest(method, path, bytes.NewReader([]byte(body)))
	res, err := client.Do(req)

	if err == nil {
		defer res.Body.Close()
	} else {
		runtime.LogError(a.ctx, err.Error())
	}

	if err != nil {
		return MTLSFetchResponse{
			Success:    false,
			StatusCode: -1,
			Body:       "",
		}
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return MTLSFetchResponse{
			Success:    false,
			StatusCode: res.StatusCode,
			Body:       "",
		}
	}

	result, err := io.ReadAll(res.Body)

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return MTLSFetchResponse{
			Success:    false,
			StatusCode: res.StatusCode,
			Body:       "",
		}
	}

	return MTLSFetchResponse{
		Success:    true,
		StatusCode: res.StatusCode,
		Body:       string(result[:]),
	}
}

func (a *App) SoftEtherStatus() string {
	res, _ := cliExec(1000, "vpncmd", "localhost /CLIENT /CMD VersionGet")

	if strings.Contains(res, "Error occurred") {
		return "Offline"
	}

	return "Running"
}

func (a *App) ConnectVPN(host string, username string, password string) {
	cliExec(1000, "vpncmd", "localhost /CLIENT /CMD NicCreate VPN69")
	cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountCreate blockguard /SERVER:localhost:433 /HUB:DEFAULT /USERNAME:admin /NICNAME:VPN69")
	cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountPasswordSet blockguard /TYPE:\"standard\" /PASSWORD:\""+password+"\"")
	cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountUsernameSet blockguard /USERNAME:"+username)
	cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountSet blockguard /HUB:DEFAULT /SERVER:\""+host+"\"")
	cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountConnect blockguard")
}

func (a *App) DisconnectVPN() {
	cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountDisconnect blockguard")
}

type Property struct {
	key   string
	value string
}

type VPNConnectionStatus struct {
	Status        string `json:"status"`
	IncomingBytes int    `json:"incomingBytes"`
	OutgoingBytes int    `json:"outgoingBytes"`
}

func (a *App) GetConnectionStatus() VPNConnectionStatus {
	res, _ := cliExec(1000, "vpncmd", "localhost /CLIENT /CMD AccountStatusGet blockguard")

	if strings.Contains(res, "Error code: 37") {
		return VPNConnectionStatus{
			Status:        "Offline",
			IncomingBytes: 0,
			OutgoingBytes: 0,
		}
	}

	lines := strings.Split(res, "\n")
	properties := []Property{}

	for _, line := range lines {
		parts := strings.Split(line, "|")

		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(parts[0], " ")
		value := strings.Trim(parts[1], " ")

		if key == "Item" {
			continue
		}

		properties = append(properties, Property{
			key:   key,
			value: value,
		})
	}

	status, _ := getProperty(properties, "Session Status")
	outgoing_raw, _ := getProperty(properties, "Outgoing Data Size")
	incoming_raw, _ := getProperty(properties, "Incoming Data Size")

	outgoing, _ := strconv.ParseInt(strings.ReplaceAll(strings.Split(outgoing_raw, " ")[0], ",", ""), 10, 32)
	incoming, _ := strconv.ParseInt(strings.ReplaceAll(strings.Split(incoming_raw, " ")[0], ",", ""), 10, 32)

	return VPNConnectionStatus{
		Status:        convertStatus(status),
		OutgoingBytes: int(outgoing),
		IncomingBytes: int(incoming),
	}
}

func convertStatus(status string) string {
	if strings.Contains(status, "Connection to VPN Server Started") ||
		strings.Contains(status, "Retrying") ||
		strings.Contains(status, "Authenticating User") ||
		strings.Contains(status, "Negotiating") {
		return "Connecting"
	}
	if strings.Contains(status, "Connection Completed (Session Established)") {
		return "Connected"
	}

	return "Failed"
}

func getProperty(properties []Property, key string) (string, error) {
	for _, prop := range properties {
		if prop.key == key {
			return prop.value, nil
		}
	}

	return "", fmt.Errorf("key not found")
}
