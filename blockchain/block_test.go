package blockchain

import (
	"reflect"
	"testing"

	"github.com/snkim/sncoin/utils"
)

func TestCreateBlock(t *testing.T){
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("x", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBolcok() should return an instance of a block")
	}
}

func TestFindBlock(t *testing.T) {
	t.Run("Block is fount", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeFindBlock: func() []byte {
				v := &Block{
					Height: 1,
				}
				return utils.ToBytes(v)
			},
		}
		block , _ := FindBlock("xx")
		if block.Height != 1 {
			t.Error("Block should be found")
		}
	})
}