package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// fileHandler - обрабатывает файлы с данными о товарах
func (r *Router) fileHandler(c *gin.Context) {
	type request struct {
		ID  int    `json:"id" binding:"required"`
		URL string `json:"url" binding:"required"`
	}

	req := request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	resChan := make(chan string, 1)
	statusChan := make(chan string, 1)
	workerIDChan := make(chan int64)
	go r.scheduler.worker(c, req.URL, req.ID, resChan, statusChan, workerIDChan)

	workerID := <-workerIDChan
	respond(c, http.StatusAccepted, map[string]int64{
		"id": workerID,
	}, "")
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
	worker, exist := r.scheduler.workers[int64(idInt)]
	if !exist {
		respond(c, http.StatusNotFound, nil, "горутины с таким id не существует")
		return
	}

	r.scheduler.mu.Lock()

	status := <-worker.status
	worker.status <- status

	r.scheduler.mu.Unlock()

	if status == STATUS_SUCCESS {
		result := <-worker.result
		respond(c, http.StatusAccepted, map[string]string{
			"status": status,
			"result": result,
		}, "")
		return
	} else if status == STATUS_ERR {
		result := <-worker.result
		respond(c, http.StatusBadRequest, nil, result)
		return
	}

	respond(c, http.StatusAccepted, map[string]string{
		"status": status,
	}, "")
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
