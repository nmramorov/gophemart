package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	orderNumber string
}

type Jobmanager struct {
	AccrualUrl string
	Jobs       chan *Job
	Cursor     *Cursor
	mu         sync.Mutex
	client     *http.Client
}

func NewJobmanager(cursor *Cursor, accrualUrl string) *Jobmanager {
	return &Jobmanager{
		AccrualUrl: accrualUrl,
		Jobs:       make(chan *Job),
		Cursor:     cursor,
		client:     &http.Client{},
	}
}

func (jm *Jobmanager) AskAccrual(url string, number string) (*AccrualResponse, int) {
	getOrder, _ := http.NewRequest(http.MethodGet, url+"/api/orders/"+number, nil)
	resp, err := jm.client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatalf("Error getting order from accrual: %e", err)
	}
	if resp.StatusCode == 429 {
		return nil, resp.StatusCode
	}
	defer resp.Body.Close()
	result := &AccrualResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		ErrorLog.Fatalf("Error decoding accrual response: %e", err)
	}
	return result, resp.StatusCode
}

func (jm *Jobmanager) RunJob(job *Job) {
	response, statusCode := jm.AskAccrual(jm.AccrualUrl, job.orderNumber)
	if statusCode == 429 {
		time.Sleep(2 * time.Second)
	}

	for response.Status != "INVALID" && response.Status != "PROCESSED" {
		response, statusCode = jm.AskAccrual(jm.AccrualUrl, job.orderNumber)
		if statusCode == 429 {
			time.Sleep(2 * time.Second)
			continue
		}
		jm.mu.Lock()
		jm.Cursor.UpdateOrder(response)
		jm.mu.Unlock()
	}
	jm.mu.Lock()
	jm.Cursor.UpdateOrder(response)
	jm.mu.Unlock()
	InfoLog.Println("Job finished")
}

func (jm *Jobmanager) AddJob(orderNumber string) {
	jm.Jobs <- &Job{orderNumber: orderNumber}
}

func (jm *Jobmanager) ManageJobs(accrualUrl string) {
	// for {
	// 	select {
	// 	case job := <-jm.Jobs:
	// 		go jm.HandleJob(job)
	// 	}
	// }
	for job := range jm.Jobs {
		go jm.RunJob(job)
	}
	close(jm.Jobs)
}
