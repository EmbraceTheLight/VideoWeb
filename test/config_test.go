package test

import (
	"VideoWeb/config"
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	cfg, err := config.ParseConfig("")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cfg)
}
