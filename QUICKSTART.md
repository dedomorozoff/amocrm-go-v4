# AmoCRM Go Client

## Быстрый старт

### Установка

```bash
go get github.com/dedomorozoff/amocrm-go-v4
```

### Использование

```go
package main

import (
    "context"
    "log"
    
    "github.com/ALipckin/amocrm-go-v4/amocrm"
)

func main() {
    client := amocrm.NewClient(
        amocrm.WithSubdomain("your-subdomain"),
        amocrm.WithPermanentToken("your-token"),
    )
    
    ctx := context.Background()
    
    // Получаем информацию об аккаунте
    account, err := client.Account.Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Account: %s", account.Name)
}
```

## Документация

Полная документация доступна в [README.md](README.md)

## Примеры

- [Базовое использование](examples/basic/) - работа с контактами, сделками, задачами
- [OAuth 2.0](examples/oauth2/) - авторизация через OAuth 2.0
- [Пакетные операции](examples/batch/) - массовое создание/обновление
- [Webhooks](examples/webhooks/) - работа с вебхуками

## Основные возможности

✅ OAuth 2.0 с автообновлением токенов  
✅ Долгосрочные токены  
✅ Rate limiting (7 req/s)  
✅ Типобезопасность  
✅ Контексты и таймауты  
✅ Пакетные операции  
✅ Логирование  

## Поддерживаемые сущности

- Контакты (Contacts)
- Компании (Companies)
- Сделки (Leads)
- Задачи (Tasks)
- Примечания (Notes)
- Вебхуки (Webhooks)
- Каталоги (Catalogs)
- Аккаунт (Account)

## Лицензия

MIT License - см. [LICENSE](LICENSE)
