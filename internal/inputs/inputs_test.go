package inputs

import (
	"testing"
)

func setRequiredEnv(t *testing.T) {
	t.Setenv("INPUT_ACTION", "acquire")
	t.Setenv("INPUT_LOCK_NAME", "deploy")
	t.Setenv("INPUT_TOKEN", "ghp_test")
	t.Setenv("GITHUB_REPOSITORY", "owner/repo")
	t.Setenv("GITHUB_SHA", "abc123")
}

func TestParse_ValidAcquire(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("INPUT_TIMEOUT", "60")
	t.Setenv("INPUT_POLL_INTERVAL", "5")
	t.Setenv("INPUT_STALE_THRESHOLD", "120")
	t.Setenv("INPUT_FAIL_ON_TIMEOUT", "false")

	cfg, err := Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Action != "acquire" {
		t.Errorf("expected acquire, got %s", cfg.Action)
	}
	if cfg.LockName != "deploy" {
		t.Errorf("expected deploy, got %s", cfg.LockName)
	}
	if cfg.Token != "ghp_test" {
		t.Errorf("expected ghp_test, got %s", cfg.Token)
	}
	if cfg.Repository != "owner/repo" {
		t.Errorf("expected owner/repo, got %s", cfg.Repository)
	}
	if cfg.SHA != "abc123" {
		t.Errorf("expected abc123, got %s", cfg.SHA)
	}
	if cfg.Timeout != 60 {
		t.Errorf("expected 60, got %d", cfg.Timeout)
	}
	if cfg.PollInterval != 5 {
		t.Errorf("expected 5, got %d", cfg.PollInterval)
	}
	if cfg.StaleThreshold != 120 {
		t.Errorf("expected 120, got %d", cfg.StaleThreshold)
	}
	if cfg.FailOnTimeout != false {
		t.Errorf("expected false, got %v", cfg.FailOnTimeout)
	}
}

func TestParse_ValidRelease_Defaults(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("INPUT_ACTION", "release")

	cfg, err := Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Action != "release" {
		t.Errorf("expected release, got %s", cfg.Action)
	}
	if cfg.Timeout != 300 {
		t.Errorf("expected default 300, got %d", cfg.Timeout)
	}
	if cfg.PollInterval != 10 {
		t.Errorf("expected default 10, got %d", cfg.PollInterval)
	}
	if cfg.StaleThreshold != 600 {
		t.Errorf("expected default 600, got %d", cfg.StaleThreshold)
	}
	if cfg.FailOnTimeout != true {
		t.Errorf("expected default true, got %v", cfg.FailOnTimeout)
	}
}

func TestParse_MissingAction(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("INPUT_ACTION", "")

	_, err := Parse()
	if err == nil {
		t.Fatal("expected error for missing action")
	}
}

func TestParse_InvalidAction(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("INPUT_ACTION", "invalid")

	_, err := Parse()
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}

func TestParse_MissingLockName(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("INPUT_LOCK_NAME", "")

	_, err := Parse()
	if err == nil {
		t.Fatal("expected error for missing lock_name")
	}
}

func TestParse_MissingToken(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("INPUT_TOKEN", "")

	_, err := Parse()
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}

func TestParse_MissingRepo(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("GITHUB_REPOSITORY", "")

	_, err := Parse()
	if err == nil {
		t.Fatal("expected error for missing GITHUB_REPOSITORY")
	}
}

func TestParse_MissingSHA(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("GITHUB_SHA", "")

	_, err := Parse()
	if err == nil {
		t.Fatal("expected error for missing GITHUB_SHA")
	}
}

// --------------- intEnv ---------------

func TestIntEnv_Set(t *testing.T) {
	t.Setenv("TEST_INT", "42")
	if got := intEnv("TEST_INT", 0); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestIntEnv_Empty(t *testing.T) {
	t.Setenv("TEST_INT", "")
	if got := intEnv("TEST_INT", 99); got != 99 {
		t.Errorf("expected default 99, got %d", got)
	}
}

func TestIntEnv_Invalid(t *testing.T) {
	t.Setenv("TEST_INT", "notanumber")
	if got := intEnv("TEST_INT", 99); got != 99 {
		t.Errorf("expected default 99, got %d", got)
	}
}

// --------------- boolEnv ---------------

func TestBoolEnv_True(t *testing.T) {
	t.Setenv("TEST_BOOL", "true")
	if got := boolEnv("TEST_BOOL", false); got != true {
		t.Errorf("expected true, got %v", got)
	}
}

func TestBoolEnv_False(t *testing.T) {
	t.Setenv("TEST_BOOL", "false")
	if got := boolEnv("TEST_BOOL", true); got != false {
		t.Errorf("expected false, got %v", got)
	}
}

func TestBoolEnv_Empty(t *testing.T) {
	t.Setenv("TEST_BOOL", "")
	if got := boolEnv("TEST_BOOL", true); got != true {
		t.Errorf("expected default true, got %v", got)
	}
}

func TestBoolEnv_Invalid(t *testing.T) {
	t.Setenv("TEST_BOOL", "notabool")
	if got := boolEnv("TEST_BOOL", true); got != true {
		t.Errorf("expected default true, got %v", got)
	}
}
