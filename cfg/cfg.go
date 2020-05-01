// (c) 2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package cfg

import (
	"github.com/ava-labs/gecko/utils/logging"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

const appName = "ortelius"

var (
	defaultLogDirectory = "/var/log/ortelius"
	defaultRedisConf    = map[string]interface{}{
		"addr":     "redis:6379",
		"database": 0,
		"password": "",
	}
)

func getConfigViper(file string) (*viper.Viper, error) {
	v := viper.NewWithOptions(viper.KeyDelimiter("_"))

	v.SetEnvPrefix(appName)
	v.AutomaticEnv()

	v.SetConfigFile(file)
	v.SetConfigType("json")

	return v, v.ReadInConfig()
}

func getLogConf(dir string) logging.Config {
	// We ignore the error because it's related to creating the default directory
	// but we are going to override it anyways
	logConf, _ := logging.DefaultConfig()
	logConf.Directory = dir
	return logConf
}

func getKafkaConf(conf map[string]interface{}) kafka.ConfigMap {
	kc := kafka.ConfigMap{}
	for k, v := range conf {
		kc[k] = v
	}
	return kc
}

func getRedisConfig(conf *viper.Viper) (opts redis.Options) {
	if conf == nil {
		return opts
	}

	redisConf := conf.Sub("redis")
	redisConf.SetEnvPrefix(appName + "_redis")
	redisConf.AutomaticEnv()

	opts.Addr = redisConf.GetString("addr")
	opts.Password = redisConf.GetString("password")
	opts.DB = redisConf.GetInt("db")

	return opts
}