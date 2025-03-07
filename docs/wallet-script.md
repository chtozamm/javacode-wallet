## Wallet Script

Скрипт для упрощённого взаимодействия с приложением.

### Использование

Перед использованием рекомендуется создать алиас:

```bash
alias wallet=./wallet.sh
```

Затем используйте `wallet` с одной из команд:

```bash
wallet {start|stop|clean|audit|bench|test}
```

### Команды:

- **start**:
  - Запускает контейнеры с помощью Docker Compose.
  - Применяет миграции базы данных с помощью Goose.
- **stop**:
  - Останавливает запущенные контейнеры.
- **clean**:
  - Останавливает контейнеры, удаляет образ приложения и том базы данных.
- **audit**:
  - Выполняет проверку кода, используя инструменты `go vet`, `staticcheck`, `gosec`.
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

### Тестирование

Ниже приведён результат выполнения `wallet test`:

```plaintext
# Created a new wallet with ID:
102e25ac-070b-4d4c-a13e-8ed21a06a4ce
Current balance: 0
# Deposit 500...
Current balance: 500
# Withdraw 150...
Current balance: 350
# Trying to withdraw 10000...
Insufficient funds to withdraw: balance 350, trying to withdraw 10000
# List all wallets:
[{"id":"102e25ac-070b-4d4c-a13e-8ed21a06a4ce","balance":350,"created_at":"2025-03-07T20:16:12.777108Z","updated_at":"2025-03-07T20:16:12.829288Z"}]
# Deleting the wallet...
# List all wallets after deletion:
[]
```
