package config

var Config *Conf

type Conf struct {
	Mongo      MongoConfig
	Blockchain BlockchainConfig
}

type MongoConfig struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	Username string `toml:"user_name"`
	Password string `toml:"password"`
	DbName   string `toml:"db_name"`
}

type BlockchainConfig struct {
	RPC        string `toml:"rpc"`
	TtContract string `toml:"tt_contract"`
}
