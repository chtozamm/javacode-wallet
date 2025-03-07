## Wallet Script

Скрипт для упрощённого взаимодействия с приложением.

### Использование

Перед использованием рекомендуется создать алиас:

```bash
alias wallet=./wallet.sh
```

Затем используйте `wallet` с одной из команд:

```bash
wallet {start|stop|bench|clean}
```

### Команды:

- **start**:
  - Запускает контейнеры с помощью Docker Compose.
  - Применяет миграции базы данных с помощью Goose.
- **stop**:
  - Останавливает запущенные контейнеры.
- **bench**:
  - Создаёт новый кошелёк через API.
  - Выполняет нагрузочное тестирование с помощью Apache Bench, отправляя 10,000 запросов с 1,000 параллельными соединениями.
  - Удаляет созданный кошелёк после завершения тестирования.
- **test**:
  - Создаёт новый кошелёк.
  - Проверяет баланс.
  - Пополняет кошелёк.
  - Снимает средства с кошелька.
  - Получает JSON массив с созданными кошельками.
  - Удаляет созданный кошелёк.
- **clean**:
  - Останавливает контейнеры, удаляет образ приложения и том базы данных.

## Тестирование

Ниже приведёт результат выполнения `wallet test`:

```plaintext
Your wallet_id: 1d4bb96c-9c86-48e4-b500-c67ae466fb2e
Balance: 0
Balance after deposit: 500
Balance after withdrawal: 350
Response when trying to withdraw more than the wallet has: Insufficient funds to withdraw: balance 350, trying to withdraw 10000
All wallets: [{"id":"1d4bb96c-9c86-48e4-b500-c67ae466fb2e","balance":350,"created_at":"2025-03-06T14:07:02.195524Z","updated_at":"2025-03-06T14:07:02.258609Z"}]
Wallet deleted successfully
All wallets after deletion: []
```

