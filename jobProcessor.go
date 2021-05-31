package main

type Job struct {
	ID      string   `json:"id,omitempty"`
	Command string   `json:"command,omitempty"`
	Message string   `json:"message,omitempty"`
	Status  string   `json:"status,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// Handle incoming jobs
func handleJob(job *Job, client *Client) {
	if len(job.Command) > 0 {
		l("["+job.ID+"] Received Command: "+job.Command, false, false)
		// Handle different commands by calling a function
		switch job.Command {
		case "ping":
			pong(job, client)
		default:
			unknownCommand(client)
		}
	}
	if len(job.Status) > 0 {
		l("["+job.ID+"] Status: "+job.Status, false, false)
	}
	if len(job.Message) > 0 {
		l("["+job.ID+"] Message: "+job.Message, false, false)
	}
}

func pong(job *Job, client *Client) {
	pong := Job{Message: "PONG", ID: job.ID}
	client.data <- pong
}

func unknownCommand(client *Client) {
	err := Job{Message: "Unknown Command", Status: "ERROR"}
	client.data <- err
}
