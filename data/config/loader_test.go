package config

import "testing"

func TestInitConfig(t *testing.T) {
	InitConfig("../../conf/config.yaml.template")
}

func BenchmarkInitConfig(b *testing.B) {
	InitConfig("../../conf/config.yaml.template")
}
