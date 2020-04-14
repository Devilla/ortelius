// (c) 2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm_index

import (
	"encoding/json"
	"time"

	"github.com/ava-labs/gecko/ids"
	"github.com/gocraft/dbr"
)

var (
	VMName = "avm"

	OutputTypesSECP2556K1Transfer OutputType = 0x000000ff

	TXTypeBase        TxType = "base"
	TXTypeCreateAsset        = "create_asset"
	TXTypeImport             = "import"
	TXTypeExport             = "export"
)

type TxType string
type OutputType uint32

type Output struct {
	AssetID string                   `json:"assetID"`
	Output  SECP2556K1TransferOutput `json:"output"`
}

type SECP2556K1TransferOutput struct {
	Amount    uint64   `json:"amount"`
	Locktime  uint32   `json:"locktime"`
	Threshold uint32   `json:"threshold"`
	Addresses []string `json:"addresses"`
}

type BaseTx struct {
	UnsignedTx struct {
		NetworkID    uint32   `json:"networkID"`
		BlockchainID string   `json:"blockchainID"`
		Outputs      []Output `json:"outputs"`
		Inputs       []struct {
			TxID        string `json:"txID"`
			OutputIndex int    `json:"outputIndex"`
			AssetID     string `json:"assetID"`
			Input       struct {
				Amount           uint64 `json:"amount"`
				SignatureIndices []int  `json:"signatureIndices"`
			} `json:"input"`
		} `json:"inputs"`
	} `json:"unsignedTx"`
	Credentials []struct {
		Signatures [][]int `json:"signatures"`
	} `json:"credentials"`
}

// rawID represents an ids.ID in the database
type rawID []byte

// Equals returns true if and only if the two rawIDs represent the same ids.ID
func (rid rawID) Equals(oRID rawID) bool {
	return string(rid) == string(oRID)
}

// MarshalJSON returns a json-marshaled string representing the ID
func (rid rawID) MarshalJSON() ([]byte, error) {
	id, err := ids.ToID(rid)
	if err != nil {
		return nil, err
	}
	return json.Marshal(id.String())
}

// rawShortID represents an ids.ShortID in the database
type rawShortID []byte

// Equals returns true if and only if the two rawShortID represent the same
// ids.rawShortID
func (rid rawShortID) Equals(oRID rawShortID) bool {
	return string(rid) == string(oRID)
}

// MarshalJSON returns a json-marshaled string representing the ID
func (rid rawShortID) MarshalJSON() ([]byte, error) {
	id, err := ids.ToShortID(rid)
	if err != nil {
		return nil, err
	}
	return json.Marshal(id.String())
}

// transaction represents a tx in the db
type transaction struct {
	ID      rawID  `json:"id"`
	ChainID rawID  `json:"chainID"`
	Type    string `json:"type"`

	CanonicalSerialization []byte `json:"canonicalSerialization"`
	JSONSerialization      []byte `json:"jsonSerialization",db:"json_serialization"`

	InputCount  uint16    `json:"inputCount"`
	OutputCount uint16    `json:"outputCount"`
	Amount      uint64    `json:"amount"`
	CreatedAt   time.Time `json:"created_at",db:"ingested_at"`
}

// output represents a tx output in the db
type output struct {
	TransactionID rawID        `json:"transactionID"`
	OutputIndex   uint16       `json:"outputIndex"`
	AssetID       rawID        `json:"assetOD"`
	OutputType    OutputType   `json:"outputType"`
	Amount        uint64       `json:"amount"`
	Locktime      uint64       `json:"locktime"`
	Threshold     uint64       `json:"threshold"`
	Addresses     []rawShortID `json:"addresses"`

	RedeemingTransactionID []byte `json:"redeemingTransactionID"`
	RedeemingSignature     []byte `json:"redeemingSignature"`
}

// output represents an address that controls a tx output in the db
type outputAddress struct {
	TransactionID      rawID
	OutputIndex        uint16
	Address            rawShortID
	RedeemingSignature dbr.NullString
}

type asset struct {
	ID      rawID `json:"id"`
	ChainID rawID `json:"chainID"`

	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Alias        string `json:"alias"`
	Denomination uint8  `json:"denomination"`

	CurrentSupply uint64 `json:"currentSupply"`
}

type chainInfo struct {
	ID        ids.ID `json:"chainID"`
	Alias     string `json:"chainAlias"`
	VM        string `json:"vm"`
	NetworkID uint32 `json:"networkID"`

	AssetCount       int64 `json:"assetCount"`
	AddressCount     int64 `json:"addressCount"`
	TransactionCount int64 `json:"transactionCount"`
	UTXOCount        int64 `json:"utxoCount"`
}

type transactionCounts struct {
	Minute uint64 `json:"minute"`
	Hour   uint64 `json:"hour"`
	Day    uint64 `json:"day"`
	Week   uint64 `json:"week"`
	Month  uint64 `json:"month"`
	Year   uint64 `json:"year"`
	All    uint64 `json:"all"`
}

type displayTx struct {
	json.RawMessage `db:"json_serialization"`
	Timestamp       time.Time `db:"ingested_at" json:"timestamp"`
	ID              rawID     `json:"id"`
}

func (dt *displayTx) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	err := json.Unmarshal(dt.RawMessage, &m)
	if err != nil {
		return nil, err
	}
	m["id"] = dt.ID
	m["timestamp"] = dt.Timestamp.UTC()
	return json.Marshal(m)
}

type recentTx struct {
	ID        rawID     `json:"id"`
	Timestamp time.Time `db:"ingested_at" json:"timestamp"`
}