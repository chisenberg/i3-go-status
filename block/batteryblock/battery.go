package batteryblock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/chisenberg/i3-go-status/block"
)

// Provider shows capacity (and status when present) for a Linux sysfs power supply, e.g. id "BAT0".
type Provider struct {
	id string
}

// New returns a BlockInterface for the given battery power_supply name under /sys/class/power_supply/.
func New(id string) *Provider {
	return &Provider{id: id}
}

// GetBlock implements block.BlockInterface.
func (p *Provider) GetBlock() *block.Block {
	if p == nil || !validBatteryID(p.id) {
		return nil
	}
	base := filepath.Join("/sys/class/power_supply", p.id)
	if st, err := os.Stat(base); err != nil || !st.IsDir() {
		return nil
	}
	pct, ok := readCapacityPercent(base)
	if !ok {
		return nil
	}
	status := readOneLine(filepath.Join(base, "status"))
	icon := statusIcon(status)

	full := fmt.Sprintf("%d%%", pct)
	short := fmt.Sprintf("%d%%", pct)
	if icon != "" {
		full = icon + " " + full
		short = icon + " " + short
	}

	b := &block.Block{
		Name:      "battery-" + p.id,
		Instance:  p.id,
		FullText:  full,
		ShortText: short,
	}
	if pct <= 15 {
		b.Color = "#FF5555"
	} else if pct <= 30 {
		b.Color = "#F1FA8C"
	}
	return b
}

// ClickBlock implements block.BlockInterface.
func (*Provider) ClickBlock(block.ClickEvent) {}

func validBatteryID(id string) bool {
	if id == "" || id == "." || id == ".." {
		return false
	}
	if filepath.Base(id) != id {
		return false
	}
	for _, r := range id {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

func readCapacityPercent(base string) (int, bool) {
	capPath := filepath.Join(base, "capacity")
	data, err := os.ReadFile(capPath)
	if err == nil {
		s := strings.TrimSpace(string(data))
		n, err := strconv.Atoi(s)
		if err != nil || n < 0 || n > 100 {
			return 0, false
		}
		return n, true
	}
	// Some hardware exposes only charge_* (µAh or similar).
	nowB, err1 := os.ReadFile(filepath.Join(base, "charge_now"))
	fullB, err2 := os.ReadFile(filepath.Join(base, "charge_full"))
	if err1 != nil || err2 != nil {
		nowB, err1 = os.ReadFile(filepath.Join(base, "energy_now"))
		fullB, err2 = os.ReadFile(filepath.Join(base, "energy_full"))
	}
	if err1 != nil || err2 != nil {
		return 0, false
	}
	now, err1 := strconv.ParseFloat(strings.TrimSpace(string(nowB)), 64)
	full, err2 := strconv.ParseFloat(strings.TrimSpace(string(fullB)), 64)
	if err1 != nil || err2 != nil || full <= 0 {
		return 0, false
	}
	pct := int((now / full) * 100)
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	return pct, true
}

func readOneLine(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

// statusIcon maps Linux sysfs power_supply status strings to a single Nerd Fonts
// Material Design glyph (DejaVuSansMono Nerd)
func statusIcon(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "charging":
		return "\U000f0084" // md-battery_charging
	case "discharging":
		return "\U000f008e" // md-battery_outline
	case "not charging":
		return "\U000f06a5" // md-power_plug
	case "full":
		return "\U000f012c" // md-check
	case "unknown":
		return "\U000f02d6" // md-help
	default:
		return ""
	}
}
