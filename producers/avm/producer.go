package avm

import (
	"encoding/json"
	"fmt"

	"github.com/ava-labs/gecko/database/nodb"
	"github.com/ava-labs/gecko/genesis"
	"github.com/ava-labs/gecko/ids"
	"github.com/ava-labs/gecko/snow"
	"github.com/ava-labs/gecko/snow/engine/common"
	"github.com/ava-labs/gecko/utils/logging"
	"github.com/ava-labs/gecko/vms/avm"
	"github.com/ava-labs/gecko/vms/platformvm"
	"github.com/ava-labs/gecko/vms/secp256k1fx"
	"github.com/ava-labs/ortelius/cfg"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// AVM produces for the AVM
type AVM struct {
	producer  *kafka.Producer
	filter    Filter
	topic     string
	chainID   ids.ID
	genesisTx *platformvm.CreateChainTx
	networkID uint32
	vmID      ids.ID
	ctx       *snow.Context
	avm       *avm.VM
}

// Initialize the producer using the configs passed as an argument
func (p *AVM) Initialize() error {
	var err error
	p.topic = cfg.Viper.GetString("chainID")
	if p.chainID, err = ids.FromString(p.topic); err != nil {
		return err
	}
	var filter map[string]interface{}
	json.Unmarshal([]byte(cfg.Viper.GetString("filter")), &filter)
	p.filter = Filter{}
	if err := p.filter.Initialize(cfg.Viper.GetString("filter")); err != nil {
		return err
	}
	p.networkID = cfg.Viper.GetUint32("networkID")
	p.vmID = avm.ID
	p.genesisTx = genesis.VMGenesis(p.networkID, p.vmID)
	p.ctx = &snow.Context{
		NetworkID: p.networkID,
		ChainID:   p.chainID,
		Log:       logging.NoLog{},
	}
	p.avm = &avm.VM{}
	echan := make(chan common.Message, 1)
	fxids := p.genesisTx.FxIDs()
	fxs := []*common.Fx{}
	for _, fxID := range fxids {
		switch {
		case fxID.Equals(secp256k1fx.ID):
			fxs = append(fxs, &common.Fx{
				Fx: secp256k1fx.Fx{},
				ID: fxID,
			})
		default:
			return fmt.Errorf("Unknown FxID: %s", secp256k1fx.ID)
		}
	}
	p.avm.Initialize(p.ctx, &nodb.Database{}, p.genesisTx.Bytes(), echan, fxs)
	kconf := cfg.Viper.Sub("kafka")
	var kafkaConf kafka.ConfigMap
	kc := kconf.AllSettings()
	for k, v := range kc {
		kafkaConf[k] = v
	}
	if p.producer, err = kafka.NewProducer(&kafkaConf); err != nil {
		return err
	}

	return nil
}

// Close shuts down the producer
func (p *AVM) Close() {
	p.producer.Close()
}

// Events returns delivery events channel
func (p *AVM) Events() chan kafka.Event {
	return p.producer.Events()
}

// Produce produces for the topic as an AVM tx
func (p *AVM) Produce(msg []byte) error {
	if p.filter.Filter(msg) {
		return nil // filter returned true, so we filter it
	}

	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}, nil)
}
