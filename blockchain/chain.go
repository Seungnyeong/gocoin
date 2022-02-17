package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/snkim/sncoin/db"
	"github.com/snkim/sncoin/utils"
)

const (
	defaultDefficulty int = 2
	difiicultyInterval int = 5
	blockInteval int = 2
	allowedRange int = 2
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height		   int `json:"height"`
	CurrentDifficulty int `json:"currentDifficulty"`
	m 				sync.Mutex
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func persistBlockchanin(b *blockchain) {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() *Block {
	block := createBlock(b.NewestHash, b.Height + 1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockchanin(b)
	return block
}

func Blocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func Txs(b *blockchain)  []*Tx {
	var txs []*Tx
	for _, block := range Blocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockchain, targetId string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetId {
			return tx
		}
	}
	return nil
}

func recalculateDifficulty(b *blockchain) int {
	allBlocks := Blocks(b)
	newestBlock := allBlocks[0]
	lastRecalulatedBlock := allBlocks[difiicultyInterval - 1]
	actualTime := (newestBlock.Timestamp / 60) - (lastRecalulatedBlock.Timestamp / 60)
	expectedTime := difiicultyInterval * blockInteval
	
	if actualTime < (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	} else {
		return b.CurrentDifficulty
	}
}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDefficulty
	} else if b.Height % difiicultyInterval == 0 {
		//recalculate the difficuty
		return recalculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func UTxOutsByAddress(address string, b *blockchain) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)
	for _, block := range Blocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}
				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}
			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}
	return uTxOuts
}

func BalanceByAddress(address string, b *blockchain) int {
	txOuts := UTxOutsByAddress(address, b)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	utils.HandleErr(json.NewEncoder(rw).Encode(b))
}


func Blockchain() *blockchain {
	once.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkpoint := db.Checkpoint()
		if checkpoint == nil {
			b.AddBlock()
		} else {
			// restore b from bytes
		
			b.restore(checkpoint)
		} 
	})
	return b
}


func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()
	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.Height = len(newBlocks)
	b.NewestHash = newBlocks[0].Hash
	persistBlockchanin(b)
	db.EmptyBlocks()

	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(block *Block) {
	b.m.Lock()
	defer b.m.Unlock()

	b.Height += 1
	b.CurrentDifficulty = block.Difficulty
	b.NewestHash = block.Hash

	persistBlockchanin(b)
	persistBlock(block)

	// mempool

}