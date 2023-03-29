package res

type Success struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

type Error struct {
	Code    uint32 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
