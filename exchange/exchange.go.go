package exchange

import (
	"context"
	"github.com/shopspring/decimal"
	"github.com/xyths/hs"
	"time"
)

type Exchange interface {
	CandleBySize(ctx context.Context, symbol string, period time.Duration, size int) (hs.Candle, error)
	AvailableBalance(ctx context.Context) (map[string]decimal.Decimal, error)
}
