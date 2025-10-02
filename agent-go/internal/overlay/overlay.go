package overlay

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/pebbe/zmq4"
	"github.com/rs/zerolog/log"

	"github.com/example/l0-btc-regtest-ai-agent/agent-go/internal/metrics"
)

type Collector struct {
	mode   string
	reg    *metrics.Registry
	rawtxE string
}

func NewCollector(mode, zmqRawTxEndpoint string, reg *metrics.Registry) *Collector {
	return &Collector{mode: mode, rawtxE: zmqRawTxEndpoint, reg: reg}
}

func (c *Collector) Start() error {
	if c.rawtxE == "" {
		return errors.New("empty zmq endpoint")
	}
	sub, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		return err
	}
	if err := sub.Connect(c.rawtxE); err != nil {
		return err
	}
	if err := sub.SetSubscribe("rawtx"); err != nil {
		return err
	}
	log.Info().Str("endpoint", c.rawtxE).Msg("ZMQ subscribed to rawtx")

	go func() {
		defer sub.Close()
		for {
			parts, err := sub.RecvMessageBytes(0)
			if err != nil {
				log.Warn().Err(err).Msg("zmq recv failed")
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if len(parts) < 2 {
				continue
			}
			raw := parts[1]
			h1 := sha256.Sum256(raw)
			h2 := sha256.Sum256(h1[:])
			rev := make([]byte, len(h2))
			for i := 0; i < len(h2); i++ {
				rev[i] = h2[len(h2)-1-i]
			}
			txid := hex.EncodeToString(rev)
			c.reg.ObserveFirstSeen(txid, time.Now())
		}
	}()
	return nil
}
