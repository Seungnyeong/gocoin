package blockchain

import (
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
}

var b *blockchain
var once sync.Once

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func (b *blockchain) persist() {
	db.SaveBlockchain(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height + 1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.Blocks()
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

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDefficulty
	} else if b.Height % difiicultyInterval == 0 {
		//recalculate the difficuty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}


func Blockchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{
				Height: 0,
			}
			checkpoint := db.Checkpoint()
			if checkpoint == nil {
				b.AddBlock("Gensis Block")
			} else {
				// restore b from bytes
			
				b.restore(checkpoint)
			}
		})
	}
	return b
}


