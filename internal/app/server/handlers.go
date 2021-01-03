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

	status := <-worker.status
	worker.status <- status

	r.scheduler.mu.Unlock()

	if status == STATUS_ERR {
		respond(c, http.StatusInternalServerError, map[string]string{
			"status": status,
			"result": <-worker.result,
		}, "")

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
