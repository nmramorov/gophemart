package jobmanager

import (
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/logger"
	"github.com/nmramorov/gophemart/internal/models"
)

type Job struct {
	orderNumber string
	username    string
}

type Jobmanager struct {
	AccrualURL string
	Jobs       chan *Job
	Cursor     *db.Cursor
	mu         sync.Mutex
	client     *resty.Client
}

func NewJobmanager(cursor *db.Cursor, accrualURL string) *Jobmanager {
	return &Jobmanager{
		AccrualURL: accrualURL,
		Jobs:       make(chan *Job),
		Cursor:     cursor,
		client:     resty.New().SetBaseURL(accrualURL),
	}
}

func (jm *Jobmanager) AskAccrual(url string, number string) (*models.AccrualResponse, int, error) {
	acc := models.AccrualResponse{}
	req := jm.client.R().
		SetResult(&acc).
		SetPathParam("number", number)

	resp, err := req.Get("/api/orders/{number}")
	if err != nil {
		logger.ErrorLog.Printf("Error getting order from accrual: %e", err)
		return nil, 0, err
	}
	logger.InfoLog.Printf("Accrual GET status code: %d", resp.StatusCode())
	if resp.StatusCode() == 429 {
		return nil, resp.StatusCode(), nil
	}
	if resp.StatusCode() == 204 {
		return &models.AccrualResponse{Status: "NEW"}, 204, nil
	}
	return &acc, resp.StatusCode(), nil
}

func (jm *Jobmanager) RunJob(job *Job) error {
	response, statusCode, err := jm.AskAccrual(jm.AccrualURL, job.orderNumber)
	if err != nil {
		return err
	}
	if statusCode == 429 {
		time.Sleep(time.Second)
	}
	for response.Status != "INVALID" && response.Status != "PROCESSED" {
		response, statusCode, err = jm.AskAccrual(jm.AccrualURL, job.orderNumber)
		if err != nil {
			return err
		}
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
	jm.Cursor.UpdateUserBalance(job.username, &models.Balance{
		Current:   response.Accrual,
		Withdrawn: 0.0,
	})
	jm.mu.Unlock()
	logger.InfoLog.Println("Job finished")
	return nil
}

func (jm *Jobmanager) AddJob(orderNumber string, username string) {
	jm.Jobs <- &Job{orderNumber: orderNumber, username: username}
}

func (jm *Jobmanager) ManageJobs(accrualURL string) {
	for job := range jm.Jobs {
		logger.InfoLog.Printf("Running job for order %s", job.orderNumber)
		go jm.RunJob(job)
	}
	close(jm.Jobs)
}
