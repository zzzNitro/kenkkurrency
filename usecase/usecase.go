package usecase

import (
	"net/http"

	"github.com/zzzNitro/kenkkurrency/repository"
)

type UseCase struct {
	Repo repository.Repository
}

func NewUseCase(repo repository.Repository) *UseCase {
	return &UseCase{Repo: repo}
}

// Handlers for different concurrency models
func (uc *UseCase) HandleControl(w http.ResponseWriter, r *http.Request) {
	// Sequential processing logic
}

func (uc *UseCase) HandleWaitGroup(w http.ResponseWriter, r *http.Request) {
	// Concurrency with WaitGroup
}

func (uc *UseCase) HandleChannels(w http.ResponseWriter, r *http.Request) {
	// Concurrency with Channels
}

func (uc *UseCase) HandleMutexes(w http.ResponseWriter, r *http.Request) {
	// Concurrency with Mutexes
}
