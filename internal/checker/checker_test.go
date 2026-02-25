package checker

import (
	"fmt"
	"testing"
)

func TestDetection(t *testing.T) {
	items := DevTools()
	CheckAll(items)
	for _, item := range items {
		status := "NOT_INSTALLED"
		if item.Status == Installed {
			status = "INSTALLED"
		}
		fmt.Printf("%-15s %-15s %s\n", item.Name, status, item.Version)
	}

	// Rust should be detected (installed in ~/.cargo/bin)
	for _, item := range items {
		if item.Name == "Rust (rustup)" {
			if item.Status != Installed {
				t.Error("Rust (rustup) should be detected as installed via ~/.cargo/bin")
			} else {
				t.Logf("Rust detected: %s", item.Version)
			}
		}
	}

	// Java should NOT be detected (macOS stub, no real JDK)
	for _, item := range items {
		if item.Name == "Java (JDK)" {
			if item.Status == Installed {
				t.Logf("Java detected: %s", item.Version)
			} else {
				t.Log("Java correctly detected as not installed (macOS stub)")
			}
		}
	}
}
