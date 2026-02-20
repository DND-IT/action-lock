package inputs

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Action         string
	LockName       string
	Timeout        int
	PollInterval   int
	StaleThreshold int
	Token          string
	Repository     string
	SHA            string
}

func Parse() (*Config, error) {
	action := os.Getenv("INPUT_ACTION")
	if action != "acquire" && action != "release" {
		return nil, fmt.Errorf("invalid action %q: must be 'acquire' or 'release'", action)
	}

	lockName := os.Getenv("INPUT_LOCK_NAME")
	if lockName == "" {
		return nil, fmt.Errorf("lock_name is required")
	}

	token := os.Getenv("INPUT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		return nil, fmt.Errorf("GITHUB_REPOSITORY not set")
	}

	sha := os.Getenv("GITHUB_SHA")
	if sha == "" {
		return nil, fmt.Errorf("GITHUB_SHA not set")
	}

	timeout := intEnv("INPUT_TIMEOUT", 300)
	pollInterval := intEnv("INPUT_POLL_INTERVAL", 10)
	staleThreshold := intEnv("INPUT_STALE_THRESHOLD", 600)

	return &Config{
		Action:         action,
		LockName:       lockName,
		Timeout:        timeout,
		PollInterval:   pollInterval,
		StaleThreshold: staleThreshold,
		Token:          token,
		Repository:     repo,
		SHA:            sha,
	}, nil
}

func intEnv(key string, defaultVal int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
