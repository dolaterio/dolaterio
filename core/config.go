package core

import (
	"time"

	"github.com/spf13/viper"
)

// ConfigStore is the storage for the configuration singleton
type ConfigStore struct {
	RedisAddress     string
	RethinkDbAddress string
	DockerHost       string
	DockerCertPath   string
	TaskTimeout      time.Duration
}

var (
	// Config is the configuration singleton
	Config ConfigStore
)

func init() {
	viper.BindEnv("redisAddress", "REDIS_ADDRESS")

	viper.BindEnv("rethinkdbAddress", "RETHINKDB_ADDRESS")

	viper.BindEnv("dockerHost", "DOCKER_HOST")
	viper.SetDefault("dockerHost", "unix:///var/run/docker.sock")

	viper.BindEnv("dockerCertPath", "DOCKER_CERT_PATH")
	viper.SetDefault("dockerCertPath", "")

	viper.SetDefault("taskTimeout", 30*time.Second)

	Config = ConfigStore{
		RedisAddress:     viper.GetString("redisAddress"),
		RethinkDbAddress: viper.GetString("rethinkdbAddress"),
		DockerHost:       viper.GetString("dockerHost"),
		DockerCertPath:   viper.GetString("dockerCertPath"),
		TaskTimeout:      viper.GetDuration("taskTimeout"),
	}
}
