package api

// NetworkInfo struct
type NetworkInfo struct {
	Network          string `json:"network"`
	Version          int    `json:"version"`
	Release          int    `json:"release"`
	Height           int    `json:"height"`
	Current          string `json:"current"`
	Blocks           int    `json:"blocks"`
	Peers            int    `json:"peers"`
	QueueLength      int    `json:"queue_length"`
	NodeStateLatency int    `json:"node_state_latency"`
}

// Block struct
type Block struct {
	HashList      []string      `json:"hash_list"`
	Nonce         string        `json:"nonce"`
	PreviousBlock string        `json:"previous_block"`
	Timestamp     int           `json:"timestamp"`
	LastRetarget  int           `json:"last_retarget"`
	Diff          string        `json:"diff"`
	Height        int           `json:"height"`
	Hash          string        `json:"hash"`
	IndepHash     string        `json:"indep_hash"`
	Txs           []interface{} `json:"txs"`
	WalletList    string        `json:"wallet_list"`
	RewardAddr    string        `json:"reward_addr"`
	Tags          []interface{} `json:"tags"`
	RewardPool    string        `json:"reward_pool"`
	WeaveSize     string        `json:"weave_size"`
	BlockSize     string        `json:"block_size"`
}

var allowedFields = map[string]bool{
	"id":        true,
	"last_tx":   true,
	"owner":     true,
	"target":    true,
	"quantity":  true,
	"type":      true,
	"data":      true,
	"reward":    true,
	"signature": true,
	"data.html": true,
}
