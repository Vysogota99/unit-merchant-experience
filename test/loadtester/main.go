package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
	"time"
)

const (
	url = "http://127.0.0.1:3000"
)

func main() {
	pull := newPull(1, 5, url)
	pull.init()
	pull.start()
}

type pull struct {
	duration   time.Duration
	nWorkers   int
	context    context.Context
	cancelFunc context.CancelFunc
	tasks      chan *task
	url        string
	statistics statistics
}

type task struct {
	url  string
	data *bytes.Buffer
}

type statistics struct {
	nErrors     int
	nSuccess    int
	nRequests   int
	mu          sync.Mutex
	statusStats map[int]int
	workTime    []float64
	durTotal    float64
	avrTime     float64
	std         float64
	nWorkers    int
}

func newPull(nWorkers int, duration int, url string) *pull {
	ctx, cFunc := context.WithCancel(context.Background())
	return &pull{
		nWorkers:   nWorkers,
		context:    ctx,
		cancelFunc: cFunc,
		tasks:      make(chan *task),
		duration:   time.Second * time.Duration(duration),
		url:        url,
		statistics: statistics{
			statusStats: map[int]int{
				200: 0,
				202: 0,
				500: 0,
				425: 0,
			},
			nWorkers: nWorkers,
			workTime: make([]float64, 0),
		},
	}
}

func (p *pull) init() {
	for i := 0; i < p.nWorkers; i++ {
		go worker(p.context, p.tasks, p)
	}
}

func (p *pull) start() {
	deadline := time.After(p.duration)
	start := time.Now()

	for {
		select {
		case <-deadline:
			p.cancelFunc()
			close(p.tasks)

			p.statistics.durTotal = time.Since(start).Seconds()
			p.statistics.showStat()
			return
		default:

			data := map[string]interface{}{
				"id":  1,
				"url": "http://nginx:80/files/1.xlsx",
			}

			body, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}

			task := &task{
				url:  p.url,
				data: bytes.NewBuffer(body),
			}

			p.tasks <- task
		}
	}
}

func worker(ctx context.Context, tasks chan *task, p *pull) {
	for {
		select {
		case task := <-tasks:
			start := time.Now()

			// получение номера воркера
			statusCodePost, id, err := post(fmt.Sprintf("%s/offer", task.url), task.data, &p.statistics)
			if err != nil {
				p.statistics.mu.Lock()
				p.statistics.nErrors++
				p.statistics.mu.Unlock()

				continue
			}

			p.statistics.mu.Lock()
			p.statistics.statusStats[statusCodePost]++
			p.statistics.mu.Unlock()

			// получение результата
			ready := false
			var code int
			for ready != true {
				code, err = get(fmt.Sprintf("%s/status/%d", task.url, id), &p.statistics)
				if err != nil {
					p.statistics.mu.Lock()
					p.statistics.nErrors++
					p.statistics.mu.Unlock()

					continue
				}
				p.statistics.mu.Lock()
				p.statistics.statusStats[code]++
				p.statistics.mu.Unlock()

				if code == 200 {
					ready = true
				} else {
					p.statistics.mu.Lock()
					p.statistics.nErrors++
					p.statistics.mu.Unlock()

					continue
				}
			}

			p.statistics.mu.Lock()
			p.statistics.workTime = append(p.statistics.workTime, time.Since(start).Seconds())
			p.statistics.nSuccess++
			p.statistics.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func get(url string, stat *statistics) (int, error) {
	stat.mu.Lock()
	stat.nRequests++
	stat.mu.Unlock()

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	code := resp.StatusCode
	log.Printf("GET %s %d\n", url, code)
	return code, nil
}

func post(url string, data *bytes.Buffer, stat *statistics) (int, int, error) {
	stat.mu.Lock()
	stat.nRequests++
	stat.mu.Unlock()

	resp, err := http.Post(url, "application/json", data)
	if err != nil {
		return 0, 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		return 0, 0, fmt.Errorf("invalid status code")
	}

	code := resp.StatusCode
	log.Printf("POST %s %d\n", url, code)

	id, err := parseID(resp)
	if err != nil {
		return 0, 0, err
	}

	return code, id, nil
}

func parseID(resp *http.Response) (int, error) {
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	type result struct {
		Result map[string]int `json:"result"`
	}

	res := &result{}

	if err := json.Unmarshal(data, res); err != nil {
		return 0, err
	}

	return res.Result["ID задачи"], nil
}

func (s *statistics) showStat() {
	fmt.Printf(`
Результаты нагрузочного тетсирования

Количество воркеров: %d
Общее время тестирования: %f (с)
==================================================================
Сделано запросов: %d

Запросов с кодом '200': %d
Запросов с кодом '202': %d
Запросов с кодом '500': %d
Запросов с кодом '425': %d
==================================================================
Среднее время работы задачи для сохранения и обработки файла: %f
Дисперсия: %f
	`, s.nWorkers, s.durTotal, s.nRequests, s.statusStats[200], s.statusStats[202], s.statusStats[500], s.statusStats[425], mean(s.workTime), std(s.workTime))
}

func mean(data []float64) float64 {
	n := len(data)
	var sum float64 = 0

	for _, val := range data {
		sum += val
	}

	return sum / float64(n)
}

func std(data []float64) float64 {
	mean := mean(data)
	n := len(data)
	var sum float64 = 0
	for _, val := range data {
		sum += math.Pow(val-mean, 2)
	}

	return sum / float64(n)
}
