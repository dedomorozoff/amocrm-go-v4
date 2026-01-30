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

	// Пакетное создание контактов
	fmt.Println("=== Пакетное создание контактов ===")

	contacts := []*amocrm.Contact{
		{
			Name:      "Контакт 1",
			FirstName: "Иван",
			LastName:  "Иванов",
		},
		{
			Name:      "Контакт 2",
			FirstName: "Петр",
			LastName:  "Петров",
		},
		{
			Name:      "Контакт 3",
			FirstName: "Сидор",
			LastName:  "Сидоров",
		},
	}

	createdContacts, err := client.Contacts.CreateBatch(ctx, contacts)
	if err != nil {
		log.Fatalf("Ошибка пакетного создания: %v", err)
	}

	fmt.Printf("Создано контактов: %d\n", len(createdContacts))
	for i, contact := range createdContacts {
		fmt.Printf("%d. %s (ID: %d)\n", i+1, contact.Name, contact.ID)
	}

	// Пакетное обновление
	fmt.Println("\n=== Пакетное обновление контактов ===")

	for i := range createdContacts {
		createdContacts[i].Name = fmt.Sprintf("Обновленный %s", createdContacts[i].Name)
	}

	// Преобразуем в указатели
	contactsToUpdate := make([]*amocrm.Contact, len(createdContacts))
	for i := range createdContacts {
		contactsToUpdate[i] = &createdContacts[i]
	}

	updatedContacts, err := client.Contacts.UpdateBatch(ctx, contactsToUpdate)
	if err != nil {
		log.Fatalf("Ошибка пакетного обновления: %v", err)
	}

	fmt.Printf("Обновлено контактов: %d\n", len(updatedContacts))
	for i, contact := range updatedContacts {
		fmt.Printf("%d. %s (ID: %d)\n", i+1, contact.Name, contact.ID)
	}

	// Пакетное создание сделок
	fmt.Println("\n=== Пакетное создание сделок ===")

	leads := []*amocrm.Lead{
		{Name: "Сделка 1", Price: 10000},
		{Name: "Сделка 2", Price: 20000},
		{Name: "Сделка 3", Price: 30000},
	}

	createdLeads, err := client.Leads.CreateBatch(ctx, leads)
	if err != nil {
		log.Fatalf("Ошибка создания сделок: %v", err)
	}

	fmt.Printf("Создано сделок: %d\n", len(createdLeads))
	for i, lead := range createdLeads {
		fmt.Printf("%d. %s (ID: %d, Сумма: %d)\n", i+1, lead.Name, lead.ID, lead.Price)
	}

	fmt.Println("\n=== Готово! ===")
}
