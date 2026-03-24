# Crypto Parser

CLI инструмент для получения цен криптовалют (Bitcoin, TON, Ethereum, Solana) с кэшированием через собственный API.


## Запуск с Docker

```bash
# 1. Собрать образ
make build

# 2. Запустить сервер (в фоне)
make server

# 3. Получить цены (в том же терминале)
make cli COIN=btc
make cli COIN=ton
make cli COIN=eth
make cli COIN=sol

# 4. Остановить сервер
make stop
```


## Запуск без Docker

```bash
# Терминал 1 — Сервер
go run api/server.go

#Терминал 2 — CLI
cd cmd/cli
go run main.go btc
```


## API эндпоинты

GET /prices — получить цены всех криптовалют

GET /price/{coin} — цена конкретной монеты (например, /price/bitcoin)

GET /health — статус сервера и кэша


## Архитектура

CLI → API (кэш 1 минута) → CoinGecko


## Docker команды
```bash
make build # — собрать образ

make server # — запустить сервер в фоне

make cli COIN=btc # — получить цену монеты

make stop # — остановить сервер

make status # — проверить статус сервера

make test # — полный тест
```


## Доступные монеты
btc или bitcoin — Bitcoin

ton — TON

eth или ethereum — Ethereum

sol или solana — Solana