package main

import (
	"encoding/json"
	"net/http"
	// "sync"
	"time"
)

type Job struct {
	orderNumber string
}

type Jobmanager struct {
	AccrualURL string
	Jobs       chan *Job
	Cursor     *Cursor
	// mu         sync.Mutex
	client     *http.Client
}

func NewJobmanager(cursor *Cursor, accrualURL string) *Jobmanager {
	return &Jobmanager{
		AccrualURL: accrualURL,
		Jobs:       make(chan *Job),
		Cursor:     cursor,
		client:     &http.Client{},
	}
}

func (jm *Jobmanager) AskAccrual(url string, number string) (*AccrualResponse, int) {
	// InfoLog.Printf("calling accrual to get order %s by %s", number, url)
	getOrder, _ := http.NewRequest(http.MethodGet, url+"/api/orders/"+number, nil)
	resp, err := jm.client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatalf("Error getting order from accrual: %e", err)
	}
	if resp.StatusCode == 429 {
		return nil, resp.StatusCode
	}
	// InfoLog.Println(resp)
	defer resp.Body.Close()
	result := &AccrualResponse{}
	// InfoLog.Println(resp.Body)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		ErrorLog.Printf("Error decoding accrual response: %e", err)
		return &AccrualResponse{Status: "NEW"}, 500
	}
	return result, resp.StatusCode
}

func (jm *Jobmanager) RunJob(job *Job) {
	response, statusCode := jm.AskAccrual(jm.AccrualURL, job.orderNumber)
	if statusCode == 429 {
		time.Sleep(2 * time.Second)
	}

	for response.Status != "INVALID" && response.Status != "PROCESSED" {
		response, statusCode = jm.AskAccrual(jm.AccrualURL, job.orderNumber)
		if statusCode == 429 {
			time.Sleep(2 * time.Second)
			continue
		}
		// jm.mu.Lock()
		jm.Cursor.UpdateOrder(response)
		// jm.mu.Unlock()
	}
	// jm.mu.Lock()
	jm.Cursor.UpdateOrder(response)
	// jm.mu.Unlock()
	InfoLog.Println("Job finished")
}

func (jm *Jobmanager) AddJob(orderNumber string) {
	jm.Jobs <- &Job{orderNumber: orderNumber}
}

func (jm *Jobmanager) ManageJobs(accrualURL string) {
	for job := range jm.Jobs {
		go jm.RunJob(job)
	}
	close(jm.Jobs)
}
