package script

import "bytes"

var empty = make([]byte, 20)

func ExtractPkScriptGenesisIdAndAddressPkh(pkscript []byte) (genesisId, addressPkh []byte) {
	scriptLen := len(pkscript)
	if scriptLen < 1024 {
		return empty, empty
	}
	if !bytes.HasSuffix(pkscript, []byte("oraclesv")) {
		return empty, empty
	}

	genesisOffset := scriptLen - 8 - 4 - 20
	addressOffset := scriptLen - 8 - 4 - 20 - 8 - 20
	if pkscript[scriptLen-8-4] != 1 {
		return empty, empty
	}

	// switch pkscript[scriptLen-8-4] {
	// case 0:
	// 	genesisOffset = genesisOffset
	// 	addressOffset = addressOffset - 40
	// case 1:
	// 	genesisOffset = genesisOffset - 20
	// 	addressOffset = addressOffset - 30
	// default:
	// 	genesisOffset = genesisOffset - 20
	// 	addressOffset = addressOffset - 20
	// }

	genesisId = make([]byte, 20)
	addressPkh = make([]byte, 20)
	copy(genesisId, pkscript[genesisOffset:genesisOffset+20])
	copy(addressPkh, pkscript[addressOffset:addressOffset+20])
	return genesisId, addressPkh
}
