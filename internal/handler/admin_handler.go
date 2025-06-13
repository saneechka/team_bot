package handler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"team_bot/internal/model"
	"team_bot/internal/repository/sqlrepo"
	"team_bot/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AuthHandler struct {
	bot           *tgbotapi.BotAPI
	repo          *sqlrepo.AuthRepository
	adminUsers    []string
	inviteService *service.InviteService
}

func NewAuthHandler(bot *tgbotapi.BotAPI, repo *sqlrepo.AuthRepository, adminUsers []string) *AuthHandler {
	return &AuthHandler{
		bot:           bot,
		repo:          repo,
		adminUsers:    adminUsers,
		inviteService: service.NewInviteService(repo),
	}
}

func (h *AuthHandler) HandleUpdate(ctx context.Context, update *tgbotapi.Update) error {
	if update.Message == nil {
		return nil
	}

	if strings.HasPrefix(update.Message.Text, "/start ") {
		return h.HandleStartWithToken(ctx, update)
	}

	switch update.Message.Text {
	case "/start":
		return h.HandleStart(ctx, update)
	case "/help":
		return h.HandleHelp(ctx, update)
	case "/join":
		return h.HandleJoinTeam(ctx, update)
	case "/admin":

		if !h.CheckUserAccess(ctx, update.Message.From.ID, update.Message.Chat.ID) {
			return nil
		}
		return h.HandleAdmin(ctx, update)
	case "/create_invite":

		if !h.CheckUserAccess(ctx, update.Message.From.ID, update.Message.Chat.ID) {
			return nil
		}
		return h.HandleCreateInvite(ctx, update)
	case "/invite_info":

		if !h.CheckUserAccess(ctx, update.Message.From.ID, update.Message.Chat.ID) {
			return nil
		}
		return h.HandleInviteInfo(ctx, update)
	default:

		if !h.CheckUserAccess(ctx, update.Message.From.ID, update.Message.Chat.ID) {
			return nil
		}
		return h.handleUnknownCommand(ctx, update)
	}
}

func (h *AuthHandler) handleUnknownCommand(ctx context.Context, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда. Используйте /help для начала работы.")
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending unknown command message: %v", err)
	}
	return nil
}

func (h *AuthHandler) CheckAdminAccess(ctx context.Context, userID int64, chatID int64) (bool, error) {
	isAdmin, err := h.repo.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("error checking admin status: %v", err)
	}

	if !isAdmin {
		msg := tgbotapi.NewMessage(chatID, "❌ Доступ запрещён. У вас нет прав администратора.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending access denied message: %v", err)
		}
		return false, nil
	}

	return true, nil
}

func (h *AuthHandler) HandleStart(ctx context.Context, update *tgbotapi.Update) error {

	isAdmin := false
	username := update.Message.From.UserName
	for _, adminUsername := range h.adminUsers {
		if username == adminUsername {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return h.HandleJoinTeam(ctx, update)
	}

	user := &model.User{
		ID:          update.Message.From.ID,
		Username:    username,
		ChatID:      update.Message.Chat.ID,
		CreatedTime: time.Now(),
		IsAdmin:     isAdmin,
	}

	if err := h.repo.SaveUser(ctx, user); err != nil {
		log.Printf("Error saving user: %v", err)
		return fmt.Errorf("error saving user: %v", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("Привет, %s! Я бот для управления командой.\n✅ Вы зарегистрированы как администратор.", username))
	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

func (h *AuthHandler) HandleAdmin(ctx context.Context, update *tgbotapi.Update) error {
	isAdmin, err := h.repo.IsAdmin(ctx, update.Message.From.ID)
	if err != nil {
		return fmt.Errorf("error checking admin status: %v", err)
	}

	var responseText string
	if isAdmin {
		responseText = "✅ Вы являетесь администратором."
	} else {
		responseText = "❌ Вы не являетесь администратором."
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending admin status message: %v", err)
	}

	return nil
}


func (h *AuthHandler) HandleStartWithToken(ctx context.Context, update *tgbotapi.Update) error {

	parts := strings.Split(update.Message.Text, " ")
	if len(parts) != 2 {
		return h.HandleStart(ctx, update)
	}

	token := parts[1]


	inviteToken, err := h.inviteService.ValidateAndUseToken(ctx, token)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при присоединении к команде: %s", err.Error()))
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending token error message: %v", err)
		}
		return h.HandleStart(ctx, update)
	}


	username := update.Message.From.UserName
	user := &model.User{
		ID:          update.Message.From.ID,
		Username:    username,
		ChatID:      update.Message.Chat.ID,
		CreatedTime: time.Now(),
		IsAdmin:     false, 
	}

	if err := h.repo.SaveUser(ctx, user); err != nil {
		log.Printf("Error saving user: %v", err)
		return fmt.Errorf("error saving user: %v", err)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("🎉 Добро пожаловать в команду, %s!\n\n"+
			"✅ Вы успешно присоединились к команде.\n"+
			"🔗 Использований токена: %d/%d",
			username, inviteToken.UsageCount, inviteToken.MaxUsage))
	msg.ReplyToMessageID = update.Message.MessageID

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending welcome message: %v", err)
	}

	return nil
}


