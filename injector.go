package gaeshim

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Inject injects the environment variables and HTTP handlers (i.e. routes) into net/http's handler.
func Inject(yamlFilePath string) error {
	yamlText, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return err
	}

	conf, err := parseConfigYAML(yamlText)
	if err != nil {
		return err
	}

	if err := applyEnvVars(conf); err != nil {
		return err
	}

	interceptRequest(conf)

	return nil
}

func applyEnvVars(conf *configuration) error {
	for k, v := range conf.EnvVariables {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
		log.Printf(`[INFO] (gae-shim) set environment variable: "%s" => "%s"`, k, v)
	}
	return nil
}

func interceptRequest(conf *configuration) {
	for _, gaeHandler := range conf.Handlers {
		// NOTE static_dir is only supported
		// TODO static_files support
		// ref: https://cloud.google.com/appengine/docs/standard/python/config/appref#handlers_static_dir
		staticDir := gaeHandler.StaticDir
		if staticDir == "" {
			continue
		}

		mimeType := gaeHandler.MimeType
		httpHeaders := gaeHandler.HTTPHeaders
		handlerURL := gaeHandler.URL
		if !strings.HasSuffix(handlerURL, "/") {
			handlerURL += "/"
		}

		log.Printf("[INFO] (gae-shim) register HTTP route: %s", handlerURL)
		http.HandleFunc(handlerURL, func(w http.ResponseWriter, r *http.Request) {
			if mimeType != "" {
				w.Header().Add("Content-Type", mimeType)
			}

			for k, v := range httpHeaders {
				w.Header().Add(k, v)
			}

			http.StripPrefix(handlerURL, http.FileServer(http.Dir(staticDir))).ServeHTTP(w, r)
		})
	}
}

func parseConfigYAML(yamlText []byte) (*configuration, error) {
	var conf configuration
	err := yaml.Unmarshal(yamlText, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
