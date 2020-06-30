package types

import "encoding/json"

type BaseCfg struct {
	Name       string          `json:"name"`
	ExtCfgFile string          `json:"ext_cfg_file"`
	CfgRaw     json.RawMessage `json:"cfg"`
}

// UnmarshalData unmarshal data to struct
func (b BaseCfg) UnmarshalData(data interface{}) error {
	return json.Unmarshal(b.CfgRaw, data)
}
