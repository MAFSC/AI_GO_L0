package aiassist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// ---- входные фичи (должны совпадать с FastAPI) ----
type Feats struct {
	AvgRttMs    float64 `json:"avg_rtt_ms"`
	MinRttMs    float64 `json:"min_rtt_ms"`
	MaxRttMs    float64 `json:"max_rtt_ms"`
	AvgLossPct  float64 `json:"avg_loss_pct"`
	DupRatio    float64 `json:"dup_ratio"`
	MempoolSize float64 `json:"mempool_size"`
	SockQueue   float64 `json:"sock_queue"`
	OverlayLoad float64 `json:"overlay_load"`
}

type inferResp struct {
	Delta float64 `json:"predicted_delta_ms"`
}

// ---- конфиг (из ENV) ----
type Config struct {
	DefaultBatch   int
	DefaultFlushMs int
	AlphaMix       float64
	MaxBatch       int
	MinFlushMs     int
	InferURL       string
	InferTimeout   time.Duration
	CacheTTL       time.Duration
}

func FromEnv() Config {
	return Config{
		DefaultBatch:   envInt("DEFAULT_BATCH", 8),
		DefaultFlushMs: envInt("DEFAULT_FLUSH_MS", 20),
		AlphaMix:       envFloat("ALPHA", 0.6),
		MaxBatch:       envInt("MAX_BATCH", 64),
		MinFlushMs:     envInt("MIN_FLUSH_MS", 5),
		InferURL:       envStr("INFER_URL", "http://localhost:18080/predict"),
		InferTimeout:   time.Duration(envInt("INFER_TIMEOUT_MS", 200)) * time.Millisecond,
		CacheTTL:       time.Duration(envInt("SUGGEST_CACHE_TTL_MS", 2000)) * time.Millisecond,
	}
}

func envStr(k, def string) string {
	if v := os.Getenv(k); v != "" { return v }
	return def
}
func envInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil { return n }
	}
	return def
}
func envFloat(k string, def float64) float64 {
	if v := os.Getenv(k); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil { return f }
	}
	return def
}

// ---- кэш результата ----
type Suggestion struct {
	Batch   int       `json:"batch"`
	FlushMs int       `json:"flush_ms"`
	DeltaMs float64   `json:"predicted_delta_ms"`
	FromAI  bool      `json:"from_ai"`
	At      time.Time `json:"at"`
}

type Agent struct {
	cfg    Config
	client *http.Client

	mu    sync.RWMutex
	cache Suggestion
}

func NewAgent(cfg Config) *Agent {
	a := &Agent{
		cfg:    cfg,
		client: &http.Client{Timeout: cfg.InferTimeout},
	}
	// начальное пустое значение (очень старое время, чтобы не сработал кэш)
	a.cache = Suggestion{At: time.Unix(0, 0)}
	return a
}

func (a *Agent) callInfer(ctx Feats) (float64, error) {
	b, _ := json.Marshal(ctx)
	req, _ := http.NewRequest("POST", a.cfg.InferURL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil { return 0, err }
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}
	var out inferResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return 0, err
	}
	return out.Delta, nil
}

func clampInt(x, lo, hi int) int {
	if x < lo { return lo }
	if x > hi { return hi }
	return x
}
func mapRange(y float64, inLo, inHi, outLo, outHi float64) float64 {
	if inHi == inLo { return outLo }
	if y < inLo { y = inLo }
	if y > inHi { y = inHi }
	r := (y - inLo) / (inHi - inLo)
	return outLo + r*(outHi-outLo)
}

// Основная функция подсказки с кэшом и фолбэком
func (a *Agent) SuggestParams(ctx Feats) (Suggestion, error) {
	// свежий кэш?
	a.mu.RLock()
	c := a.cache
	a.mu.RUnlock()
	if time.Since(c.At) < a.cfg.CacheTTL {
		return c, nil
	}

	// онлайновый инференс
	y, err := a.callInfer(ctx)
	if err != nil {
		sug := Suggestion{
			Batch:   a.cfg.DefaultBatch,
			FlushMs: a.cfg.DefaultFlushMs,
			DeltaMs: 0,
			FromAI:  false,
			At:      time.Now(),
		}
		a.mu.Lock(); a.cache = sug; a.mu.Unlock()
		return sug, err
	}

	// эвристика → базовые значения
	baseBatch := int(mapRange(y, 0, 1000, 16, 4))
	baseFlush := int(mapRange(y, 0, 1000, 10, 60))

	// смешивание с дефолтом
	b := int(a.cfg.AlphaMix*float64(a.cfg.DefaultBatch) + (1-a.cfg.AlphaMix)*float64(baseBatch))
	f := int(a.cfg.AlphaMix*float64(a.cfg.DefaultFlushMs) + (1-a.cfg.AlphaMix)*float64(baseFlush))

	sug := Suggestion{
		Batch:   clampInt(b, 1, a.cfg.MaxBatch),
		FlushMs: clampInt(f, a.cfg.MinFlushMs, 250),
		DeltaMs: y,
		FromAI:  true,
		At:      time.Now(),
	}
	a.mu.Lock(); a.cache = sug; a.mu.Unlock()
	return sug, nil
}

// Заглушка, чтобы не ругались импорты, если файл урежут.
var _ = errors.New
