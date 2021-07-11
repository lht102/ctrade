# ctrade
crypto trading bot

## Goal
Listen on different varieties of data source, execute trade automatically.

### Current implementation
Listen on CoinbasePro's new coin listing tweet -> create buy order in binance futures

## Running the application
Create `.env.xxx` from `.env.sample`.
```
ENV=xxx make run
```
To use binance production environment, please create `.env.prod` run the following command.
```
ENV=prod make run
```

## Disclaimer
USE THE SOFTWARE AT YOUR OWN RISK.

## License
[MIT](https://choosealicense.com/licenses/mit/)
