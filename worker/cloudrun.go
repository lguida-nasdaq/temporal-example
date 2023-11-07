package main

import (
	"log"

	"github.com/metriodev/temporal/cloudrun"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// New worker for the deplpy queue
	w := worker.New(c, cloudrun.DeployQueueName, worker.Options{})
	// clourRun Activities with a GCP client
	cloudRunActivities := cloudrun.NewActivities(new(cloudrun.GoogleCloudRunServiceFactory))
	// Register the cloudRun activities
	w.RegisterActivity(cloudRunActivities)
	// Register the deploy workflow
	w.RegisterWorkflow(cloudrun.DeployClientWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}