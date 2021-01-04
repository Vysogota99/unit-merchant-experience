package server

import (
	"fmt"
	"log"

	"sync"
	"time"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/data"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store"
)

const (
	STATUS_DOWNLOAD_FILE = "В данный момент качается файл с данными xlsx"
	STATUS_PARSE_FILE    = "В данный момент парсится файл с данными xlsx"
	STATUS_DB            = "В данный момент проиходит работа с базой данных"
	STATUS_VALIDATE      = "В данный момент происходит валидация данных"
	STATUS_SUCCESS       = "Готово"
	STATUS_ERR           = "Ошибка"
	STATUS_SLEEP         = "Спит"
)

type scheduler struct {
	nWorkers int
	mu       sync.Mutex
	workers  map[int]*worker
	tasks    chan *task
	store    store.Store
}

type worker struct {
	id     int
	status chan string
	result chan string
	store  store.Store
}

type task struct {
	ownerID  int
	url      string
	workerID chan int
}

func newScheduler(n int, store store.Store) *scheduler {
	return &scheduler{
		nWorkers: n,
		workers:  make(map[int]*worker),
		tasks:    make(chan *task),
		store:    store,
	}
}

func (s *scheduler) initPull() {
	for i := 0; i < s.nWorkers; i++ {
		worker := worker{
			id:     i,
			status: make(chan string, 1),
			result: make(chan string),
			store:  s.store,
		}

		worker.status <- STATUS_SLEEP
		s.workers[i] = &worker
		go worker.start(&s.tasks)
	}
}

// start - запускает определенный воркер. В качестве аргумента принимает указатель на канал с указателем на задачу.
//			Таким образом все worker-ы ссылаются на один канал и берут оттуда задчи. При получении task, ворвер сразу
//			сохраняет в него свой id, для того, чтобы по нему получить результаты работы из map workers
func (w *worker) start(tasks *chan *task) {
	for task := range *tasks {
		log.Printf("Worker #%d начал работу", w.id)

		task.workerID <- w.id

		w.updateStatus(STATUS_DOWNLOAD_FILE)

		ids, err := w.store.Offer().GetOffersIDSBySalerID(task.ownerID)
		if err != nil {
			w.updateStatus(STATUS_ERR)
			w.result <- err.Error()
			log.Printf("Worker #%d закончил работу с ошибкой: %s", w.id, err.Error())

			w.updateStatus(STATUS_SLEEP)
			continue
		}

		log.Println(ids)
		

		filePath := fmt.Sprintf("../static/%d_%d.xlsx", time.Now().Unix(), task.ownerID)
		if err := data.DownloadFile(filePath, task.url); err != nil {
			w.updateStatus(STATUS_ERR)
			w.result <- err.Error()
			log.Printf("Worker #%d закончил работу с ошибкой: %s", w.id, err.Error())

			w.updateStatus(STATUS_SLEEP)
			continue
		}

		w.updateStatus(STATUS_PARSE_FILE)
		dataToValidate, err := data.ReadXLSX(filePath)
		if err != nil {
			w.updateStatus(STATUS_ERR)
			w.result <- err.Error()
			log.Printf("Worker #%d закончил работу с ошибкой: %s", w.id, err.Error())

			w.updateStatus(STATUS_SLEEP)
			continue
		}

		w.updateStatus(STATUS_VALIDATE)

		dataToDB, err := data.Validate(dataToValidate)
		if err != nil {
			w.updateStatus(STATUS_ERR)
			w.result <- err.Error()
			log.Printf("Worker #%d закончил работу с ошибкой: %s", w.id, err.Error())

			w.updateStatus(STATUS_SLEEP)
			continue
		}

		w.updateStatus(STATUS_SUCCESS)
		log.Println(dataToDB)
		log.Printf("Worker #%d закончил работу успешно", w.id)
	}
}

// updateStatus - обновляет статус worker-а
func (w *worker) updateStatus(newStatus string) {
	<-w.status
	w.status <- newStatus
}
