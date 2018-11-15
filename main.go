package main

import (
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/krzysztof-gzocha/curnot/pkg/config"
)

const configFile = "config.yml"

func main() {
	cfg := config.Config{}
	err := configor.Load(&cfg, configFile)
	if err != nil {
		fmt.Printf("Could not load %s: %s\n", configFile, err.Error())
		return
	}
}
