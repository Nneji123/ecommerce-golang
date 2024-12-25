package emailsequence

import (
	"fmt"
	"time"
)

func (executor SequenceExecutor) Run(sequence Sequence) {
	for _, step := range sequence.Steps {
		if step.Email != nil {
			executor.SendEmail(*step.Email)
		} else if step.Trigger != nil {
			executor.ProcessTrigger(*step.Trigger)
		} else if step.Goal != nil {
			executor.ProcessGoal(*step.Goal)
		}
	}
}

func (executor SequenceExecutor) SendEmail(email Email) {
	// Implement logic to send the email
	fmt.Println("Sending email:", email.Subject)
}

func (executor SequenceExecutor) ProcessTrigger(trigger Trigger) {
	time.Sleep(5 * time.Second)
	fmt.Println("Sleeping for 5 seconds")
	fmt.Println("Processing trigger:", string(trigger.EventType))
}

func (executor SequenceExecutor) ProcessGoal(goal Goal) {
	// Implement logic to process the goal
	fmt.Println("Processing goal:", goal.Description)
}

func main() {
	// Testing the SequenceExecutor
	sequence := Sequence{
		Name: "Test Sequence",
		Steps: []SequenceStep{
			{
				Email: &Email{
					Subject: "Test Subject",
					Body:    "Test Body",
				},
			},
			{
				Trigger: &Trigger{
					EventType: OpenedEmail,
					Condition: Yes,
				},
			},
			{
				Goal: &Goal{
					Description: "Test Goal",
				},
			},
		},
	}

	executor := SequenceExecutor{}
	executor.Run(sequence)
}
