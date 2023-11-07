package cloudrun

import (
	"log"

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
	w := worker.New(c, DeployQueueName, worker.Options{})
	// clourRun Activities with a GCP client
	cloudRunActivities := NewActivities(new(GoogleCloudRunServiceFactory))
	// Register the cloudRun activities
	w.RegisterActivity(cloudRunActivities)
	// Register the deploy workflow
	w.RegisterWorkflow(DeployClientWorkflow)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}