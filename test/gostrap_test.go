package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SexyBobRiK/gostrap"
)

func TestLetsGoFullLifecycle(t *testing.T) {
	var boot *gostrap.Bootstrap
	tmpFile := filepath.Join(os.TempDir(), "gostrap.yml")
	if err := os.WriteFile(tmpFile, []byte(configContentYAML), 0644); err != nil {
		t.Fatal(err)
	}
	t.Run("Initialization", func(t *testing.T) {
		var err error
		boot, err = gostrap.LetsGo(tmpFile)
		if err != nil {
			t.Fatalf("❌ Error initialization: %v", err)
		}
		if boot.Config.ConfigName != "gostrap" {
			t.Error("⚠️ Config name incorrect")
		}
	})
	t.Run("Pulse", func(t *testing.T) {
		if err := boot.Pulse(); err != nil {
			t.Fatalf("❌ Error pulse: %v", err)
		}
		t.Log("✅ Server started")
	})
	t.Run("ShutDown", func(t *testing.T) {
		time.Sleep(10 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		boot.ShutDown(ctx)
		t.Log("⚠️ Server shut down")
	})

}
