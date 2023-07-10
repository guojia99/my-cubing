/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/7/10 下午2:43.
 *  * Author: guojia(https://github.com/guojia99)
 */

package svc

import (
	"os"

	json "github.com/json-iterator/go"
)

type Config struct {
	GinMode    string   `json:"GinMode"`
	Port       int      `json:"Port"`
	StaticPath string   `json:"StaticPath"`
	DB         DBConfig `json:"DB"`
}

type DBConfig struct {
	Driver string `json:"Driver"`
	DSN    string `json:"DSN"`
}

func (c *Config) Load(file string) error {
	configBody, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(configBody, &c)
	return err
}
