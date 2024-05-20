package fact

import (
	"context"
	"fmt"
	"github.com/compico/kpi/pkg/storage"
	"log"
	"time"
)

type Service struct {
	Storage    storage.Storage[*Fact, Collection]
	ApiClient  *Client
	Repository *Repository
}

func NewService(allocSize int, repo *Repository, apiClient *Client) *Service {
	return &Service{
		Storage:    storage.New[*Fact, Collection](allocSize),
		Repository: repo,
		ApiClient:  apiClient,
	}
}

// SaveCollection Метод для сохранения фактов в память
// Если в памяти элементов больше чем n, тогда отправляет записи в БД
// И память чистится
func (s *Service) SaveCollection(ctx context.Context, collection Collection) error {
	s.Storage.Add(collection...)

	// Для тестов поставил 10. Оптимально ставить примерно 500
	if s.Storage.Len() > 10 {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return s.Repository.Insert(ctx, s.Storage.Flush())
	}

	return nil
}

// SendFactToApiServer основной метод для сохранения фактов в api сервер
//
// Берёт первый элемент из базы у которого indicator_to_mo_fact_id равен нулю
// отправляет факт в api сервер, получает id и сохраняет его в найденной записи
func (s *Service) SendFactToApiServer(c context.Context) error {
	fact := &Fact{}
	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	err := s.Repository.GetFirstUnsentFact(ctx, fact)
	if err != nil {
		return err
	}

	if fact.Id == 0 {
		return nil
	}

	ctx, cancel2 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel2()
	factResp, err := s.ApiClient.SendFact(ctx, fact)
	if err != nil {
		return err
	}

	if factResp.Status != "OK" {
		return fmt.Errorf("fact status not OK: %s, %s", factResp.Status, factResp.Messages.Error)
	}

	fact.IndicatorToMoFactId = factResp.Data.IndicatorToMoFactId

	ctx, cancel3 := context.WithTimeout(ctx, 5*time.Second)
	defer cancel3()
	err = s.Repository.UpdateFact(ctx, fact)

	log.Printf("New fact received: %v", fact.IndicatorToMoFactId)

	return s.SendFactToApiServer(c)
}
