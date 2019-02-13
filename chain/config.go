package chain

import (
	"flag"
	"os"

	"github.com/mediacoin-pro/core/crypto"
)

type Config struct {
	NetworkID      int
	ChainID        uint64
	MasterKey      string
	Dir            string
	VerifyTxsLevel int

	_mkey *crypto.PublicKey
}

var DefaultConfig = &Config{
	ChainID:        1,
	NetworkID:      1,
	MasterKey:      "0x04662a375cfa894f8c22adcce4f0c8949750da9df48942b35ef8f1bc42cafe6fcc4b3d5103a0ee7ff2ca0e3e3845c1df639f81f3d1140ad81840e9af8495d27ea6",
	VerifyTxsLevel: VerifyTxLevel1,
}

func NewConfig() *Config {
	cfg := *DefaultConfig
	cfg.Dir = os.Getenv("HOME") + "/mdc.bc"

	flag.StringVar(&cfg.Dir, "bc-dir", cfg.Dir, "blockchain dir")
	return &cfg
}

func (c *Config) MasterPubKey() *crypto.PublicKey {
	if c._mkey == nil {
		c._mkey = crypto.MustParsePublicKey(c.MasterKey)
	}
	return c._mkey
}

const (
	VerifyTxLevel1 = 1
)

func GenesisBlockHeader(cfg *Config) *BlockHeader {
	return &BlockHeader{
		Version: 0,
		Network: cfg.NetworkID,
		ChainID: cfg.ChainID,
	}
}
