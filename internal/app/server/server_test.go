package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Vysogota99/unit-merchant-experience/internal/app/data"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store/postgres"
	"github.com/stretchr/testify/assert"
)

const (
	serverPort         = ":8081"
	connStringPostgres = "user=user password=password dbname=app sslmode=disable"
	nWorkers           = 10
)

func TestFilehandler(t *testing.T) {
	var store store.Store = postgres.New(connStringPostgres)
	scheduler := newScheduler(nWorkers)
	scheduler.initPull()
	router := NewRouter(serverPort, store, scheduler)

	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"id":  1,
		"url": "http://nginx:80/files/1.xlsx",
	}

	body, err := json.Marshal(data)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/file", bytes.NewBuffer(body))
	router.Setup().ServeHTTP(w, req)
	assert.Equal(t, http.StatusAccepted, req.Response.StatusCode)
}

func TestDownloadFile(t *testing.T) {
	ownerID := 123
	filePath := fmt.Sprintf("./%d_%d.xlsx", time.Now().Unix(), ownerID)
	url := "http://127.0.0.1/files/1.xlsx"

	err := data.DownloadFile(filePath, url)
	assert.NoError(t, err)
}
