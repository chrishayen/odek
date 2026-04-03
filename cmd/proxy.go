package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy"
	"github.com/router-for-me/CLIProxyAPI/v6/sdk/config"
	log "github.com/sirupsen/logrus"
)

// startProxy embeds CLIProxyAPI and runs it in-process.
// When quiet is true, proxy logs are suppressed.
func startProxy(quiet bool, skipModelCheck ...bool) (func(), error) {
	configPath := os.Getenv("CLIPROXY_CONFIG")
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("load proxy config: %w", err)
	}

	if quiet {
		cfg.Debug = false
		cfg.RequestLog = false
		gin.SetMode(gin.ReleaseMode)
		log.SetOutput(io.Discard)
	}

	ready := make(chan struct{})
	svc, err := cliproxy.NewBuilder().
		WithConfig(cfg).
		WithConfigPath(configPath).
		WithHooks(cliproxy.Hooks{
			OnAfterStart: func(s *cliproxy.Service) { close(ready) },
		}).
		Build()
	if err != nil {
		return nil, fmt.Errorf("build proxy: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		if err := svc.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		}
	}()

	select {
	case <-ready:
	case err := <-errCh:
		cancel()
		return nil, fmt.Errorf("proxy: %w", err)
	case <-time.After(15 * time.Second):
		cancel()
		return nil, fmt.Errorf("proxy failed to start within 15s")
	}

	if os.Getenv("ANTHROPIC_BASE_URL") == "" {
		os.Setenv("ANTHROPIC_BASE_URL", fmt.Sprintf("http://127.0.0.1:%d", cfg.Port))
	}
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		os.Setenv("ANTHROPIC_API_KEY", "sk-local-proxy")
	}

	if len(skipModelCheck) > 0 && skipModelCheck[0] {
		return cancel, nil
	}

	// Wait for models to load.
	proxyURL := os.Getenv("ANTHROPIC_BASE_URL")
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	for i := 0; i < 100; i++ {
		req, _ := http.NewRequest("GET", proxyURL+"/v1/models", nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode == 200 && strings.Contains(string(body), `"id"`) {
				return cancel, nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	cancel()
	return nil, fmt.Errorf("proxy started but models failed to load within 10s")
}

// startCallbackRelay forwards OAuth callbacks to the proxy.
func startCallbackRelay(cfg *config.Config) {
	target, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", cfg.Port))
	proxy := httputil.NewSingleHostReverseProxy(target)
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/anthropic/callback"
		proxy.ServeHTTP(w, r)
	})
	if err := http.ListenAndServe(":54545", mux); err != nil {
		fmt.Fprintf(os.Stderr, "callback relay error: %v\n", err)
	}
}
