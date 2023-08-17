package xlog

import "time"

type XLog struct {
	CreatedTime time.Time `json:"CreatedTime"`
	Area        string    `json:"Area"`
	Title       string    `json:"Title"`
	Messages    string    `json:"Messages"`
}
