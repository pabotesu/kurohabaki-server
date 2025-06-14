package etcd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pabotesu/kurohabaki-server/config"
	"github.com/pkg/errors"
)

var etcdCmd *exec.Cmd

func StartEtcd(cfg config.EtcdConfig) error {
	args := []string{
		fmt.Sprintf("--listen-client-urls=http://%s:%d,http://127.0.0.1:%d",
			cfg.AdvertiseClientIP, cfg.AdvertiseClientPort, cfg.ListenClientPort),
		fmt.Sprintf("--advertise-client-urls=http://%s:%d", cfg.AdvertiseClientIP, cfg.AdvertiseClientPort),
	}

	etcdCmd = exec.Command("etcd", args...)
	etcdCmd.Stdout = newPrefixedLogger("etcd", "etcd_bootstrap.log")
	etcdCmd.Stderr = etcdCmd.Stdout

	err := etcdCmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start etcd")
	}

	// Wait until etcd is ready to accept client requests
	endpoint := fmt.Sprintf("%s:%d", cfg.AdvertiseClientIP, cfg.AdvertiseClientPort)
	if err := waitForEtcdReady(endpoint, 10*time.Second); err != nil {
		return errors.Wrap(err, "etcd did not become ready in time")
	}

	return nil
}

func StopEtcd() error {
	if etcdCmd == nil || etcdCmd.Process == nil {
		return nil
	}
	log.Println("🛑 Stopping embedded etcd...")
	if err := etcdCmd.Process.Kill(); err != nil {
		return errors.Wrap(err, "failed to stop etcd")
	}
	_, err := etcdCmd.Process.Wait()
	return err
}

func waitForEtcdReady(endpoint string, timeout time.Duration) error {
	url := fmt.Sprintf("http://%s/health", endpoint)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for etcd to be ready: %w", ctx.Err())
		case <-ticker.C:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				return nil
			}
		}
	}
}

func newPrefixedLogger(prefix string, logFile string) *logWriter {
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("⚠ Failed to open log file %s: %v", logFile, err)
		return &logWriter{prefix: prefix, w: os.Stdout}
	}
	return &logWriter{prefix: prefix, w: f}
}

type logWriter struct {
	prefix string
	w      *os.File
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		_, err := fmt.Fprintf(l.w, "%s %s\n", l.prefix, line)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
