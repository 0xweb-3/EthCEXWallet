package config

type ChainId struct {
	Scroll          uint64 `mapstructure:"scroll"`
	Polygon         uint64 `mapstructure:"polygon"`
	PolygonSepolia  uint64 `mapstructure:"polygon_sepolia"`
	Ethereum        uint64 `mapstructure:"ethereum"`
	EthereumSepolia uint64 `mapstructure:"ethereum_sepolia"`
	Base            uint64 `mapstructure:"base"`
	BaseSepolia     uint64 `mapstructure:"base_sepolia"`
	Manta           uint64 `mapstructure:"manta"`
	MantaSepolia    uint64 `mapstructure:"manta_sepolia"`
	MantleSepolia   uint64 `mapstructure:"mantle_sepolia"`
	Mantle          uint64 `mapstructure:"mantle"`
	ZkFairSepolia   uint64 `mapstructure:"zk_fair_sepolia"`
	ZkFair          uint64 `mapstructure:"zk_fair"`
	OkxSepolia      uint64 `mapstructure:"okx_sepolia"`
	Okx             uint64 `mapstructure:"okx"`
	Op              uint64 `mapstructure:"op"`
	OpTest          uint64 `mapstructure:"op_test"`
	Linea           uint64 `mapstructure:"linea"`
	arb             uint64 `mapstructure:"arb"`
}

type Config struct {
	Name           string  `mapstructure:"name"`
	EthRpcUrl      string  `mapstructure:"eth_rpc_url"`
	MaxRequestTime int     `mapstructure:"max_request_time"`
	ChainId        ChainId `mapstructure:"chain_id"`
}
