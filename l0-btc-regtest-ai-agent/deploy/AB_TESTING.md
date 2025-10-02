# A/B тестирование
Запускаются два агента:
- `agent_control` — базовая политика
- `agent_experiment` — та же политика, но с флагом EXPERIMENT_GROUP=experiment (впаяйте сюда свои AI-настройки)

Сравнивайте метрики на дашборде Grafana **L0 BTC Regtest**.
