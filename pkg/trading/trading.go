package trading

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/lht102/ctrade/api"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

var (
	errEmptyPriceList = errors.New("empty price list")
	errSymbolNotFound = errors.New("futures symbol not found")
)

type BinanceFuturesManager struct {
	futuresClient *futures.Client
	futuresOpts   futuresOptions
	logger        *zap.Logger

	mu               sync.Mutex
	supportedSymbols map[string]futures.Symbol
}

func NewBinanceFuturesManager(futuresClient *futures.Client, logger *zap.Logger, opts ...FuturesOption) (*BinanceFuturesManager, error) {
	options := newDefaultFuturesOptions()
	for _, o := range opts {
		o.apply(&options)
	}

	supportedSymbols, err := getSymbolsInfo(futuresClient)
	if err != nil {
		return nil, err
	}

	return &BinanceFuturesManager{
		futuresClient:    futuresClient,
		futuresOpts:      options,
		logger:           logger,
		supportedSymbols: supportedSymbols,
	}, nil
}

func (m *BinanceFuturesManager) ConsumeBuySignal(buySignal api.BuySignal) error {
	symbol := buySignal.Symbol + "USDT"

	return m.createLongPosition(symbol)
}

func (m *BinanceFuturesManager) createLongPosition(symbol string) error {
	ctx := context.Background()

	futuresSymbol, err := m.getSymbol(symbol)
	if err != nil {
		return err
	}

	price, err := getPrice(m.futuresClient, symbol)
	if err != nil {
		return err
	}

	qtyPrecision := futuresSymbol.QuantityPrecision
	qty := decimal.NewFromInt(1).
		Div(price).
		Mul(decimal.NewFromFloat(m.futuresOpts.eachTradeAmountInUSD)).
		Round(int32(qtyPrecision))

	_, err = m.futuresClient.NewChangeLeverageService().
		Symbol(symbol).
		Leverage(m.futuresOpts.leverage).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("change leverage: %w", err)
	}

	createOrderService := m.futuresClient.NewCreateOrderService()

	if !m.futuresOpts.willExecuteOrder {
		m.logger.Sugar().Infof("Trying to buy %s at ~%s with %s amount", symbol, price.String(), qty.String())

		return nil
	}

	createOrderResp, err := createOrderService.
		Symbol(symbol).
		Side(futures.SideTypeBuy).
		Type(futures.OrderTypeMarket).
		Quantity(qty.String()).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("create buy order: %w", err)
	}
	m.logger.Sugar().Infof("Executed a %s buy order at ~%s with %s amount", symbol, price.String(), qty.String())

	getOrderResp, err := m.futuresClient.NewGetOrderService().
		OrderID(createOrderResp.OrderID).
		Symbol(symbol).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("get order: %w", err)
	}

	avgPrice, err := decimal.NewFromString(getOrderResp.AvgPrice)
	if err != nil {
		return fmt.Errorf("convert average price string to decimal: %w", err)
	}

	tickSize, err := decimal.NewFromString(futuresSymbol.PriceFilter().TickSize)
	if err != nil {
		return fmt.Errorf("convert tick size string to decimal: %w", err)
	}

	multiplier := decimal.NewFromFloat(m.futuresOpts.takeProfitPriceChangedPercentage).
		Div(decimal.NewFromInt(100)). // nolint: gomnd
		Add(decimal.NewFromInt(1))
	stopPrice := roundToTickSize(avgPrice.Mul(multiplier), tickSize)

	_, err = createOrderService.
		Symbol(symbol).
		Side(futures.SideTypeSell).
		Type(futures.OrderTypeTakeProfitMarket).
		TimeInForce(futures.TimeInForceTypeGTC).
		ClosePosition(true).
		StopPrice(stopPrice.String()).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("create take profit order: %w", err)
	}

	return nil
}

func (m *BinanceFuturesManager) getSymbol(symbol string) (futures.Symbol, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	s, ok := m.supportedSymbols[symbol]
	if !ok {
		return futures.Symbol{}, errSymbolNotFound
	}

	return s, nil
}

func (m *BinanceFuturesManager) UpdateSupportedSymbols() error {
	symbols, err := getSymbolsInfo(m.futuresClient)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.supportedSymbols = symbols

	return nil
}

func getSymbolsInfo(futuresClient *futures.Client) (map[string]futures.Symbol, error) {
	resp, err := futuresClient.
		NewExchangeInfoService().
		Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get exchange info: %w", err)
	}

	res := make(map[string]futures.Symbol, len(resp.Symbols))
	for _, s := range resp.Symbols {
		res[s.Symbol] = s
	}

	return res, nil
}

func getPrice(futuresClient *futures.Client, symbol string) (decimal.Decimal, error) {
	res, err := futuresClient.NewListPricesService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("binance futures list prices: %w", err)
	}

	if len(res) == 0 {
		return decimal.Decimal{}, errEmptyPriceList
	}

	p, err := decimal.NewFromString(res[0].Price)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("convert price string to decimal: %w", err)
	}

	return p, nil
}

func roundToTickSize(price decimal.Decimal, tickSize decimal.Decimal) decimal.Decimal {
	return price.DivRound(tickSize, 0).Mul(tickSize)
}
