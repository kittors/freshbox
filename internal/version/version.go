package version

// Version is the current freshbox release version.
// This can be overridden at build time via:
//
//	go build -ldflags "-X github.com/kittors/freshbox/internal/version.Version=v2.0.0"
var Version = "v1.0.0"
