package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Job struct {
	Order *Order
}

type Jobmanager struct {
	AccrualUrl string
	Jobs       chan *Job
	Cursor     *Cursor
}

func NewJobmanager(cursor *Cursor, accrualUrl string) *Jobmanager {
	return &Jobmanager{
		AccrualUrl: accrualUrl,
		Jobs:       make(chan *Job),
		Cursor:     cursor,
	}
}

func AskAccrual(url string, number string) (*AccrualResponse, int) {
	client := &http.Client{}
	getOrder, _ := http.NewRequest(http.MethodGet, url+"/api/orders/"+number, nil)
	resp, err := client.Do(getOrder)
	if err != nil {
		ErrorLog.Fatal(err)
	}
	if resp.StatusCode == 429 {
		return nil, resp.StatusCode
	}
	defer resp.Body.Close()
	result := &AccrualResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		ErrorLog.Fatal(err)
	}
	return result, resp.StatusCode
}

func (jm *Jobmanager) RunJob(job *Job) {
	response, statusCode := AskAccrual(jm.AccrualUrl, job.Order.Number)
	if statusCode == 429 {
		time.Sleep(2 * time.Second)
	}

	for response.Status != "INVALID" && response.Status != "PROCESSED" {
		response, statusCode = AskAccrual(jm.AccrualUrl, job.Order.Number)
		if statusCode == 429 {
			time.Sleep(2 * time.Second)
			continue
		}
		jm.Cursor.UpdateOrder(response)
	}
	jm.Cursor.UpdateOrder(response)
	InfoLog.Println("Job finished")
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
}
