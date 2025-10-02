# L0 AI Go Agent — Bitcoin Regtest Only

Этот репозиторий — сфокусированный каркас **AI-ускорителя L0** для **тестовой сети биткоина (regtest)**.
Включает:
- Go-агент (оверлей + политика + метрики Prometheus)
- Адаптер Bitcoin (RPC + ZMQ заглушки)
- Docker-полигон: 5 нод bitcoind (regtest) + 2 экземпляра агента (A/B), Prometheus + Grafana
- Скрипты: настройка регтеста, генерация блоков и транзакций, netem задержки
- ONNX-инференс-заглушка (можно заменить на вашу модель)

## Быстрый старт
```bash
cd deploy/docker
docker compose up --build -d
# инициализация regtest, кошельков, соединений и генерация блоков/tx
../scripts/bootstrap_regtest.sh
../scripts/spam_tx.sh
# Grafana: http://localhost:3000 (дашборд L0 BTC Regtest)
# Prometheus: http://localhost:9090
```
