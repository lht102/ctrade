package trading

const (
	defaultTakeProfitPriceChangedPercentage = 5.0
	defaultEachTradeAmountInUSD             = 500.0
	defaultLeverage                         = 5
)

type FuturesOption interface {
	apply(*futuresOptions)
}

type futuresOptions struct {
	takeProfitPriceChangedPercentage float64
	eachTradeAmountInUSD             float64
	leverage                         int
	willExecuteOrder                 bool
}

func newDefaultFuturesOptions() futuresOptions {
	return futuresOptions{
		takeProfitPriceChangedPercentage: defaultTakeProfitPriceChangedPercentage,
		eachTradeAmountInUSD:             defaultEachTradeAmountInUSD,
		leverage:                         defaultLeverage,
		willExecuteOrder:                 false,
	}
}

type takeProfitPriceChangedPercentageOption float64

func (c takeProfitPriceChangedPercentageOption) apply(opts *futuresOptions) {
	opts.takeProfitPriceChangedPercentage = float64(c)
}

func WithTakeProfitPriceChangedPercentage(f float64) FuturesOption {
	return takeProfitPriceChangedPercentageOption(f)
}

type eachTradeAmountInUSDOption float64

func (c eachTradeAmountInUSDOption) apply(opts *futuresOptions) {
	opts.eachTradeAmountInUSD = float64(c)
}

func WithEachTradeAmountInUSD(f float64) FuturesOption {
	return eachTradeAmountInUSDOption(f)
}

type leverageOption int

func (c leverageOption) apply(opts *futuresOptions) {
	opts.leverage = int(c)
}

func WithLeverage(l int) FuturesOption {
	return leverageOption(l)
}

type willExecuteOrderOption bool

func (c willExecuteOrderOption) apply(opts *futuresOptions) {
	opts.willExecuteOrder = bool(c)
}

func WithWillExecuteOrder(f bool) FuturesOption {
	return willExecuteOrderOption(f)
}
