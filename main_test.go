package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestInitConfig(t *testing.T) {
	expected := Config{
		URL:                "http://localhost",
		ExpectedStatusCode: 200,
		Timeout:            2 * time.Second,
	}
	result := initConfig()
	if expected != result {
		t.Errorf("Expected %+v but got %+v", expected, result)
	}
}

func TestLoadEnvConfig(t *testing.T) {
	os.Setenv("HC_URL", "http://example.com")
	defer os.Unsetenv("HC_URL")

	os.Setenv("HC_EXPECTED_STATUS_CODE", "201")
	defer os.Unsetenv("HC_EXPECTED_STATUS_CODE")

	os.Setenv("HC_TIMEOUT", "3s")
	defer os.Unsetenv("HC_TIMEOUT")

	expected := Config{
		URL:                "http://example.com",
		ExpectedStatusCode: 201,
		Timeout:            3 * time.Second,
	}

	config := initConfig()
	err := loadEnvConfig(&config)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if expected != config {
		t.Errorf("Expected %+v but got %+v", expected, config)
	}

	// Test with an invalid expected status code
	os.Setenv("HC_EXPECTED_STATUS_CODE", "invalid")
	err = loadEnvConfig(&config)
	if err == nil {
		t.Error("Expected error due to invalid HC_EXPECTED_STATUS_CODE, but got nil")
	}

	// Test with an invalid timeout
	os.Setenv("HC_TIMEOUT", "invalid")
	err = loadEnvConfig(&config)
	if err == nil {
		t.Error("Expected error due to invalid HC_TIMEOUT, but got nil")
	}
}

func TestPerformHealthCheck(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer testServer.Close()

	cfg := Config{
		URL:                testServer.URL,
		ExpectedStatusCode: 200,
		Timeout:            2 * time.Second,
	}
	err := performHealthCheck(context.Background(), cfg, testServer.Client())
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Case: Unexpected status code
	cfg.ExpectedStatusCode = 201
	err = performHealthCheck(context.Background(), cfg, testServer.Client())
	if err == nil {
		t.Error("Expected error due to unexpected status code, but got nil")
	}

	// Case: Invalid URL
	cfg.URL = "invalid-url"
	err = performHealthCheck(context.Background(), cfg, testServer.Client())
	if err == nil {
		t.Error("Expected error due to invalid URL, but got nil")
	}
}

func TestLoadConfig(t *testing.T) {
	os.Setenv("HC_URL", "http://test.com")
	defer os.Unsetenv("HC_URL")

	os.Setenv("HC_EXPECTED_STATUS_CODE", "201")
	defer os.Unsetenv("HC_EXPECTED_STATUS_CODE")

	os.Setenv("HC_TIMEOUT", "3s")
	defer os.Unsetenv("HC_TIMEOUT")

	expected := Config{
		URL:                "http://test.com",
		ExpectedStatusCode: 201,
		Timeout:            3 * time.Second,
	}
	result, err := loadConfig()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if expected != result {
		t.Errorf("Expected %+v but got %+v", expected, result)
	}
}
