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
		amocrm.WithPermanentToken("your-permanent-token-here"),
		amocrm.WithDebug(true), // включаем отладочный режим
	)

	ctx := context.Background()

	// Получаем информацию об аккаунте
	fmt.Println("=== Получение информации об аккаунте ===")
	account, err := client.Account.Get(ctx)
	if err != nil {
		log.Fatalf("Ошибка получения аккаунта: %v", err)
	}
	fmt.Printf("Аккаунт: %s (ID: %d)\n", account.Name, account.ID)
	fmt.Printf("Поддомен: %s\n", account.Subdomain)
	fmt.Printf("Валюта: %s\n\n", account.Currency)

	// Создаем контакт
	fmt.Println("=== Создание контакта ===")
	contact := &amocrm.Contact{
		Name:      "Иван Иванов",
		FirstName: "Иван",
		LastName:  "Иванов",
		CustomFieldsValues: []amocrm.CustomFieldValue{
			{
				FieldID: 123456, // ID поля "Телефон" (замените на реальный)
				Values: []amocrm.FieldValue{
					{
						Value:    "+79001234567",
						EnumCode: "WORK",
					},
				},
			},
			{
				FieldID: 123457, // ID поля "Email" (замените на реальный)
				Values: []amocrm.FieldValue{
					{
						Value: "ivan@example.com",
					},
				},
			},
		},
	}

	createdContact, err := client.Contacts.Create(ctx, contact)
	if err != nil {
		log.Fatalf("Ошибка создания контакта: %v", err)
	}
	fmt.Printf("Создан контакт: %s (ID: %d)\n\n", createdContact.Name, createdContact.ID)

	// Получаем список контактов
	fmt.Println("=== Получение списка контактов ===")
	contacts, err := client.Contacts.List(ctx, &amocrm.ContactsFilter{
		Limit: 10,
		Query: "Иван",
	})
	if err != nil {
		log.Fatalf("Ошибка получения контактов: %v", err)
	}
	fmt.Printf("Найдено контактов: %d\n", len(contacts))
	for _, c := range contacts {
		fmt.Printf("- %s (ID: %d)\n", c.Name, c.ID)
	}
	fmt.Println()

	// Создаем сделку
	fmt.Println("=== Создание сделки ===")
	lead := &amocrm.Lead{
		Name:       "Новая сделка",
		Price:      100000,
		PipelineID: 1,   // ID воронки (замените на реальный)
		StatusID:   142, // ID статуса (замените на реальный)
	}

	createdLead, err := client.Leads.Create(ctx, lead)
	if err != nil {
		log.Fatalf("Ошибка создания сделки: %v", err)
	}
	fmt.Printf("Создана сделка: %s (ID: %d, Сумма: %d)\n\n", createdLead.Name, createdLead.ID, createdLead.Price)

	// Привязываем контакт к сделке
	fmt.Println("=== Привязка контакта к сделке ===")
	err = client.Leads.LinkContacts(ctx, createdLead.ID, []int{createdContact.ID})
	if err != nil {
		log.Fatalf("Ошибка привязки контакта: %v", err)
	}
	fmt.Printf("Контакт %d привязан к сделке %d\n\n", createdContact.ID, createdLead.ID)

	// Создаем задачу
	fmt.Println("=== Создание задачи ===")
	task := &amocrm.Task{
		Text:         "Позвонить клиенту",
		EntityID:     createdLead.ID,
		EntityType:   string(amocrm.EntityTypeLead),
		TaskTypeID:   int(amocrm.TaskTypeCall),
		CompleteTill: 1735689600, // Unix timestamp (замените на актуальный)
	}

	createdTask, err := client.Tasks.Create(ctx, task)
	if err != nil {
		log.Fatalf("Ошибка создания задачи: %v", err)
	}
	fmt.Printf("Создана задача: %s (ID: %d)\n\n", createdTask.Text, createdTask.ID)

	// Добавляем примечание к сделке
	fmt.Println("=== Добавление примечания ===")
	note := &amocrm.Note{
		EntityID: createdLead.ID,
		NoteType: amocrm.NoteTypeCommon,
		Params: map[string]interface{}{
			"text": "Клиент заинтересован в продукте",
		},
	}

	createdNote, err := client.Notes.Create(ctx, amocrm.EntityTypeLead, note)
	if err != nil {
		log.Fatalf("Ошибка создания примечания: %v", err)
	}
	fmt.Printf("Создано примечание (ID: %d)\n\n", createdNote.ID)

	// Обновляем контакт
	fmt.Println("=== Обновление контакта ===")
	createdContact.Name = "Иван Петрович Иванов"
	updatedContact, err := client.Contacts.Update(ctx, createdContact)
	if err != nil {
		log.Fatalf("Ошибка обновления контакта: %v", err)
	}
	fmt.Printf("Обновлен контакт: %s (ID: %d)\n\n", updatedContact.Name, updatedContact.ID)

	fmt.Println("=== Готово! ===")
}
