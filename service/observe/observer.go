package observe

import (
	"fmt"
	"hireforwork-server/models"
)

// Observer interface defines the contract for all observers
type Observer interface {
	OnJobPosted(job *models.Jobs)
}

// Subject interface defines the contract for the subject being observed
type Subject interface {
	Register(observer Observer)
	Unregister(observer Observer)
	Notify(job *models.Jobs)
}

// JobEventManager implements the Subject interface
type JobEventManager struct {
	observers map[Observer]bool
}

func NewJobEventManager() *JobEventManager {
	return &JobEventManager{
		observers: make(map[Observer]bool),
	}
}

func (jem *JobEventManager) Register(observer Observer) {
	jem.observers[observer] = true
}

func (jem *JobEventManager) Unregister(observer Observer) {
	delete(jem.observers, observer)
}

func (jem *JobEventManager) Notify(job *models.Jobs) {
	fmt.Printf("Notifying %d observers...\n", len(jem.observers))
	for observer := range jem.observers {
		fmt.Println("Calling observer.OnJobPosted...")
		observer.OnJobPosted(job)
	}
}
