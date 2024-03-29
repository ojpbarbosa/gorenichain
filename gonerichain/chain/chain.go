package chain

import (
	"crypto/sha512"
	"fmt"
	"gonerichain/block"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Chain[T comparable] struct {
	First, Last *block.Block[T]
	Length      int64
}

func New[T comparable]() *Chain[T] {
	chain := new(Chain[T])

	return chain
}

func Compute[T comparable](block *block.Block[T]) {
	attempt, computing := int64(0), true

	fmt.Printf("🧮 %s\n\n", block.Hash())

	for computing {
		hash := fmt.Sprintf("%x",
			sha512.Sum512(([]byte(fmt.Sprint(block.Nonce + attempt)))))

		attempt++

		if hash[:4] == "0000" {
			fmt.Printf("\n✅ %s\n\n", hash)

			computing = false
		} else {
			fmt.Printf("\t🔗 %s\n", hash)
		}
	}
}

func (chain *Chain[T]) AddBlock(data T) {
	block := block.NewBlock(chain.Last, data)

	Compute(block)

	if chain.Length == 0 {
		chain.First = block
	}

	block.Previous = chain.Last
	chain.Last = block
	chain.Length++
}

func (chain *Chain[T]) Print() {
	block := chain.Last

	for block != nil {
		block.Print()

		block = block.Previous

		if block != nil {
			fmt.Println("\t\t\t↑")
		}
	}
}

func (chain *Chain[T]) PrintBlockByData(data T) {
	block := chain.Last

	found := false

	for block != nil && !found {
		if block.Data == data {
			found = true
		} else {
			block = block.Previous
		}
	}

	if found {
		block.Print()
	} else {
		t := table.NewWriter()

		t.SetOutputMirror(os.Stdout)
		t.SetStyle(table.StyleRounded)

		t.AppendHeader(table.Row{"BLOCK NOT FOUND", "VALUE"})
		t.AppendRow(table.Row{"Data", data})

		t.Render()
	}
}
