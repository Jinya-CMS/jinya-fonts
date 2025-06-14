package api

import (
	"encoding/json"
	"github.com/bamzi/jobrunner"
	"jinya-fonts/config"
	"jinya-fonts/storage"
	"net/http"
	"time"
)

type status struct {
	JobIsScheduled   bool       `json:"jobIsScheduled"`
	JobNextExecution *time.Time `json:"jobNextExecution"`
	Online           bool       `json:"online"`
	RedisUp          bool       `json:"redisUp"`
	ServingWebsite   bool       `json:"servingWebsite"`
}

func getStatus(w http.ResponseWriter, _ *http.Request) {
	currentStatus := status{
		JobIsScheduled: len(jobrunner.Entries()) > 0,
		Online:         true,
		RedisUp:        storage.CheckRedis(),
		ServingWebsite: config.LoadedConfiguration.ServeWebsite,
	}

	if currentStatus.JobIsScheduled {
		jobrunnerStatus := jobrunner.StatusPage()
		if len(jobrunnerStatus) > 0 {
			currentStatus.JobNextExecution = &jobrunnerStatus[0].Next
		}
	}

	err := json.NewEncoder(w).Encode(currentStatus)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