func (h *AuthHandler) HandleCreateInvite(ctx context.Context, update *tgbotapi.Update) error {

	hasAccess, err := h.CheckAdminAccess(ctx, update.Message.From.ID, update.Message.Chat.ID)
	if err != nil || !hasAccess {
		return err
	}


	token, err := h.inviteService.CreateInviteLink(ctx, update.Message.From.ID, 24, 50)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при создании пригласительной ссылки: %v", err))
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending error message: %v", err)
		}
		return fmt.Errorf("error creating invite link: %v", err)
	}


	botInfo, err := h.bot.GetMe()
	if err != nil {
		log.Printf("Error getting bot info: %v", err)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("🔗 <b>Пригласительная ссылка создана!</b>\n\n"+
				"<b>Токен:</b> <code>%s</code>\n"+
				"<b>Срок действия:</b> до %s\n"+
				"<b>Лимит использований:</b> %d\n\n"+
				"Отправьте этот токен пользователям для присоединения к команде.",
				token.Token,
				token.ExpiresAt.Format("02.01.2006 15:04"),
				token.MaxUsage))
		msg.ParseMode = "HTML"
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending invite link: %v", err)
		}
		return nil
	}


	inviteLink := h.inviteService.FormatInviteLink(botInfo.UserName, token.Token)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("🔗 <b>Пригласительная ссылка создана!</b>\n\n"+
			"<b>Ссылка:</b> %s\n"+
			"<b>Токен:</b> <code>%s</code>\n"+
			"<b>Срок действия:</b> до %s\n"+
			"<b>Лимит использований:</b> %d\n\n"+
			"Отправьте эту ссылку пользователям для присоединения к команде.",
			inviteLink,
			token.Token,
			token.ExpiresAt.Format("02.01.2006 15:04"),
			token.MaxUsage))
	msg.ParseMode = "HTML"

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending invite link: %v", err)
	}

	return nil
}


func (h *AuthHandler) HandleInviteInfo(ctx context.Context, update *tgbotapi.Update) error {

	hasAccess, err := h.CheckAdminAccess(ctx, update.Message.From.ID, update.Message.Chat.ID)
	if err != nil || !hasAccess {
		return err
	}


	token, err := h.inviteService.GetInviteLink(ctx)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при получении информации о ссылке: %v", err))
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending error message: %v", err)
		}
		return fmt.Errorf("error getting invite info: %v", err)
	}

	if token == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"ℹ️ Активных пригласительных ссылок нет.\n\n"+
				"Используйте /create_invite для создания новой ссылки.")
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending no invite message: %v", err)
		}
		return nil
	}


	timeLeft := time.Until(token.ExpiresAt)
	var statusText string
	if timeLeft <= 0 {
		statusText = "❌ Истек"
	} else {
		hours := int(timeLeft.Hours())
		minutes := int(timeLeft.Minutes()) % 60
		statusText = fmt.Sprintf("✅ Активна (%dч %dм)", hours, minutes)
	}

	botInfo, err := h.bot.GetMe()
	var inviteLink string
	if err != nil {
		inviteLink = "Ошибка получения ссылки"
	} else {
		inviteLink = h.inviteService.FormatInviteLink(botInfo.UserName, token.Token)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("📋 <b>Информация о пригласительной ссылке</b>\n\n"+
			"<b>Ссылка:</b> %s\n"+
			"<b>Токен:</b> <code>%s</code>\n"+
			"<b>Статус:</b> %s\n"+
			"<b>Использований:</b> %d/%d\n"+
			"<b>Создана:</b> %s\n"+
			"<b>Истекает:</b> %s",
			inviteLink,
			token.Token,
			statusText,
			token.UsageCount,
			token.MaxUsage,
			token.CreatedAt.Format("02.01.2006 15:04"),
			token.ExpiresAt.Format("02.01.2006 15:04")))
	msg.ParseMode = "HTML"

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending invite info: %v", err)
	}

	return nil
}


