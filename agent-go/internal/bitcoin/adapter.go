package bitcoin

import (
	"context"
	"time"
	"net/http"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"github.com/rs/zerolog/log"
	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/metrics"
)

type Config struct {
	RPCHost string
	RPCPort string
	RPCUser string
	RPCPass string
	ZMQRawTx string
	ZMQRawBlk string
}

type Adapter struct {
	cfg Config
	reg *metrics.Registry
	httpc *http.Client
}

func NewAdapter(cfg Config, reg *metrics.Registry) *Adapter {
	return &Adapter{
		cfg: cfg,
		reg: reg,
		httpc: &http.Client{ Timeout: 5 * time.Second },
	}
}

// Run subscribes to ZMQ and polls RPC periodically (placeholders).
// fb is feedback sink to policy.
func (a *Adapter) Run(ctx context.Context, fb func(latMs float64, dups int)) {
	t := time.NewTicker(2*time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// simple RPC ping to mempoolinfo as placeholder
			start := time.Now()
			_, err := a.rpc("getmempoolinfo", []any{})
			lat := float64(time.Since(start).Milliseconds())
			if err != nil {
				log.Warn().Err(err).Msg("rpc getmempoolinfo failed")
			}
			fb(lat, 0)
		}
	}
}

type rpcReq struct{
	Jsonrpc string `json:"jsonrpc"`
	ID int `json:"id"`
	Method string `json:"method"`
	Params any `json:"params"`
}

func (a *Adapter) rpc(method string, params any) (json.RawMessage, error) {
	url := fmt.Sprintf("http://%s:%s@%s:%s/",
		a.cfg.RPCUser, a.cfg.RPCPass, a.cfg.RPCHost, a.cfg.RPCPort)
	body, _ := json.Marshal(rpcReq{Jsonrpc:"1.0", ID:1, Method:method, Params:params})
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.httpc.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	all, _ := io.ReadAll(resp.Body)
	var out struct{ Result json.RawMessage; Error any }
	if err := json.Unmarshal(all, &out); err != nil { return nil, err }
	if out.Error != nil { return nil, fmt.Errorf("rpc error: %v", out.Error) }
	return out.Result, nil
}
