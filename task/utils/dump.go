package utils

import (
	"encoding/binary"
	"encoding/hex"
	"satoblock/model"
	"satoblock/utils"

	"go.uber.org/zap"
)

func ParseEndDumpUtxo(log *zap.Logger, newUtxoDataMap map[string]*model.TxoData) {
	for keyStr, data := range newUtxoDataMap {
		key := []byte(keyStr)

		log.Info("utxo",
			zap.Uint32("h", data.BlockHeight),
			zap.String("tx", utils.HashString(key[:32])),
			zap.Uint32("i", binary.LittleEndian.Uint32(key[32:])),
			zap.Uint64("v", data.Satoshi),
			zap.Int("n", len(data.ScriptType)),
		)
	}
}

func ParseEndDumpScriptType(log *zap.Logger, calcMap map[string]int) {
	for keyStr, value := range calcMap {
		key := []byte(keyStr)

		log.Info("script type",
			zap.String("s", hex.EncodeToString(key)),
			zap.Int("n", len(keyStr)),
			zap.Int("num", value),
		)
	}
}
