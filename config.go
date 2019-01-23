package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type proxyConfig struct {
	BindPort  int
	EndPoints []string
}

func newProxyConfig(filepath string) (*proxyConfig, error) {
	pc := &proxyConfig{}
	if _, err := toml.DecodeFile(filepath, pc); err != nil {
		fmt.Printf("parse proxy config err: %s, path: %s\n",
			err.Error(), filepath)
		return nil, err
	}

	return pc, nil
}
