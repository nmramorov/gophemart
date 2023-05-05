package main

import (

	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Job struct {
	orderNumber string
	username    string
}

type Jobmanager struct {
	AccrualURL string
	Jobs       chan *Job
	Cursor     *Cursor
	mu         sync.Mutex
	client *resty.Client
}

func NewJobmanager(cursor *Cursor, accrualURL string) *Jobmanager {
	return &Jobmanager{
		AccrualURL: accrualURL,
		Jobs:       make(chan *Job),
		Cursor:     cursor,
		client:     resty.New().SetBaseURL(accrualURL),
	}
}

func (jm *Jobmanager) AskAccrual(url string, number string) (*AccrualResponse, int) {
	acc := AccrualResponse{}
	req := jm.client.R().
		SetResult(&acc).
		SetPathParam("number", number)

	resp, err := req.Get("/api/orders/{number}")
	if err != nil {
		ErrorLog.Fatalf("Error getting order from accrual: %e", err)
	}
	InfoLog.Printf("Accrual GET status code: %d", resp.StatusCode())
	if resp.StatusCode() == 429 {
		return nil, resp.StatusCode()
	}
	if resp.StatusCode() == 204 {
		return &AccrualResponse{Status: "NEW"}, 204
	}
	return &acc, resp.StatusCode()
}

func (jm *Jobmanager) RunJob(job *Job) {
	response, statusCode := jm.AskAccrual(jm.AccrualURL, job.orderNumber)
	if statusCode == 429 {
		time.Sleep(time.Second)
	}
	for response.Status != "INVALID" && response.Status != "PROCESSED" {
		response, statusCode = jm.AskAccrual(jm.AccrualURL, job.orderNumber)
		if statusCode == 429 {
			time.Sleep(time.Second)
			continue
		}
		jm.mu.Lock()
		jm.Cursor.UpdateOrder(job.username, response)
		jm.mu.Unlock()
	}
	jm.mu.Lock()
	jm.Cursor.UpdateOrder(job.username, response)
	jm.Cursor.UpdateUserBalance(job.username, &Balance{
		Current:   response.Accrual,
		Withdrawn: 0.0,
	})
	jm.mu.Unlock()
	InfoLog.Println("Job finished")
}

func (jm *Jobmanager) AddJob(orderNumber string, username string) {
	jm.Jobs <- &Job{orderNumber: orderNumber, username: username}
}

func (jm *Jobmanager) ManageJobs(accrualURL string) {
	for job := range jm.Jobs {
		InfoLog.Printf("Running job for order %s", job.orderNumber)
		go jm.RunJob(job)
	}
	close(jm.Jobs)
}
