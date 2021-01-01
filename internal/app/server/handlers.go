package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// fileHandler - обрабатывает файлы с данными о товарах
func (r *Router) fileHandler(c *gin.Context) {
	type request struct {
		ID  int    `json:"id" binding:"required"`
		URL string `json:"url" binding:"required"`
	}

	req := request{}
	if err := c.ShouldBindJSON(req); err != nil {
		respond(c, http.StatusBadRequest, nil, err.Error())
		return
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
