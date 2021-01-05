package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/data"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store/mock"
	"github.com/Vysogota99/unit-merchant-experience/internal/app/store/postgres"
	"github.com/stretchr/testify/assert"
)

const (
	serverPort         = ":8081"
	connStringPostgres = "user=user1 password=password dbname=app sslmode=disable"
	nWorkers           = 10
)

func TestFilehandler(t *testing.T) {
	db, err := sql.Open("postgres", connStringPostgres)
	store := postgres.New(db)

	scheduler := newScheduler(nWorkers, store)
	scheduler.initPull()
	router := NewRouter(serverPort, store, scheduler)

	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"id":  1,
		"url": "http://127.0.0.1:80/files/1.xlsx",
	}

	body, err := json.Marshal(data)
	assert.NoError(t, err)

	req, _ := http.NewRequest("POST", "/offer", bytes.NewBuffer(body))
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

func TestGetOffers(t *testing.T) {
	db, mockDB, err := sqlmock.New()
	assert.NoError(t, err)
	store := mock.New(db, mockDB)

	scheduler := newScheduler(nWorkers, store)
	router := NewRouter(serverPort, store, scheduler)

	type testCase struct {
		name   string
		params map[string]string
		code   int
	}

	tCases := []testCase{
		testCase{
			name: "test1",
			params: map[string]string{
				"offer_id": "",
				"saler_id": "",
				"offer":    "",
			},
			code: 200,
		},
		testCase{
			name: "test2",
			params: map[string]string{
				"offer_id": "-1",
				"saler_id": "asd",
				"offer":    "",
			},
			code: 422,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/offer?offer_id=%s&saler_id=%s&offer=%s", tc.params["offer_id"], tc.params["saler_id"], tc.params["offer"])
			req, _ := http.NewRequest("GET", url, nil)
			router.Setup().ServeHTTP(w, req)
		})
	}
}
