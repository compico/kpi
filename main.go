package main

import (
	"context"
	"github.com/compico/kpi/internal/fact"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Вместо обычного main.go обычно использую cli тулзы по типу urfave/cli
func main() {
	// Обычно вывожу генерацию в отдельный файл конфига
	dsn := "admin:111111@tcp(localhost:3306)/app?charset=utf8&parseTime=True&loc=Local"

	// Использовал ORM так-как не особо хотел напрягаться в написании SQL запросов
	// В целом инструмент хороший, но не покрывает все бизнес-задачи, приходится его дополнять
	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 500,
	})
	if err != nil {
		log.Fatalf("connect database error: %s", err.Error())
	}

	// Можно было-бы сделать другой способ миграции, обычно использую goose для этого
	// Но это заняло ещё бы где-то 4 часа
	err = orm.AutoMigrate(&fact.Fact{})
	if err != nil {
		log.Fatalf("auto migrate error: %s", err.Error())
	}

	// Для таких вещей обычно использую google/wire, но что-бы не тратить лишний час
	// решил инициализировать всё в main.go
	factRepository := fact.NewRepository(orm)

	// Сюда вписать токен
	kpiApiClient := fact.NewClient(os.Getenv("API_TOKEN"))
	factService := fact.NewService(2000, factRepository, kpiApiClient)

	// тут технически подошла бы любая библиотека для http севера,
	// в т.ч стандартный net/http,
	// но так-как сейчас изучаю fiber - использовал его
	app := fiber.New()

	// В качестве handler - можно сделать отдельную структуру, которая её возвращает
	// опять же для ускорения разработки, всё вывел в main.go
	app.Post("/api/v1/fact", func(c fiber.Ctx) error {
		// аллоцирую память под коллекцию фактов, цифра по сути из головы
		factRequest := make(fact.Collection, 0, 500)

		// Тут происходит unmarshal json'а
		// По-хорошему нужно использовать валидации, но решил не запариваться
		// ещё есть хороший инструмент как easyjson от mailru, но пока не изучал его
		err := c.Bind().JSON(&factRequest)
		if err != nil {
			return err
		}

		if len(factRequest) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"status": "empty_request"})
		}

		go func() {
			ctx := context.Background()
			// Сохраняю всё сначала в память
			// Тут особенности mysql. Очень глупо на каждый http-запрос
			// делать INSERT в базу, если можно это как-то по быстрому оптимизировать
			err := factService.SaveCollection(ctx, factRequest)
			if err != nil {
				log.Printf("save collection error: %s", err.Error())
			}
		}()

		return c.Status(fiber.StatusOK).JSON(map[string]string{"status": "ok"})
	})

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("app start error: %s", err.Error())
		}
	}()

	// Создаю тикеры для разных задач

	// saver будет брать данные из памяти и класть их в базу
	saver := time.NewTicker(5 * time.Second)
	// sender будет брать данные из базы и отправлять их на api сервер
	sender := time.NewTicker(1 * time.Second)
	done := make(chan struct{}, 1)
	ctx, stopSendingFact := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-saver.C:
				ctx := context.Background()
				err := factRepository.Insert(ctx, factService.Storage.Flush())
				if err != nil {
					log.Printf("insert error: %s", err.Error())
				}
			case <-sender.C:
				err := factService.SendFactToApiServer(ctx)
				if err != nil {
					log.Printf("send fact to api error: %s", err.Error())
				}
			case <-done:
				return
			}
		}
	}()

	// Тут мы говорим программе,
	// что пока мы не получим по stdin сигнал Ctrl+C/Ctrl+D - ничего не делаем
	// тот же сигнал делает кубер при остановке пода
	// или докер при остановке контейнера
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	// Как только получаем сигнал
	<-sigCh

	close(done)
	stopSendingFact()
	err = app.Server().Shutdown()

	if err != nil {
		log.Fatalf("server shutdown error: %s", err.Error())
	}
}
