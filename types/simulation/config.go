package simulation

import "testing"

// Config contains the necessary configuration flags for the simulator
type Config struct {
	GenesisFile string // custom simulation genesis file; cannot be used with params file
	ParamsFile  string // custom simulation params file which overrides any random params; cannot be used with genesis

	ExportParamsPath   string // custom file path to save the exported params JSON
	ExportParamsHeight int    // height to which export the randomly generated params
	ExportStatePath    string // custom file path to save the exported app state JSON
	ExportStatsPath    string // custom file path to save the exported simulation statistics JSON

	Seed               int64  // simulation random seed
	InitialBlockHeight int    // initial block to start the simulation
	GenesisTime        int64  // genesis time to start the simulation
	NumBlocks          int    // number of new blocks to simulate from the initial block height
	BlockSize          int    // operations per block
	ChainID            string // chain-id used on the simulation

	Lean   bool // lean simulation log output
	Commit bool // have the simulation commit

	DBBackend   string // custom db backend type
	BlockMaxGas int64  // custom max gas for block
	FuzzSeed    []byte
	T           testing.TB
}

func (c Config) Clone() Config {
	return Config{
		GenesisFile:        c.GenesisFile,
		ParamsFile:         c.ParamsFile,
		ExportParamsPath:   c.ExportParamsPath,
		ExportParamsHeight: c.ExportParamsHeight,
		ExportStatePath:    c.ExportStatePath,
		ExportStatsPath:    c.ExportStatsPath,
		Seed:               c.Seed,
		InitialBlockHeight: c.InitialBlockHeight,
		GenesisTime:        c.GenesisTime,
		NumBlocks:          c.NumBlocks,
		BlockSize:          c.BlockSize,
		ChainID:            c.ChainID,
		Lean:               c.Lean,
		Commit:             c.Commit,
		DBBackend:          c.DBBackend,
		BlockMaxGas:        c.BlockMaxGas,
	}
}
