package main

type Job struct {
	ID      string   `json:"id,omitempty"`
	Command string   `json:"command,omitempty"`
	Message string   `json:"message,omitempty"`
	Status  string   `json:"status,omitempty"`
	Args    []string `json:"args,omitempty"`
}

func (job *Job) reset() {
	job.Args = nil
	job.Command = ""
	job.ID = ""
	job.Message = ""
	job.Status = ""
}

// Handle incoming jobs
func handleJob(job Job, client *Client) {
	if len(job.Command) > 0 {
		l("["+job.ID+"] Received Command: "+job.Command, false, false)
		// Handle different commands by calling a function
		switch job.Command {
		case "ping":
			go pong(job, client)
		default:
			go unknownCommand(job, client)
		}
	}
	if len(job.Status) > 0 {
		l("["+job.ID+"] Status: "+job.Status, false, false)
	}
	if len(job.Message) > 0 {
		l("["+job.ID+"] Message: "+job.Message, false, false)
	}
}

func pong(job Job, client *Client) {
	pong := Job{Message: "PONG", ID: job.ID}
	client.data <- pong
}

func unknownCommand(job Job, client *Client) {
	err := Job{Message: "Unknown Command", Status: "ERROR", ID: job.ID}
	client.data <- err
}
