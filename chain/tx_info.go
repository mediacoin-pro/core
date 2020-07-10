package chain

import "github.com/mediacoin-pro/core/common/bin"

type TxInfo struct {
	*Transaction
}

func (t TxInfo) IsNil() bool {
	return t.Transaction == nil
}

func (t TxInfo) Encode() []byte {
	if t.Transaction == nil {
		return bin.Encode(
			0,
			0,
			0,
			0,
			0,
		)
	}

	// caching tx-users info (see: Transaction.UsernameByID())
	t.Transaction.SenderNick()

	return bin.Encode(
		t.Transaction,
		t.blockNum,
		t.blockIdx,
		t.blockTs,
		t._users,
	)
}

func (t *TxInfo) Decode(data []byte) (err error) {
	var (
		blockNum uint64 //
		blockIdx int    //
		blockTs  int64  //
		users    map[uint64]string
	)
	err = bin.Decode(data,
		&t.Transaction,
		&blockNum,
		&blockIdx,
		&blockTs,
		&users,
	)
	if t.Transaction != nil {
		t.Transaction.blockNum = blockNum
		t.Transaction.blockIdx = blockIdx
		t.Transaction.blockTs = blockTs
		t.Transaction._users = users
	}
	return
}
