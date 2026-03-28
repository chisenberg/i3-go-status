package timeblock

import (
	"time"

	"github.com/chisenberg/i3-go-status/block"
)

// Time renders the current clock in the status bar.
type Time struct{}

// NewTime returns a BlockInterface for the time block.
func New() *Time {
	return &Time{}
}

// GetBlock implements BlockInterface.
func (*Time) GetBlock() *block.Block {
	return &block.Block{
		Name:     "time",
		FullText: time.Now().Format("15:04:05"),
		Color:    "#FFFFFF",
	}
}

// ClickBlock implements BlockInterface.
func (*Time) ClickBlock(block.ClickEvent) {}
