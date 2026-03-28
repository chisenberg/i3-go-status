package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/chisenberg/i3-go-status/block"
	"github.com/chisenberg/i3-go-status/block/batteryblock"
	"github.com/chisenberg/i3-go-status/block/netblock"
	"github.com/chisenberg/i3-go-status/block/timeblock"
)

func main() {

	// print header
	header := map[string]any{
		"version":      1,
		"click_events": true,
	}
	h, _ := json.Marshal(header)
	fmt.Println(string(h))

	// start infinite JSON array
	fmt.Println("[")

	// create providers
	providers := []block.BlockInterface{
		batteryblock.New("BAT0"),
		netblock.New(),
		timeblock.New(),
	}

	// start reading clicks
	go readClicks(providers)

	for {
		var blocks []block.Block
		for _, p := range providers {
			if p == nil {
				continue
			}
			if b := p.GetBlock(); b != nil {
				blocks = append(blocks, *b)
			}
		}

		b, _ := json.Marshal(blocks)

		// print blocks
		fmt.Printf("%s,\n", b)

		time.Sleep(1 * time.Second)
	}
}

func readClicks(providers []block.BlockInterface) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		var e block.ClickEvent
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			continue
		}
		if e.Name == "" {
			continue
		}
		for _, p := range providers {
			if p == nil {
				continue
			}
			b := p.GetBlock()
			if b == nil || b.Name != e.Name {
				continue
			}
			p.ClickBlock(e)
		}
	}
}
