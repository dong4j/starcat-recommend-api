// Package main 是 starcat-recommend-api 的入口。
//
// 本服务按 Starcat 现有 weekly / trending / sharing / wiki 四个 API 的契约实现:
// /healthz 公开健康检查, /api/v1/ping 与业务接口均走 Bearer API Key 鉴权。
package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/starcat-app/starcat-recommend-api/internal/handler"
	"github.com/starcat-app/starcat-recommend-api/internal/middleware"
	"github.com/starcat-app/starcat-recommend-api/internal/provider"
	"github.com/starcat-app/starcat-recommend-api/internal/version"
)

const (
	defaultPort            = "5005"
	defaultSimRepoEndpoint = "https://simrepo.dera.page/collections/repos/points/recommend"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("[env] no .env file found, using OS environment only")
	} else {
		log.Printf("[env] .env loaded")
	}

	port := envOrDefault("PORT", defaultPort)
	apiKeys := requiredListEnv("API_KEYS")
	simRepoAPIKey := requiredEnv("SIMREPO_API_KEY")
	simRepoEndpoint := envOrDefault("SIMREPO_ENDPOINT", defaultSimRepoEndpoint)

	baseProvider := provider.NewSimRepoProvider(simRepoEndpoint, simRepoAPIKey, nil)
	recommendProvider := provider.NewCachedProvider(
		baseProvider,
		durationEnv("CACHE_TTL_SUCCESS_SECONDS", 7*24*time.Hour),
		durationEnv("CACHE_TTL_EMPTY_SECONDS", time.Hour),
		durationEnv("CACHE_TTL_ERROR_SECONDS", 10*time.Minute),
	)

	authMW := middleware.NewBearerAuth(apiKeys)
	recommendHandler := handler.NewRecommendHandler(recommendProvider)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.Handle("GET /api/v1/ping", authMW.Wrap(handler.HandlePingV1(version.Service, version.Version)))
	mux.Handle("GET /api/v1/repos/{repo_id}/recommendations", authMW.Wrap(http.HandlerFunc(recommendHandler.HandleRecommendations)))

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Received shutdown signal, closing service...")
		os.Exit(0)
	}()

	log.Printf("starcat-recommend-api %s starting on port %s", version.Version, port)
	log.Printf("Endpoints:")
	log.Printf("  GET /api/v1/ping                              - Connectivity probe for Starcat client (auth required)")
	log.Printf("  GET /api/v1/repos/{repo_id}/recommendations  - Similar repository recommendations (auth required)")
	log.Printf("  GET /healthz                                  - Health check (public)")
	handler := middleware.CORS(mux)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func requiredEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		log.Fatalf("%s env is required", key)
	}
	return value
}

func requiredListEnv(key string) []string {
	raw := requiredEnv(key)
	return strings.Split(raw, ",")
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		log.Printf("[env] invalid %s=%q, using fallback %s", key, value, fallback)
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
