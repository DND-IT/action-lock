package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dnd-it/action-lock/internal/inputs"
	"github.com/dnd-it/action-lock/internal/lock"
	"github.com/dnd-it/action-lock/internal/outputs"
)

func main() {
	cfg, err := inputs.Parse()
	if err != nil {
		outputs.Error(err.Error())
		os.Exit(1)
	}

	client := lock.New(cfg.Repository, cfg.Token)
	lockRef := fmt.Sprintf("refs/locks/%s", cfg.LockName)

	switch cfg.Action {
	case "acquire":
		acquired := acquire(client, cfg)
		outputs.Set("acquired", fmt.Sprintf("%t", acquired))
		outputs.Set("lock_ref", lockRef)
		if !acquired {
			outputs.Error(fmt.Sprintf("Failed to acquire lock %q within %ds", cfg.LockName, cfg.Timeout))
			os.Exit(1)
		}
	case "release":
		release(client, cfg)
		outputs.Set("acquired", "false")
		outputs.Set("lock_ref", lockRef)
	}
}

func acquire(client *lock.Client, cfg *inputs.Config) bool {
	deadline := time.Now().Add(time.Duration(cfg.Timeout) * time.Second)
	interval := time.Duration(cfg.PollInterval) * time.Second

	for {
		acquired, err := client.Acquire(cfg.LockName, cfg.SHA)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: lock attempt failed: %v\n", err)
		}
		if acquired {
			fmt.Printf("Lock %q acquired\n", cfg.LockName)
			return true
		}

		// Check for stale lock
		age, err := client.LockAge(cfg.LockName)
		if err == nil && age > cfg.StaleThreshold {
			fmt.Printf("Stale lock detected (%ds old, threshold %ds), removing...\n", age, cfg.StaleThreshold)
			if err := client.Release(cfg.LockName); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to remove stale lock: %v\n", err)
			}
			continue
		}

		if time.Now().After(deadline) {
			return false
		}

		remaining := time.Until(deadline).Seconds()
		fmt.Printf("Lock %q held by another process, retrying in %ds... (%.0fs remaining)\n", cfg.LockName, cfg.PollInterval, remaining)
		time.Sleep(interval)
	}
}

func release(client *lock.Client, cfg *inputs.Config) {
	if err := client.Release(cfg.LockName); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to release lock: %v\n", err)
		return
	}
	fmt.Printf("Lock %q released\n", cfg.LockName)
}
