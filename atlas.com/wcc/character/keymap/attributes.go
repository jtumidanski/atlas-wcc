package keymap

type attributes struct {
	Key    int32 `json:"key"`
	Type   int8  `json:"type"`
	Action int32 `json:"action"`
}
