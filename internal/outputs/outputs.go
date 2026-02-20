package outputs

import (
	"fmt"
	"os"
)

func Set(key, value string) {
	path := os.Getenv("GITHUB_OUTPUT")
	if path == "" {
		fmt.Printf("::set-output name=%s::%s\n", key, value)
		return
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "::error::Failed to open GITHUB_OUTPUT: %v\n", err)
		return
	}
	defer func() { _ = f.Close() }()

	_, _ = fmt.Fprintf(f, "%s=%s\n", key, value)
}

func Notice(msg string) {
	fmt.Printf("::notice::%s\n", msg)
}

func Error(msg string) {
	fmt.Printf("::error::%s\n", msg)
}
