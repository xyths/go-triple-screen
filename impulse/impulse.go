package impulse

import (
	"github.com/markcheno/go-talib"
	indicator "github.com/xyths/go-indicators"
)

// the impulse system

type Config struct {
	EmaPeriod        int
	MacdFastPeriod   int
	MacdSlowPeriod   int
	MacdSignalPeriod int
}

func Impulse(config Config, inReal []float64) []int {
	ema := talib.Ema(inReal, config.EmaPeriod)
	_, _, hist := talib.Macd(inReal, config.MacdFastPeriod, config.MacdSlowPeriod, config.MacdSignalPeriod)
	return indicator.Impulse(ema, hist)
}
