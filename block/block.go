package block

// Block is one i3bar / i3blocks JSON block.
type Block struct {
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text,omitempty"`
	Color               string `json:"color,omitempty"`
	Background          string `json:"background,omitempty"`
	Border              string `json:"border,omitempty"`
	BorderTop           int    `json:"border_top,omitempty"`
	BorderBottom        int    `json:"border_bottom,omitempty"`
	BorderLeft          int    `json:"border_left,omitempty"`
	BorderRight         int    `json:"border_right,omitempty"`
	Separator           bool   `json:"separator,omitempty"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Markup              string `json:"markup,omitempty"` // "none" or "pango"
}

// ClickEvent is one i3bar click event line read from stdin (see i3bar protocol).
type ClickEvent struct {
	Name      string `json:"name"`
	Instance  string `json:"instance"`
	Button    int    `json:"button"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	RelativeX int    `json:"relative_x"`
	RelativeY int    `json:"relative_y"`
	OutputX   int    `json:"output_x"`
	OutputY   int    `json:"output_y"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

// BlockInterface produces an optional status block and handles bar clicks for that block.
// GetBlock: nil means omit the block. ClickBlock is invoked when a stdin click event targets this block's Name.
type BlockInterface interface {
	GetBlock() *Block
	ClickBlock(e ClickEvent)
}
