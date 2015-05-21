package core

import "github.com/spf13/viper"

// ConfigStore is the storage for the configuration singleton
type ConfigStore struct {
	RethinkDbAddress string
}

var (
	// Config is the configuration singleton
	Config ConfigStore
)

func init() {
	viper.BindEnv("rethinkdbAddress", "RETHINKDB_ADDRESS")
	viper.SetDefault("rethinkdbAddress", "content")

	Config = ConfigStore{
		RethinkDbAddress: viper.GetString("rethinkdbAddress"),
	}
}
