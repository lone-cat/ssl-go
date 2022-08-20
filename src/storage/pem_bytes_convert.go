package storage

import (
	"encoding/pem"
	"errors"
)

func PEMBlockToBytes(pemBlock *pem.Block) (bts []byte, err error) {
	if pemBlock == nil {
		err = errors.New(`nil pem block passed`)
		return
	}
	bts = pem.EncodeToMemory(pemBlock)
	return
}

func PEMBlocksToBytesSlice(pemBlocks []*pem.Block) (data [][]byte, err error) {
	if pemBlocks == nil {
		err = errors.New(`nil pem blocks slice passed`)
	}
	data = make([][]byte, 0)
	var bts []byte
	for _, pemBlock := range pemBlocks {
		bts, err = PEMBlockToBytes(pemBlock)
		if err != nil {
			break
		}
		data = append(data, bts)
	}

	return
}

func BytesToPEMBlocks(data []byte) (pemBlocks []*pem.Block, err error) {
	left := data
	pemBlocks = make([]*pem.Block, 0)
	var pemBlock *pem.Block
	for {
		pemBlock, left = pem.Decode(left)
		if pemBlock == nil {
			break
		}

		pemBlocks = append(pemBlocks, pemBlock)
	}

	return
}

func BytesSliceToPEMBlocks(data [][]byte) (pemBlocks []*pem.Block, err error) {
	if data == nil {
		err = errors.New(`nil data passed`)
		return
	}

	pemBlocks = make([]*pem.Block, 0)
	var pemBlockSlice []*pem.Block
	for _, bts := range data {
		pemBlockSlice, err = BytesToPEMBlocks(bts)
		if err != nil {
			break
		}

		pemBlocks = append(pemBlocks, pemBlockSlice...)
	}

	return
}
