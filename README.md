# AmoCRM API Go Client v4

![AmoCRM Logo](https://www.amocrm.ru/static/images/logo.svg)

[![Go Reference](https://pkg.go.dev/badge/github.com/dedomorozoff/amocrm-go-v4.svg)](https://pkg.go.dev/github.com/dedomorozoff/amocrm-go-v4)
[![Go Report Card](https://goreportcard.com/badge/github.com/dedomorozoff/amocrm-go-v4)](https://goreportcard.com/report/github.com/dedomorozoff/amocrm-go-v4)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Go-библиотека для работы с REST API [amoCRM](https://www.amocrm.ru) **v4** с поддержкой:
- ✅ OAuth 2.0 авторизации
- ✅ Долгосрочных токенов (рекомендуется для серверных интеграций)
- ✅ Автоматического обновления токенов
- ✅ Rate limiting (троттлинг запросов)
- ✅ Контекстов и таймаутов
- ✅ Типобезопасности
- ✅ Логирования

## Установка

```bash
go get github.com/ALipckin/amocrm-go-v4
```

## Быстрый старт

### Авторизация по долгосрочному токену (рекомендуется)

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ALipckin/amocrm-go-v4/amocrm"
)

func main() {
    // Создаем клиент с долгосрочным токеном
    client := amocrm.NewClient(
        amocrm.WithSubdomain("testsubdomain"),
        amocrm.WithPermanentToken("your-permanent-token"),
    )
    
    ctx := context.Background()
    
    // Получаем информацию об аккаунте
    account, err := client.Account.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Account: %s\n", account.Name)
}
```

### Авторизация по OAuth 2.0

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/ALipckin/amocrm-go-v4/amocrm"
    "github.com/ALipckin/amocrm-go-v4/amocrm/storage"
)

func main() {
    // Создаем хранилище токенов
    tokenStorage := storage.NewFileStorage("./tokens")
    
    // Создаем клиент с OAuth 2.0
    client := amocrm.NewClient(
        amocrm.WithSubdomain("testsubdomain"),
        amocrm.WithOAuth2(
            "client-id",
            "client-secret",
            "https://example.com/oauth2/callback",
        ),
        amocrm.WithTokenStorage(tokenStorage),
    )
    
    ctx := context.Background()
    
    // Первичная авторизация с кодом
    err := client.Auth.ExchangeCode(ctx, "authorization-code")
    if err != nil {
        log.Fatal(err)
    }
    
    // Последующие запросы автоматически обновляют токен при необходимости
    contacts, err := client.Contacts.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d contacts\n", len(contacts))
}
```

## Основные возможности

### Работа с контактами

```go
// Создание контакта
contact := &amocrm.Contact{
    Name: "Иван Иванов",
    CustomFieldsValues: []amocrm.CustomFieldValue{
        {
            FieldID: 123,
            Values: []amocrm.FieldValue{
                {Value: "+79001234567", EnumCode: "WORK"},
            },
        },
    },
}

createdContact, err := client.Contacts.Create(ctx, contact)
if err != nil {
    log.Fatal(err)
}

// Получение контакта по ID
contact, err := client.Contacts.GetByID(ctx, 12345)

// Обновление контакта
contact.Name = "Петр Петров"
updatedContact, err := client.Contacts.Update(ctx, contact)

// Поиск контактов
contacts, err := client.Contacts.List(ctx, &amocrm.ContactsFilter{
    Query: "Иван",
    Limit: 50,
})

// Пакетное создание
contacts := []*amocrm.Contact{contact1, contact2, contact3}
created, err := client.Contacts.CreateBatch(ctx, contacts)
```

### Работа со сделками

```go
// Создание сделки
lead := &amocrm.Lead{
    Name:       "Новая сделка",
    Price:      100000,
    PipelineID: 1,
    StatusID:   142,
}

createdLead, err := client.Leads.Create(ctx, lead)

// Привязка контактов к сделке
err = client.Leads.LinkContacts(ctx, leadID, []int{contactID1, contactID2})

// Привязка компании
err = client.Leads.LinkCompany(ctx, leadID, companyID)
```

### Работа с компаниями

```go
company := &amocrm.Company{
    Name: "ООО Рога и Копыта",
    CustomFieldsValues: []amocrm.CustomFieldValue{
        {
            FieldID: 456,
            Values: []amocrm.FieldValue{
                {Value: "info@company.com"},
            },
        },
    },
}

createdCompany, err := client.Companies.Create(ctx, company)
```

### Работа с задачами

```go
task := &amocrm.Task{
    Text:         "Позвонить клиенту",
    CompleteTill: time.Now().Add(24 * time.Hour).Unix(),
    EntityID:     leadID,
    EntityType:   amocrm.EntityTypeLead,
    TaskTypeID:   amocrm.TaskTypeCall,
}

createdTask, err := client.Tasks.Create(ctx, task)
```

### Работа с примечаниями

```go
note := &amocrm.Note{
    EntityID:   leadID,
    EntityType: amocrm.EntityTypeLead,
    NoteType:   amocrm.NoteTypeCommon,
    Params: map[string]interface{}{
        "text": "Важное примечание",
    },
}

createdNote, err := client.Notes.Create(ctx, note)
```

### Webhooks

```go
// Добавление webhook
webhook := &amocrm.Webhook{
    Destination: "https://example.com/webhook",
    Settings: []string{"add_lead", "update_lead"},
}

err := client.Webhooks.Subscribe(ctx, webhook)

// Получение списка webhooks
webhooks, err := client.Webhooks.List(ctx)

// Удаление webhook
err = client.Webhooks.Unsubscribe(ctx, webhookID)
```

## Конфигурация

### Опции клиента

```go
client := amocrm.NewClient(
    // Обязательные
    amocrm.WithSubdomain("testsubdomain"),
    
    // Авторизация (выберите один из методов)
    amocrm.WithPermanentToken("token"),
    // или
    amocrm.WithOAuth2("client-id", "client-secret", "redirect-uri"),
    
    // Опциональные
    amocrm.WithHTTPClient(customHTTPClient),
    amocrm.WithTokenStorage(customStorage),
    amocrm.WithRateLimit(7), // запросов в секунду
    amocrm.WithTimeout(30 * time.Second),
    amocrm.WithLogger(customLogger),
    amocrm.WithDebug(true),
)
```

### Хранение токенов

#### FileStorage (по умолчанию)

```go
storage := storage.NewFileStorage("./tokens")
```

#### Собственное хранилище

Реализуйте интерфейс `TokenStorage`:

```go
type TokenStorage interface {
    Save(ctx context.Context, domain string, token *Token) error
    Load(ctx context.Context, domain string) (*Token, error)
    HasToken(ctx context.Context, domain string) (bool, error)
}
```

Пример с базой данных:

```go
type DatabaseStorage struct {
    db *sql.DB
}

func (s *DatabaseStorage) Save(ctx context.Context, domain string, token *amocrm.Token) error {
    // Сохранение в БД
    return nil
}

func (s *DatabaseStorage) Load(ctx context.Context, domain string) (*amocrm.Token, error) {
    // Загрузка из БД
    return nil, nil
}

func (s *DatabaseStorage) HasToken(ctx context.Context, domain string) (bool, error) {
    // Проверка наличия
    return false, nil
}
```

## Rate Limiting

Библиотека автоматически ограничивает количество запросов согласно рекомендациям AmoCRM (не более 7 запросов в секунду):

```go
client := amocrm.NewClient(
    amocrm.WithSubdomain("test"),
    amocrm.WithRateLimit(7), // можно изменить
)
```

## Обработка ошибок

```go
contact, err := client.Contacts.GetByID(ctx, 12345)
if err != nil {
    switch e := err.(type) {
    case *amocrm.APIError:
        fmt.Printf("API Error: %s (code: %d)\n", e.Message, e.StatusCode)
    case *amocrm.ValidationError:
        fmt.Printf("Validation Error: %s\n", e.Message)
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Логирование

```go
import "log/slog"

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

client := amocrm.NewClient(
    amocrm.WithSubdomain("test"),
    amocrm.WithLogger(logger),
    amocrm.WithDebug(true), // включает подробное логирование запросов/ответов
)
```

## Примеры

Больше примеров в директории [examples/](examples):

- [OAuth 2.0 авторизация](./examples/oauth2/)
- [Работа с контактами](./examples/contacts/)
- [Работа со сделками](./examples/leads/)
- [Webhooks](./examples/webhooks/)
- [Пакетные операции](./examples/batch/)

## Документация

Официальная документация AmoCRM API v4: https://www.amocrm.ru/developers/content/crm_platform/api-reference

## Требования

- Go 1.21 или выше

## Лицензия

MIT License - см. [LICENSE](LICENSE)

## Автор

Создано на основе [amocrm-api-php-v4](https://github.com/dedomorozoff/amocrm-api-php-v4)

## Поддержка

Если вы нашли баг или хотите предложить улучшение, создайте [issue](https://github.com/yourusername/amocrm-go/issues)
