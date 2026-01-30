package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ALipckin/amocrm-go-v4/amocrm"
)

func main() {
	// Создаем клиент
	client := amocrm.NewClient(
		amocrm.WithSubdomain("testsubdomain"),
		amocrm.WithPermanentToken("your-token"),
	)

	ctx := context.Background()

	// Создаем webhook
	fmt.Println("=== Создание webhook ===")
	webhook := &amocrm.Webhook{
		Destination: "https://example.com/webhook",
		Settings: []string{
			"add_lead",
			"update_lead",
			"delete_lead",
			"add_contact",
			"update_contact",
		},
	}

	err := client.Webhooks.Subscribe(ctx, webhook)
	if err != nil {
		log.Fatalf("Ошибка создания webhook: %v", err)
	}
	fmt.Println("Webhook создан успешно")

	// Получаем список webhooks
	fmt.Println("\n=== Список webhooks ===")
	webhooks, err := client.Webhooks.List(ctx)
	if err != nil {
		log.Fatalf("Ошибка получения webhooks: %v", err)
	}

	fmt.Printf("Найдено webhooks: %d\n", len(webhooks))
	for i, wh := range webhooks {
		fmt.Printf("%d. URL: %s\n", i+1, wh.Destination)
		fmt.Printf("   События: %v\n", wh.Settings)
		fmt.Printf("   Отключен: %v\n", wh.Disabled)
		fmt.Printf("   ID: %s\n\n", wh.ID)
	}

	// Пример обработки входящего webhook
	fmt.Println("=== Обработка входящего webhook ===")
	fmt.Println("Пример структуры для обработки webhook:")

	exampleWebhookData := `
	{
		"leads": {
			"add": [
				{
					"id": 12345,
					"name": "Новая сделка",
					"status_id": 142,
					"price": 100000,
					"responsible_user_id": 123,
					"created_at": 1234567890,
					"updated_at": 1234567890
				}
			]
		},
		"account": {
			"id": 123,
			"subdomain": "testsubdomain"
		}
	}
	`

	fmt.Println(exampleWebhookData)

	fmt.Println("\n=== Готово! ===")
}
