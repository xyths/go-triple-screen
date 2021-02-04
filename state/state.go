package state

import (
	"context"
	"fmt"
	"github.com/xyths/hs"
	"go.mongodb.org/mongo-driver/mongo"
)

type State struct {
	db   *mongo.Database
	coll *mongo.Collection

	LongTermRule int // 1: long/green, -1: short/red, 0: blue/neutral
	// intermediate signal/command
	Signal int // 1: long, -1: short, 0: no signal
}

const (
	RuleLong    = 1
	RuleShort   = -1
	RuleNeutral = 0

	SignalLong  = 1
	SignalShort = -1
	SignalWait  = 0

	keyLongTermRule = "longTerm"
	keySignal       = "signal"
)

func (s *State) Init(db *mongo.Database) {
	s.db = db
	s.coll = s.db.Collection("state")
}

func (s State) Format(color bool) string {
	return s.formatLongTerm(color) + "\n" + s.formatMiddleTerm(color) + "\n" + s.formatShortTerm(color)
}

func (s *State) Load(ctx context.Context) error {
	if err := hs.LoadKey(ctx, s.coll, keyLongTermRule, &s.LongTermRule); err != nil {
		return err
	}
	if err := hs.LoadKey(ctx, s.coll, keySignal, &s.LongTermRule); err != nil {
		return err
	}
	return nil
}

func (s *State) Clear(ctx context.Context) error {
	// long-term/intermediate/short-term status no-need to delete, will keep in there
	//if err := hs.LoadKey(ctx, s.coll, "longTerm", &s.LongTermRule); err != nil {
	//	return err
	//}
	if err := hs.DeleteKey(ctx, s.coll, keySignal); err != nil {
		return err
	}
	return nil
}

func (s *State) UpdateLongState(ctx context.Context, newState int) error {
	s.LongTermRule = newState
	return hs.SaveKey(ctx, s.coll, keyLongTermRule, s.LongTermRule)
}

func (s *State) UpdateSignal(ctx context.Context, newSignal int) error {
	s.Signal = newSignal
	return hs.SaveKey(ctx, s.coll, keySignal, s.Signal)
}

func (s State) formatLongTerm(color bool) string {
	var status, yes, no string
	switch s.LongTermRule {
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

func (s State) formatMiddleTerm(color bool) string {
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

func (s State) formatShortTerm(color bool) string {
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
