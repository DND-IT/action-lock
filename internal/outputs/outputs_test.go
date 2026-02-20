package outputs

import (
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	origStdout := os.Stdout
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = origStdout

	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	_ = r.Close()
	return string(buf[:n])
}

// --------------- Set ---------------

func TestSet_WithGitHubOutput(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "gh-output")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	_ = tmp.Close()

	t.Setenv("GITHUB_OUTPUT", tmp.Name())

	Set("acquired", "true")

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if got := string(data); got != "acquired=true\n" {
		t.Errorf("expected 'acquired=true\\n', got %q", got)
	}
}

func TestSet_WithGitHubOutput_Append(t *testing.T) {
	tmp, err := os.CreateTemp(t.TempDir(), "gh-output")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	_ = tmp.Close()

	t.Setenv("GITHUB_OUTPUT", tmp.Name())

	Set("acquired", "true")
	Set("lock_age", "60")

	data, err := os.ReadFile(tmp.Name())
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if got := string(data); got != "acquired=true\nlock_age=60\n" {
		t.Errorf("expected two lines, got %q", got)
	}
}

func TestSet_WithoutGitHubOutput(t *testing.T) {
	t.Setenv("GITHUB_OUTPUT", "")

	got := captureStdout(t, func() {
		Set("acquired", "true")
	})
	if got != "::set-output name=acquired::true\n" {
		t.Errorf("expected set-output fallback, got %q", got)
	}
}

// --------------- Notice ---------------

func TestNotice(t *testing.T) {
	got := captureStdout(t, func() {
		Notice("lock acquired")
	})
	if got != "::notice::lock acquired\n" {
		t.Errorf("expected notice, got %q", got)
	}
}

// --------------- Error ---------------

func TestError(t *testing.T) {
	got := captureStdout(t, func() {
		Error("something failed")
	})
	expected := "::error::something failed\n"
	if !strings.Contains(got, expected) {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
