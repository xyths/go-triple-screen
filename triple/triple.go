package triple

import (
	"context"
	"errors"
	"fmt"
	indicator "github.com/xyths/go-indicators"
	"github.com/xyths/go-triple-screen/impulse"
	"github.com/xyths/go-triple-screen/state"
	"github.com/xyths/hs"
	"github.com/xyths/hs/broadcast"
	"github.com/xyths/hs/exchange/gateio"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type Timeframe struct {
	Long   string
	Middle string
	Short  string
}
type StrategyConf struct {
	Timeframe Timeframe
}

type Config struct {
	Exchange hs.ExchangeConf
	Mongo    hs.MongoConf
	Strategy StrategyConf
	Log      hs.LogConf
	Robots   []hs.BroadcastConf
}

type Trader struct {
	config Config

	longTermConfig impulse.Config

	Sugar  *zap.SugaredLogger
	db     *mongo.Database
	robots []broadcast.Broadcaster
	ex     Exchange
	symbol string

	intervalLong   time.Duration
	intervalMiddle time.Duration
	intervalShort  time.Duration

	state state.State
}

func NewTrader(config Config) (*Trader, error) {
	return &Trader{config: config}, nil
}

func (t *Trader) Init(ctx context.Context) error {
	l, err := hs.NewZapLogger(t.config.Log)
	if err != nil {
		return err
	}
	t.Sugar = l.Sugar()
	t.Sugar.Info("Logger initialized")
	db, err := hs.ConnectMongo(ctx, t.config.Mongo)
	if err != nil {
		return err
	}
	t.db = db
	t.Sugar.Info("database connected")
	t.state.Init(t.db)
	t.Sugar.Info("state loaded")
	if len(t.config.Exchange.Symbols) == 0 {
		return errors.New("no symbol in config.exchange")
	}
	t.symbol = t.config.Exchange.Symbols[0]
	if err := t.initEx(); err != nil {
		return err
	}
	t.Sugar.Info("exchange initialized")
	t.initRobots(ctx)
	t.Sugar.Info("robots initialized")

	t.intervalLong = time.Hour * 24 * 7
	t.intervalMiddle = time.Hour * 24
	t.intervalShort = time.Hour

	t.longTermConfig = impulse.Config{
		EmaPeriod:        12,
		MacdSlowPeriod:   26,
		MacdFastPeriod:   12,
		MacdSignalPeriod: 9,
	}

	t.Sugar.Info("trader initialized")
	return nil
}

// Start serve until ctx.Done
func (t *Trader) Start(ctx context.Context) error {
	t.Sugar.Info("trader started")
	//t.loadState(ctx)
	//t.checkState(ctx)

	t.doWork(ctx, true)
	wakeTime := time.Now()
	longTermWakeTime := wakeTime.Truncate(t.intervalLong).Add(t.intervalLong)
	wakeTime = wakeTime.Truncate(t.intervalMiddle)

	wakeTime = wakeTime.Add(t.intervalMiddle)
	sleepTime := time.Until(wakeTime)
	t.Sugar.Debugf("next check time: %s", wakeTime.String())
	t.Sugar.Debugf("next long-term check time: %s", longTermWakeTime.String())

	for {
		select {
		case <-ctx.Done():
			t.Sugar.Info(ctx.Err())
			return ctx.Err()
		case <-time.After(sleepTime):
			checkLong := false
			if wakeTime == longTermWakeTime {
				checkLong = true
			}
			t.doWork(ctx, checkLong)
			wakeTime = wakeTime.Add(t.intervalMiddle)
			sleepTime = time.Until(wakeTime)
			t.Sugar.Debugf("next check time: %s", wakeTime.String())
		}
	}
}

// stop the running service if it's not
func (t *Trader) Stop(ctx context.Context) error {
	t.Sugar.Info("trader stopped")
	return nil
}

// release all resource
func (t *Trader) Close(ctx context.Context) error {
	t.Sugar.Info("trader closed")
	return nil
}

// print state
func (t *Trader) Print(ctx context.Context) error {
	fmt.Printf("%s", t.state.Format(true))
	return nil
}

// cancel all orders, clear state in database
func (t *Trader) Clear(ctx context.Context) error {
	return nil
}

func (t *Trader) initEx() error {
	switch t.config.Exchange.Name {
	case "gate":
		t.ex = gateio.NewSpotV4(t.config.Exchange.Key, t.config.Exchange.Secret, "", t.Sugar)
	default:
		return errors.New("exchange not support")
	}
	return nil
}

func (t *Trader) initRobots(ctx context.Context) {
	for _, conf := range t.config.Robots {
		t.robots = append(t.robots, broadcast.New(conf))
	}
	t.Sugar.Info("Broadcasters initialized")
}

// doWork do real work.
// 1. check long-term status
// 2. check candle state
// 3. buy or sell (market price)
func (t *Trader) doWork(ctx context.Context, checkLong bool) {
	if checkLong {
		t.updateTide(ctx)
	}
	t.updateWave(ctx)
}

func (t *Trader) updateTide(ctx context.Context) {
	candle, err := t.ex.CandleBySize(ctx, t.symbol, t.intervalLong, 200)
	if err != nil {
		t.Sugar.Errorf("update tide error: %s", err)
		return
	}
	//t.Sugar.Debugf("candle len = %d", candle.Length())
	//for i := candle.Length() - 4; i < candle.Length(); i++ {
	//	t.Sugar.Debugf("[%d] %d %f %f %f %f %f", i,
	//		candle.Timestamp[i], candle.Open[i], candle.High[i], candle.Low[i], candle.Close[i], candle.Volume[i])
	//}
	// use close here, maybe hl2 is better?
	rules := impulse.Impulse(t.longTermConfig, candle.Close)
	l := len(rules)
	if l < 3 {
		return
	}
	newState := rules[l-2]
	//for i := 0; i < l; i++ {
	//	t.Sugar.Debugf("[%d] timestamp %d rule %d", i, candle.Timestamp[i], rules[i])
	//}
	t.Sugar.Infof("new long-term state is %d", newState)
	if err := t.state.UpdateLongState(ctx, newState); err != nil {
		t.Sugar.Errorf("update state error: %s", err)
		return
	}
}

func (t *Trader) updateWave(ctx context.Context) {
	candle, err := t.ex.CandleBySize(ctx, t.symbol, t.intervalMiddle, 2000)
	if err != nil {
		return
	}
	efi := indicator.Efi(2, candle.Close, candle.Volume)
	t.Sugar.Debugf("intermediate candle len = %d", candle.Length())
	for i := candle.Length() - 20; i < candle.Length(); i++ {
		t.Sugar.Debugf("[%d] %d %f %f %f %f %f - %f", i,
			candle.Timestamp[i], candle.Open[i], candle.High[i], candle.Low[i], candle.Close[i], candle.Volume[i], efi[i])
	}
	l := len(efi)
	if l < 2 {
		return
	}
	signal := efi[l-2]
	if signal < 0 && (t.state.LongTermRule == state.RuleLong || t.state.LongTermRule == state.RuleNeutral) {
		if err := t.state.UpdateSignal(ctx, state.SignalLong); err != nil {
			t.Sugar.Errorf("update signal (%d) error: %s", state.SignalLong, err)
		}
	} else if signal > 0 && (t.state.LongTermRule == state.RuleShort || t.state.LongTermRule == state.RuleNeutral) {
		if err := t.state.UpdateSignal(ctx, state.SignalShort); err != nil {
			t.Sugar.Errorf("update signal (%d) error: %s", state.SignalShort, err)
		}
	}
}
