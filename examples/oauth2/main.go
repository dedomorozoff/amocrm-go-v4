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
			"your-client-id",
			"your-client-secret",
			"https://example.com/oauth2/callback",
		),
		amocrm.WithTokenStorage(tokenStorage),
		amocrm.WithDebug(true),
	)

	ctx := context.Background()

	// Проверяем, есть ли сохраненный токен
	domain := "testsubdomain.amocrm.ru"
	hasToken, err := tokenStorage.HasToken(ctx, domain)
	if err != nil {
		log.Fatalf("Ошибка проверки токена: %v", err)
	}

	if !hasToken {
		// Первичная авторизация
		fmt.Println("=== Первичная авторизация ===")

		// Получаем URL для авторизации
		authURL, err := client.Auth.GetAuthorizationURL("random-state", "")
		if err != nil {
			log.Fatalf("Ошибка получения URL авторизации: %v", err)
		}

		fmt.Printf("Перейдите по ссылке для авторизации:\n%s\n\n", authURL)
		fmt.Print("Введите код авторизации: ")

		var authCode string
		fmt.Scanln(&authCode)

		// Обмениваем код на токен
		err = client.Auth.ExchangeCode(ctx, authCode)
		if err != nil {
			log.Fatalf("Ошибка обмена кода: %v", err)
		}

		fmt.Println("Авторизация успешна! Токен сохранен.\n")
	} else {
		fmt.Println("=== Используем сохраненный токен ===\n")
	}

	// Получаем информацию об аккаунте
	account, err := client.Account.GetWithUsers(ctx)
	if err != nil {
		log.Fatalf("Ошибка получения аккаунта: %v", err)
	}

	fmt.Printf("Аккаунт: %s\n", account.Name)
	fmt.Printf("Текущий пользователь ID: %d\n", account.CurrentUserID)

	if account.Embedded != nil && len(account.Embedded.Users) > 0 {
		fmt.Println("\nПользователи:")
		for _, user := range account.Embedded.Users {
			fmt.Printf("- %s (%s)\n", user.Name, user.Email)
		}
	}

	// Получаем список контактов
	fmt.Println("\n=== Список контактов ===")
	contacts, err := client.Contacts.List(ctx, &amocrm.ContactsFilter{
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Ошибка получения контактов: %v", err)
	}

	fmt.Printf("Найдено контактов: %d\n", len(contacts))
	for i, contact := range contacts {
		fmt.Printf("%d. %s (ID: %d)\n", i+1, contact.Name, contact.ID)
	}

	// Получаем список сделок
	fmt.Println("\n=== Список сделок ===")
	leads, err := client.Leads.List(ctx, &amocrm.LeadsFilter{
		Limit: 5,
	})
	if err != nil {
		log.Fatalf("Ошибка получения сделок: %v", err)
	}

	fmt.Printf("Найдено сделок: %d\n", len(leads))
	for i, lead := range leads {
		fmt.Printf("%d. %s (ID: %d, Сумма: %d)\n", i+1, lead.Name, lead.ID, lead.Price)
	}

	// Проверяем текущий токен
	currentToken := client.Auth.GetCurrentToken()
	if currentToken != nil {
		fmt.Printf("\n=== Информация о токене ===\n")
		fmt.Printf("Тип токена: %s\n", currentToken.TokenType)
		fmt.Printf("Истекает: %s\n", currentToken.ExpiresAt.Format("2006-01-02 15:04:05"))

		if currentToken.IsExpired() {
			fmt.Println("⚠️  Токен истек, будет автоматически обновлен при следующем запросе")
		} else {
			fmt.Println("✅ Токен действителен")
		}
	}

	fmt.Println("\n=== Готово! ===")
}
