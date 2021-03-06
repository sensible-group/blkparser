package script

import (
	"bytes"
	"encoding/binary"
	"satoblock/utils"
)

var empty = make([]byte, 1)

func ExtractPkScriptGenesisIdAndAddressPkh(txid, pkscript []byte) (isNFT bool, codeHash, genesisId, addressPkh []byte, value uint64) {
	scriptLen := len(pkscript)
	if scriptLen < 2048 {
		return false, empty, empty, empty, 0
	}
	dataLen := 0
	genesisIdLen := 0
	genesisOffset := scriptLen - 8 - 4
	valueOffset := scriptLen - 8 - 4 - 8
	addressOffset := scriptLen - 8 - 4 - 8 - 20

	if (bytes.HasSuffix(pkscript, []byte("sensible")) || bytes.HasSuffix(pkscript, []byte("oraclesv"))) &&
		pkscript[scriptLen-8-4] == 1 { // PROTO_TYPE == 1

		if pkscript[scriptLen-72-36-1-1] == 0x4c && pkscript[scriptLen-72-36-1] == 108 {
			genesisIdLen = 36        // new ft
			dataLen = 1 + 1 + 1 + 72 // opreturn + 0x4c + pushdata + data
		} else if pkscript[scriptLen-72-20-1-1] == 0x4c && pkscript[scriptLen-72-20-1] == 92 {
			genesisIdLen = 20        // old ft
			dataLen = 1 + 1 + 1 + 72 // opreturn + 0x4c + pushdata + data
		} else if pkscript[scriptLen-50-36-1-1] == 0x4c && pkscript[scriptLen-50-36-1] == 86 {
			genesisIdLen = 36        // old ft
			dataLen = 1 + 1 + 1 + 50 // opreturn + 0x4c + pushdata + data
		} else if pkscript[scriptLen-92-20-1-1] == 0x4c && pkscript[scriptLen-92-20-1] == 112 {
			genesisIdLen = 20        // old ft
			dataLen = 1 + 1 + 1 + 92 // opreturn + 0x4c + pushdata + data
		} else {
			genesisIdLen = 40        // error ft
			dataLen = 1 + 1 + 1 + 72 // opreturn + 0x4c + pushdata + data
		}

		genesisOffset -= genesisIdLen
		valueOffset -= genesisIdLen
		addressOffset -= genesisIdLen

	} else if pkscript[scriptLen-1] < 2 && pkscript[scriptLen-37-1] == 37 && pkscript[scriptLen-37-1-40-1] == 40 && pkscript[scriptLen-37-1-40-1-1] == OP_RETURN {
		// nft issue
		isNFT = true
		genesisIdLen = 40
		genesisOffset = scriptLen - 37 - 1 - genesisIdLen
		valueOffset = scriptLen - 1 - 8
		addressOffset = scriptLen - 1 - 8 - 8 - 20

		dataLen = 1 + 1 + 1 + 37 // opreturn + pushdata + pushdata + data
	} else if pkscript[scriptLen-1] == 1 && pkscript[scriptLen-61-1] == 61 && pkscript[scriptLen-61-1-40-1] == 40 && pkscript[scriptLen-61-1-40-1-1] == OP_RETURN {
		// nft transfer
		isNFT = true
		genesisIdLen = 40
		genesisOffset = scriptLen - 61 - 1 - genesisIdLen
		valueOffset = scriptLen - 1 - 32 - 8
		addressOffset = scriptLen - 1 - 32 - 8 - 20

		dataLen = 1 + 1 + 1 + 61 // opreturn + pushdata + pushdata + data
	} else {
		return false, empty, empty, empty, 0
	}

	genesisId = make([]byte, genesisIdLen)
	addressPkh = make([]byte, 20)
	copy(genesisId, pkscript[genesisOffset:genesisOffset+genesisIdLen])
	copy(addressPkh, pkscript[addressOffset:addressOffset+20])

	value = binary.LittleEndian.Uint64(pkscript[valueOffset : valueOffset+8])

	codeHash = utils.GetHash160(pkscript[:scriptLen-genesisIdLen-dataLen])

	// logger.Log.Info("sensible",
	// 	// zap.String("script", hex.EncodeToString(pkscript)),
	// 	zap.String("txid", utils.HashString(txid)),
	// 	zap.String("hash", hex.EncodeToString(codeHash)),
	// 	zap.String("genesis", hex.EncodeToString(genesisId)),
	// 	zap.String("address", hex.EncodeToString(addressPkh)),
	// 	zap.Bool("nft", isNFT),
	// 	zap.Uint64("v", value),
	// 	zap.Uint8("last_byte", pkscript[scriptLen-1]),
	// )

	return isNFT, codeHash, genesisId, addressPkh, value
}
