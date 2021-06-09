package main

import (
	"fmt"
	"time"
)

// TODO Rework job processor to include outgoing requests and meausre job completion time

type Job struct {
	ID      string   `json:"id,omitempty"`
	Command string   `json:"command,omitempty"`
	Message string   `json:"message,omitempty"`
	Status  string   `json:"status,omitempty"`
	Args    []string `json:"args,omitempty"`
	Started int64    `json:"started,omitempty"`
	Client  *Client
}

func (job *Job) reset() {
	job.Args = nil
	job.Command = ""
	job.ID = ""
	job.Message = ""
	job.Status = ""
}

func (job *Job) new() {
	job.ID = getID()
	job.Status = "new"
	job.Started = time.Now().UnixNano()
	job.send()
}

func (job *Job) send() {
	jobs[job.ID] = *job
	job.Client.data <- *job
}

func (job *Job) acknowledge() {
	job.Status = "acknowledged"
	job.send()
}

func (job *Job) done() {
	l("["+job.ID+"] Job has been completed in "+time.Since(time.Unix(0, job.Started)).String(), false, true)
	delete(jobs, job.ID)
}

var jobs = make(map[string]Job)

// Handle incoming jobs
func handleJob(job Job, client *Client) {
	if len(jobs[job.ID].Command) == 0 {
		l("["+job.ID+"] New job: "+job.Command, false, false)
	}
	job.Client = client
	jobs[job.ID] = job
	if job.Status == "acknowledged" {
		l("["+job.ID+"] Has been acknowledged by server", false, true)
	} else if len(job.Command) > 0 {
		// Handle different commands by calling a function
		switch job.Command {
		case "ping":
			pong(job.ID)
		// If command is not found an error message will be sent back to the sender
		default:
			unknownCommand(job.ID)
		}
	}
}

func pong(id string) {
	job := jobs[id]
	if job.Status == "completed" {
		l("["+job.ID+"] "+job.Message, false, true)
		job.done()
	} else if job.Status == "new" {
		job.Message = "pong"
		job.Status = "completed"
		job.send()
		job.done()
	}
}

func unknownCommand(id string) {
	job := jobs[id]
	l("["+job.ID+"] Unknown Command '"+job.Command+"'", false, true)
	if job.Status != "ERROR" {
		job.Message = "Unknown Command"
		job.Status = "ERROR"
		job.send()
	}
	job.done()
}

func listJobs() {
	for id, job := range jobs {
		fmt.Println("[" + id + "] status: " + job.Status + "| Command: " + job.Command)
	}
}
