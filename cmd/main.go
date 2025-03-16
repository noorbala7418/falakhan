package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/noorbala7418/falakhan/internal/model"
	"github.com/noorbala7418/falakhan/pkg/config"
	"github.com/noorbala7418/falakhan/pkg/s3"
	"github.com/sirupsen/logrus"
)

func initLog(logLevel string) {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	switch logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	cfgPath, cfgErr := config.ParseFlags()
	if cfgErr != nil {
		logrus.Fatal(cfgErr)
	}
	cfg, cfgErr := config.NewConfig(cfgPath)
	if cfgErr != nil {
		logrus.Fatal(cfgErr)
	}

	initLog(cfg.LogLevel)

	logrus.Info("falakhan started.")
	logrus.Info("log level is: ", cfg.LogLevel)

	http.HandleFunc("GET /", homePage)
	http.HandleFunc("POST /upload/{dir}", uploadFile(cfg))

	logrus.Info("Start Webserver. Listen on ", cfg.Listen)
	err := http.ListenAndServe(cfg.Listen, nil)

	if errors.Is(err, http.ErrServerClosed) {
		logrus.Info("Server Shutdown. Goodbye!", time.Now())
	} else if err != nil {
		logrus.Error("Error in running webserver. ", err)
		os.Exit(2)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	appInfo := map[string]string{
		"app": "Falakhan",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appInfo)
}

func routeMapper(path string, config model.Config) (*model.S3Credential, error) {
	for i := 0; i < len(config.Routes); i++ {
		if path == config.Routes[i].Name {
			logrus.Debugf("function routeChecker: related path %s detected", path)
			result, err := s3CredentialFinder(config.Routes[i].S3Config, config.S3Credentials)

			if err != nil {
				logrus.Error("function routeMapper: Error in find s3credential. path: ", path, ". error: ", err)
				return nil, fmt.Errorf("function routeMapper: Error in find s3credential. path: %s. error: %w", path, err)
			}
			return result, nil
		}
	}
	logrus.Debug("function routeMapper: No route found with path ", path)
	return &model.S3Credential{}, nil
}

func s3CredentialFinder(configName string, config []model.S3Credential) (*model.S3Credential, error) {
	for i := 0; i < len(config); i++ {
		if configName == config[i].Name {
			logrus.Debug("function s3CredentialFinder: config found with name ", configName)
			return &config[i], nil
		}
	}
	logrus.Errorf("function s3CredentialFinder: no config found with name %s", configName)
	return nil, fmt.Errorf("function s3CredentialFinder: no config found with name %s", configName)
}

func uploadFile(config *model.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		directory := r.PathValue("dir")

		// Maximum upload of 10 MB files
		r.ParseMultipartForm(10 << 20)

		// Get handler for filename, size and headers
		file, handler, err := r.FormFile("file")
		if err != nil {
			logrus.Error("function uploadFile. Error Retrieving the File: ", err)
			resultStatus := map[string]string{
				"status":       "failed",
				"error-reason": err.Error(),
			}
			result, _ := json.Marshal(resultStatus)

			http.Error(w, string(result), http.StatusInternalServerError)
			return
		}

		defer file.Close()
		logrus.Debug("function uploadFile. File name: ", handler.Filename)
		logrus.Debug("function uploadFile. File Size: ", handler.Size)
		logrus.Debug("function uploadFile. MIME Header: ", handler.Header)

		s3Credentials, s3MapErr := routeMapper(directory, *config)

		if (s3MapErr != nil) || (s3Credentials == nil) {
			logrus.Error("function uploadFile. Error in finding related config: ", s3MapErr)
			resultStatus := map[string]string{
				"status":       "failed",
				"error-reason": "Current route is not related",
			}
			result, _ := json.Marshal(resultStatus)

			http.Error(w, string(result), http.StatusNotAcceptable)
			return
		}

		s3Location, s3Err := s3.ObjectS3Store(*s3Credentials, directory+"/"+handler.Filename, file)

		if s3Err != nil {
			logrus.Error("function uploadFile. Error in pushing to s3: ", s3Err)
			resultStatus := map[string]string{
				"status":       "failed",
				"error-reason": s3Err.Error(),
			}
			result, _ := json.Marshal(resultStatus)

			http.Error(w, string(result), http.StatusInternalServerError)
			return
		}

		resultStatus := map[string]string{
			"status":   "done",
			"location": s3Location,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resultStatus)
	}
}
