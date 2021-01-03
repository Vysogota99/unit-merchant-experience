package server

import (
	"fmt"
	"log"

	"sync"
	"sync/atomic"
	"time"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/data"
	"github.com/gin-gonic/gin"
)

const (
	STATUS_DOWNLOAD_FILE = "В данный момент качается файл с данными xlsx"
	STATUS_PARSE_FILE    = "В данный момент парсится файл с данными xlsx"
	STATUS_DB            = "В данный момент проиходит работа с базой данных"
	STATUS_VALIDATE      = "В данный момент происходит валидация данных"
	STATUS_SUCCESS       = "Готово"
	STATUS_ERR           = "Ошибка"
)

type scheduler struct {
	nWorkers       int
	workerID       int64
	nActiveWorkers int64
	workers        map[int64]*worker
	mu             sync.Mutex
}

type worker struct {
	id     int64
	status chan string
	result chan string
}

func newScheduler(n int) *scheduler {
	return &scheduler{
		nWorkers: n,
		workers:  make(map[int64]*worker),
	}
}



func (s *scheduler) worker(c *gin.Context, url string, ownerID int, status chan string, result chan string, workerID chan<- int64) {
	atomic.AddInt64(&s.nActiveWorkers, 1)
	defer atomic.AddInt64(&s.nActiveWorkers, -1)

	s.mu.Lock()

	s.workerID++
	worker := &worker{
		id:     s.workerID,
		status: status,
		result: result,
	}

	tmpWorkerID := s.workerID
	s.workers[tmpWorkerID] = worker
	defer delete(s.workers, tmpWorkerID)

	s.mu.Unlock()

	log.Printf("Активных воркеров: %d\n", s.nActiveWorkers)
	workerID <- worker.id
	status <- STATUS_DOWNLOAD_FILE

	filePath := fmt.Sprintf("../static/%d_%d.xlsx", time.Now().Unix(), ownerID)
	if err := data.DownloadFile(filePath, url); err != nil {
		updateStatus(&s.mu, status, STATUS_ERR)
		result <- err.Error()
		return
	}

	updateStatus(&s.mu, status, STATUS_PARSE_FILE)

	dataToValidate, err := data.ReadXLSX(filePath)
	if err != nil {
		updateStatus(&s.mu, status, STATUS_ERR)
		result <- err.Error()
		return
	}
	time.Sleep(time.Second * 5)

	updateStatus(&s.mu, status, STATUS_VALIDATE)

	dataToDB, err := data.Validate(dataToValidate)
	if err != nil {
		updateStatus(&s.mu, status, STATUS_ERR)
		result <- err.Error()
		return
	}

	time.Sleep(time.Second * 5)

	updateStatus(&s.mu, status, STATUS_SUCCESS)
	result <- fmt.Sprintf("%v", dataToDB)
	log.Printf("Горутина %d завершила работу\n", tmpWorkerID)
}

// updateStatus - обновляет статус worker-а в канале
func updateStatus(mu *sync.Mutex, status chan string, newStatus string) {
	mu.Lock()
	defer mu.Unlock()
	<-status
	status <- newStatus
}