func (h *AuthHandler) HandleJoinTeam(ctx context.Context, update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		"🔗 <b>Присоединение к команде</b>\n\n"+
			"Для присоединения к команде вам нужна пригласительная ссылка от администратора.\n\n"+
			"<b>Как присоединиться:</b>\n"+
			"1. Получите пригласительную ссылку от администратора\n"+
			"2. Нажмите на неё или используйте команду /start с токеном\n\n"+
			"<b>Пример:</b> <code>/start abc123def456</code>")
	msg.ParseMode = "HTML"

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending join info: %v", err)
	}

	return nil
}


func (h *AuthHandler) CheckUserAccess(ctx context.Context, userID int64, chatID int64) bool {

	exists, err := h.repo.UserExists(ctx, userID)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)

		msg := tgbotapi.NewMessage(chatID, "❌ Ошибка доступа. Попробуйте команду /join для присоединения к команде.")
		if _, sendErr := h.bot.Send(msg); sendErr != nil {
			log.Printf("Error sending access error message: %v", sendErr)
		}
		return false
	}

	if !exists {

		msg := tgbotapi.NewMessage(chatID,
			"❌ Доступ запрещён. Вы не зарегистрированы в системе.\n\n"+
				"Используйте команду /join для получения информации о присоединении к команде.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending access denied message: %v", err)
		}
		return false
	}

	return true
}

func (h *AuthHandler) Start(ctx context.Context) {
	log.Println("Starting bot...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping bot...")
			h.bot.StopReceivingUpdates()
			return
		case update := <-updates:
			if err := h.HandleUpdate(ctx, &update); err != nil {
				log.Printf("Error handling update: %v", err)
			}
		}
	}
}


func (h *AuthHandler) HandleHelp(ctx context.Context, update *tgbotapi.Update) error {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	

	exists, err := h.repo.UserExists(ctx, userID)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)

		msg := tgbotapi.NewMessage(chatID, h.getGuestHelpText())
		msg.ParseMode = "HTML"
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending help message: %v", err)
		}
		return nil
	}

	if !exists {

		msg := tgbotapi.NewMessage(chatID, h.getGuestHelpText())
		msg.ParseMode = "HTML"
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending guest help message: %v", err)
		}
		return nil
	}


	isAdmin, err := h.repo.IsAdmin(ctx, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)

		msg := tgbotapi.NewMessage(chatID, h.getUserHelpText())
		msg.ParseMode = "HTML"
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending user help message: %v", err)
		}
		return nil
	}

	var helpText string
	if isAdmin {
		helpText = h.getAdminHelpText()
	} else {
		helpText = h.getUserHelpText()
	}

	msg := tgbotapi.NewMessage(chatID, helpText)
	msg.ParseMode = "HTML"
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending help message: %v", err)
	}

	return nil
}


func (h *AuthHandler) getGuestHelpText() string {
	return "🤖 <b>Помощь - Гость</b>\n\n" +
		"<b>Доступные команды:</b>\n\n" +
		"/start - Запуск бота и регистрация\n" +
		"/help - Показать эту справку\n" +
		"/join - Информация о присоединении к команде\n\n" +
		"<b>Как присоединиться к команде:</b>\n" +
		"1. Получите пригласительную ссылку от администратора\n" +
		"2. Нажмите на неё или используйте /start с токеном\n\n" +
		"<b>Пример:</b> <code>/start abc123def456</code>"
}


func (h *AuthHandler) getUserHelpText() string {
	return "🤖 <b>Помощь - Участник команды</b>\n\n" +
		"<b>Доступные команды:</b>\n\n" +
		"/start - Перезапуск бота\n" +
		"/help - Показать эту справку\n" +
		"/admin - Проверить статус администратора\n\n" +
		"<b>Статус:</b> ✅ Вы зарегистрированы как участник команды"
}


func (h *AuthHandler) getAdminHelpText() string {
	return "🤖 <b>Помощь - Администратор</b>\n\n" +
		"<b>Доступные команды:</b>\n\n" +
		"/start - Перезапуск бота\n" +
		"/help - Показать эту справку\n" +
		"/admin - Проверить статус администратора\n\n" +
		"<b>Управление приглашениями:</b>\n" +
		"/create_invite - Создать пригласительную ссылку\n" +
		"/invite_info - Информация о текущей ссылке\n\n" +
		"<b>Общие команды:</b>\n" +
		"/join - Информация о присоединении к команде\n\n" +
		"<b>Статус:</b> 👑 Вы являетесь администратором\n\n" +
		"<b>Возможности администратора:</b>\n" +
		"• Создание пригласительных ссылок (24 часа, до 100 использований)\n" +
		"• Просмотр статистики использования ссылок\n" +
		"• Управление доступом к боту"
}
