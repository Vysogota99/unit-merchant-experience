package server

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// offerHandler - обрабатывает файлы с данными о товарах
func (r *Router) offerHandler(c *gin.Context) {
	type request struct {
		ID  int    `json:"id" binding:"required"`
		URL string `json:"url" binding:"required"`
	}

	req := request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	task := task{
		url:      req.URL,
		ownerID:  req.ID,
		workerID: make(chan int),
	}

	r.scheduler.tasks <- &task
	workerID := <-task.workerID
	respond(c, http.StatusAccepted, map[string]int{
		"ID задачи": workerID,
	}, "")
}

func (r *Router) getOfferHandler(c *gin.Context) {
	offerID := c.Query("offer_id")
	salerID := c.Query("saler_id")
	offer := c.Query("offer")

	if offerID != "" {
		res, err := strconv.ParseInt(offerID, 10, 64)
		if err != nil {
			respond(c, http.StatusUnprocessableEntity, "", err.Error())
			return
		}

		if res < 0 {
			respond(c, http.StatusUnprocessableEntity, "", "Порядковый номер не может быть меньше 0")
			return
		}
	}

	if salerID != "" {
		res, err := strconv.ParseInt(salerID, 10, 64)
		if err != nil {
			respond(c, http.StatusUnprocessableEntity, "", err.Error())
			return
		}

		if res < 0 {
			respond(c, http.StatusUnprocessableEntity, "", "Порядковый номер не может быть меньше 0")
			return
		}
	}
	result, err := r.store.Offer().GetOffers(context.Background(), offerID, salerID, offer)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err.Error())
		return
	}

	respond(c, http.StatusOK, result, "")
}

func (r *Router) fileStatusHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		respond(c, http.StatusBadRequest, nil, "no id in query params")
		return
	}

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		respond(c, http.StatusInternalServerError, nil, err.Error())
		return
	}
	worker, exist := r.scheduler.workers[int(idInt)]
	if !exist {
		respond(c, http.StatusNotFound, nil, "задачи с таким id не существует")
		return
	}

	r.scheduler.mu.Lock()

	log.Printf("Получение статуса воркера #%d", idInt)
	status := <-worker.status

	// обновляем статус для успешно выполненного воркера, тем самым переводим его в
	//	режим ожидания новой задачи
	if status == STATUS_SUCCESS {
		worker.status <- STATUS_SLEEP
	} else {
		worker.status <- status
	}
	log.Printf("статус воркера #%d: %s", idInt, status)

	r.scheduler.mu.Unlock()

	if status == STATUS_ERR {
		<-worker.result
		respond(c, http.StatusInternalServerError, map[string]interface{}{
			"status": status,
		}, "")
	} else if status == STATUS_SLEEP {
		respond(c, http.StatusOK, map[string]interface{}{
			"status": status,
		}, "")
	} else if status == STATUS_SUCCESS {

		respond(c, http.StatusAccepted, map[string]interface{}{
			"status":     status,
			"статистика": <-worker.result,
		}, "")
	} else {
		respond(c, http.StatusProcessing, map[string]interface{}{
			"status": status,
		}, "")
	}
}

func respond(c *gin.Context, code int, result interface{}, err string) {
	if err == "EOF" {
		result = "Неправильное тело запроса"
	}

	c.JSON(
		code,
		gin.H{
			"result": result,
			"error":  err,
		},
	)

}
