package bcclient

import (
	"github.com/mediacoin-pro/core/chain"
	"github.com/mediacoin-pro/core/chain/txobj"
)

type Client interface {
	UserByNick(nick string) (user *txobj.User, err error)
	PublishTx(tx *chain.Transaction) (err error)
}
