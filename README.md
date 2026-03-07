# Loglint

> **Статический анализатор логов для Go** — проверяет сообщения в `log/slog` и `go.uber.org/zap` на соответствие стандартам качества, читаемости и безопасности.

### Установка

```bash
# Установка последней версии (standalone)
go install github.com/barashF/test-linter/cmd/loglint@latest

```

### Сборка из исходников

```bash
git clone https://github.com/barashF/test-linter.git
cd test-linter
go build -o loglint ./cmd/loglint
```

# Использование

```bash

loglint ./...

#С кастомными чувствительными ключевыми словами

loglint -sensitive-keywords="api_key,secret_token,my_password" ./...

#Отключить отдельные правила

loglint -disable-rules="lowercase,special-chars" ./...

#Применить авто-исправления

loglint --fix ./...
```

---

## Возможности

| Правило                 | Описание                                      | Пример                                          | Авто-исправление |
| ----------------------- | --------------------------------------------- | ----------------------------------------------- | ---------------- |
| 🔤 **Lowercase**        | Сообщение должно начинаться с маленькой буквы | ❌ `"Starting server"` → ✅ `"starting server"` | ✅ Да            |
| 🌐 **English only**     | Только английский язык (без кириллицы)        | ❌ `"ошибка"` → ✅ `"error"`                    | ❌ Нет           |
| 🎯 **No special chars** | Без эмодзи, `!!!`, `...`, `@#$%`              | ❌ `"failed!!!"` → ✅ `"failed"`                | ✅ Да            |
| 🔐 **Sensitive data**   | Детектирует утечки паролей, токенов, секретов | ❌ `"password=123"` → 🚫                        | ❌ Нет           |

---

### Интеграция с golangci-lint

```bash
#1. Склонируй форк golangci-lint

git clone https://github.com/barashF/golangci-lint.git
cd golangci-lint

#2. Собери бинарник (требуется Go 1.26+)

go build -o golangci-lint ./cmd/golangci-lint

#Запуск с включённым линтером

golangci-lint run --enable=loglint ./...

#С применением авто-исправлений**

golangci-lint run --enable=loglint --fix ./...
```

### Пример .golangci-lint.yml

```yml

version: "2"

linters:
enable: - loglint

linters-settings:
loglint: # Кастомные чувствительные ключевые слова
sensitive-keywords: - password - api_key - secret_token - access_token - my_custom_secret

    # Отключить отдельные правила
    disable-rules:
      - lowercase      # Не проверять регистр первой буквы
      - cyrillic       # Не проверять кириллицу

```

## Примеры использования

### Проверка логов

```go
package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.Encoding = "console"

	zapLogger, err := config.Build()
	zapLogger.Info("Test")
	zapLogger.Info("password")
	if err != nil {
		log.Fatalf("failed initialize logger: %v", err)
	}
}

```

```bash
golangci-lint run --enable=loglint ./...
service-payment/cmd/main.go:17:17: log message must start with lowercase letter (loglint)
	zapLogger.Info("Test")
	               ^
service-payment/cmd/main.go:18:17: log message may contain sensitive keyword "password" (loglint)
	zapLogger.Info("password")
	               ^
2 issues:
* loglint: 2
```

### авто-исправление

```bash
golangci-lint run --fix
service-payment/cmd/main.go:18:17: log message may contain sensitive keyword "password" (loglint)
	zapLogger.Info("password")
	               ^
1 issues:
* loglint: 1

```

```go
package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	config := zap.NewProductionConfig()

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.Encoding = "console"

	zapLogger, err := config.Build()
	zapLogger.Info("test") #успешное исправление регистра
	zapLogger.Info("password")
	if err != nil {
		log.Fatalf("failed initialize logger: %v", err)
	}
}
```
