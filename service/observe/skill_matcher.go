package observe

import (
	"context"
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/models"
	service "hireforwork-server/service/modules"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type JobNotification struct {
	Career interface{}
	Job    *models.Jobs
}

// SkillMatcherObserver implements the Observer interface
type SkillMatcherObserver struct {
	db               *db.DB
	careerCollection *mongo.Collection
	notificationChan chan JobNotification
	wg               sync.WaitGroup
}

func NewSkillMatcherObserver(db *db.DB) *SkillMatcherObserver {
	observer := &SkillMatcherObserver{
		db:               db,
		careerCollection: db.GetCollection("Career"),
		notificationChan: make(chan JobNotification, 100), // Buffer size of 100
	}

	// Start the notification worker
	observer.startNotificationWorker()
	return observer
}

func (s *SkillMatcherObserver) startNotificationWorker() {
	// Start multiple workers for parallel processing
	for i := 0; i < 5; i++ {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for notification := range s.notificationChan {
				s.processNotification(notification)
			}
		}()
	}
}

func (s *SkillMatcherObserver) OnJobPosted(job *models.Jobs) {
	fmt.Println("Checking for matching careers...")
	fmt.Println("Job Requirements:", job.JobRequirement)

	// Create a filter to match careers with at least one matching skill
	filter := bson.M{
		"isDeleted": false, // Only get active accounts
		"profile.skills": bson.M{
			"$elemMatch": bson.M{
				"$in": job.JobRequirement,
			},
		},
	}
	// Get matching careers from the database
	cursor, err := s.careerCollection.Find(context.Background(), filter, nil)
	if err != nil {
		fmt.Printf("Error fetching careers: %v\n", err)
		return
	}
	defer cursor.Close(context.Background())

	// Decode all careers into a slice
	var careers []map[string]interface{}
	if err = cursor.All(context.Background(), &careers); err != nil {
		fmt.Printf("Error decoding careers: %v\n", err)
		return
	}

	fmt.Printf("Found %d careers with matching skills\n", len(careers))

	// Send notifications asynchronously
	for _, career := range careers {
		// Send to channel instead of direct processing
		s.notificationChan <- JobNotification{
			Career: career,
			Job:    job,
		}
	}
}

func (s *SkillMatcherObserver) processNotification(notification JobNotification) {
	career := notification.Career
	job := notification.Job

	careerMap := career.(map[string]interface{})

	// Get career email
	careerEmail, ok := careerMap["careerEmail"].(string)
	if !ok {
		return
	}

	// Get career name
	firstName, _ := careerMap["careerFirstName"].(string)
	lastName, _ := careerMap["lastName"].(string)
	name := firstName + " " + lastName

	// Add this at the top of processNotification function
	hostURL := os.Getenv("HOST_URL") // Match the exact env var name
	if hostURL == "" {
		hostURL = "http://localhost:8080" // Default fallback if env var not set
		fmt.Println("Warning: HOSTING_URL not set, using default:", hostURL)
	}

	// Create email subject
	subject := fmt.Sprintf("Job Alert: New %s position matching your skills", job.JobTitle)
	// Create email body using HTML template
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2557a7;">Job Alert: New Position Match!</h2>
        
        <p>Dear %s,</p>
        
        <p>We found a job that matches your skills! Based on your profile, you have the skills this position requires.</p>
        
        <div style="background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 20px 0;">
            <h3 style="color: #2557a7; margin-top: 0;">%s</h3>
        </div>
        
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s/jobs/%s" style="background-color: #2557a7; color: white; padding: 12px 25px; text-decoration: none; border-radius: 5px; font-weight: bold;">View Job & Apply</a>
        </div>
        
        <p style="color: #666; font-size: 0.9em;">
            Best regards,<br>
            HireForWork Team
        </p>
    </div>
</body>
</html>`,
		name,
		job.JobTitle,
		hostURL,
		job.Id.Hex(),
	)

	// Send email using the service
	if err := service.SendRecommendationJob(careerEmail, subject, body); err != nil {
		fmt.Printf("Error sending email to %s: %v\n", careerEmail, err)
		return
	}

	fmt.Printf("Sent job notification to %s (%s)\n", name, careerEmail)
}

// Cleanup method to properly shut down the observer
func (s *SkillMatcherObserver) Shutdown() {
	close(s.notificationChan)
	s.wg.Wait()
}
