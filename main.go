package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Config holds the configuration values required for the health check.
type Config struct {
	URL                string        // Target URL for the health check.
	ExpectedStatusCode int           // Expected HTTP status code from the target.
	Timeout            time.Duration // Timeout duration for the request.
}

// HTTPClient is an interface that defines the Get method for initiating HTTP GET requests.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// setupLogger initializes a logger instance with JSON formatting for structured logging.
func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// loadConfig returns a Config instance populated with values obtained from environment variables or command-line flags.
func loadConfig() (Config, error) {
	cfg := initConfig()
	if err := loadEnvConfig(&cfg); err != nil {
		return Config{}, err
	}
	parseFlags(&cfg)
	return cfg, nil
}

// initConfig initializes a Config instance with default values.
func initConfig() Config {
	return Config{
		URL:                "http://localhost",
		ExpectedStatusCode: http.StatusOK,
		Timeout:            2 * time.Second,
	}
}

// loadEnvConfig overwrites Config fields with values obtained from environment variables, if they exist.
func loadEnvConfig(cfg *Config) error {
	if val, exists := os.LookupEnv("HC_URL"); exists {
		cfg.URL = val
	}
	if val, exists := os.LookupEnv("HC_EXPECTED_STATUS_CODE"); exists {
		if code, err := strconv.Atoi(val); err == nil {
			cfg.ExpectedStatusCode = code
		} else {
			return fmt.Errorf("invalid value for HC_EXPECTED_STATUS_CODE (%s), using default: %w", val, err)
		}
	}
	if val, exists := os.LookupEnv("HC_TIMEOUT"); exists {
		if duration, err := time.ParseDuration(val); err == nil {
			cfg.Timeout = duration
		} else {
			return fmt.Errorf("invalid value for HC_TIMEOUT (%s), using default: %w", val, err)
		}
	}
	return nil
}

// parseFlags overwrites Config fields with values obtained from command-line flags.
func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.URL, "url", cfg.URL, "Target URL for health check")
	flag.IntVar(&cfg.ExpectedStatusCode, "status", cfg.ExpectedStatusCode, "Expected HTTP status code")
	flag.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "Request timeout duration")
	flag.Parse()
}

// performHealthCheck initiates a health check against the specified URL and logs the result.
func performHealthCheck(ctx context.Context, cfg Config, client HTTPClient) error {
	target, err := url.Parse(cfg.URL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	resp, err := client.Get(target.String())
	if err != nil {
		return fmt.Errorf("unable to reach %s: %w", target, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != cfg.ExpectedStatusCode {
		return fmt.Errorf("unexpected status code received: expected %d, got %d", cfg.ExpectedStatusCode, resp.StatusCode)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Health check successful",
		slog.String("target", target.String()),
		slog.Int("statusCode", resp.StatusCode))
	return nil
}

func main() {
	logger := setupLogger()
	slog.SetDefault(logger)

	cfg, err := loadConfig()
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError, "Error loading config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	client := &http.Client{Timeout: cfg.Timeout}
	ctx := context.Background()

	if err := performHealthCheck(ctx, cfg, client); err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "Health check failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	os.Exit(0)
}
