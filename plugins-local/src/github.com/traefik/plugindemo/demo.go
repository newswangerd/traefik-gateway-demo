// Package plugindemo a demo plugin.
package plugindemo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"io/ioutil"
	"encoding/json"


)

type TokenResponse struct {
	GatewayToken string `json:"gateway_token"`
}

// Config the plugin configuration.
type Config struct {
	GatewayAPI string `json:"gatewayAPI,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		GatewayAPI: "",
	}
}

// Demo a Demo plugin.
type Demo struct {
	next        http.Handler
	gatewayPath string
	name        string
	template    *template.Template
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &Demo{
		gatewayPath: config.GatewayAPI,
		next:        next,
		name:        name,
		template:    template.New("demo").Delims("[[", "]]"),
	}, nil
}

func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := "http://localhost:7080/api/gateway-token"

	client := http.Client{}

	auth_req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	res, getErr := client.Do(auth_req)
	if getErr != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	token := TokenResponse{}
	jsonErr := json.Unmarshal(body, &token)
	if jsonErr != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", token.GatewayToken)
	os.Stdout.WriteString(token.GatewayToken)
	// req.Header.Set("Authorization", "Basic YWRtaW46YWRtaW4=")
	os.Stdout.WriteString("________-----------_________---------")
	// os.Stdout.WriteString(a.gatewayPath)

	a.next.ServeHTTP(rw, req)
}
