package state

import (
	"context"
	"fmt"
	"github.com/xyths/hs"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	TrendUp      int = 1
	TrendDown        = -1
	TrendNeutral     = 0

	SignalWait  = 0
	SignalLong  = 1
	SignalSell  = -1
	SignalShort = -2
	SignalCover = 2

	keyTide   = "tide"
	keyWave   = "wave"
	keySignal = "signal"
)

type State struct {
	db   *mongo.Database
	coll *mongo.Collection

	tide int // 1: up/green, -1: down/red, 0: neutral/blue
	wave int // 1: up, -1: down, 0: neutral
	// intermediate signal/command
	Signal int // 1: long, -1: short, 0: no signal
}

func (s *State) Init(db *mongo.Database) {
	s.db = db
	s.coll = s.db.Collection("state")
}

func (s State) Format(color bool) string {
	return s.formatTide(color) + "\n" + s.formatWave(color) + "\n" + s.formatRipple(color)
}

func (s *State) Load(ctx context.Context) error {
	if err := hs.LoadKey(ctx, s.coll, keyTide, &s.tide); err != nil {
		return err
	}
	if err := hs.LoadKey(ctx, s.coll, keySignal, &s.tide); err != nil {
		return err
	}
	return nil
}

func (s *State) Clear(ctx context.Context) error {
	// long-term/intermediate/short-term status no-need to delete, will keep in there

	if err := hs.DeleteKey(ctx, s.coll, keySignal); err != nil {
		return err
	}
	return nil
}

func (s *State) Tide() int {
	return s.tide
}

func (s *State) UpdateTide(ctx context.Context, newState int) error {
	s.tide = newState
	return hs.SaveKey(ctx, s.coll, keyTide, s.tide)
}

func (s *State) Wave() int {
	return s.wave
}

func (s *State) UpdateWave(ctx context.Context, newWave int) error {
	s.wave = newWave
	return hs.SaveKey(ctx, s.coll, keyWave, s.wave)
}

func (s *State) UpdateSignal(ctx context.Context, newSignal int) error {
	s.Signal = newSignal
	return hs.SaveKey(ctx, s.coll, keySignal, s.Signal)
}

func (s State) formatTide(color bool) string {
	var status, yes, no string
	switch s.tide {
	case 1:
		status = green(color, "Green")
		yes = "long, stand aside"
		no = "short"
	case -1:
		status = red(color, "Red")
		yes = "short, stand aside"
		no = "long"
	case 0:
		status = blue(color, "Blue")
		yes = "long, short"
		no = ""
	}
	long := fmt.Sprintf("FIRST SCREEN - MARKET TIDE\n\tStatus: %s\n\tYes: %s\n\tNo: %s", status, yes, no)
	return long
}

func (s State) formatWave(color bool) string {
	var signal string
	switch s.Signal {
	case SignalLong:
		signal = green(color, "Long")
	case SignalShort:
		signal = red(color, "Short")
	case SignalWait:
		signal = ""
	}
	return fmt.Sprintf("SECOND SCREEN - MARKET WAVE\n\tSignal: %s", signal)
}

func (s State) formatRipple(color bool) string {
	return "THIRD SCREEN - MARKET RIPPLE"
}

func red(color bool, token string) string {
	if color {
		return "\033[31m" + token + "\033[0m"
	} else {
		return token
	}
}
func green(color bool, token string) string {
	if color {
		return "\033[32m" + token + "\033[0m"
	} else {
		return token
	}
}
func blue(color bool, token string) string {
	if color {
		return "\033[34m" + token + "\033[0m"
	} else {
		return token
	}
}
