package config

// DBCfg cfg for database connect
type DBCfg struct {
	Address  string `json:"address"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Cfg struct {
	DB DBCfg `json:"db"`
}
