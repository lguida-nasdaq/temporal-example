package main

import (
	"context"
	"fmt"
	"log"

	"github.com/metriodev/temporal/cloudrun"
	"go.temporal.io/sdk/client"
)

func main() {
	// create a new temporal client
	// set up the worker
	c, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	clientName := "fake-client"
	image := "gcr.io/google-samples/hello-app:2.0"

	ctx := context.Background()

	workflowRun, err := c.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID: "deploy_fake_client_git_sha",
		TaskQueue: cloudrun.DeployQueueName,
	}, cloudrun.DeployClientWorkflow, clientName, image)

	if err != nil {
		log.Fatal("Error starting workflow: ", err.Error())
	}

	// Wait for workflow execution
	fmt.Println("Waiting for worklow completion")
	err = workflowRun.Get(ctx, nil)
	if err != nil {
		log.Fatal("Error on workflow: ", err.Error())
	}

	fmt.Printf("Sucessufly deployed %s with %s\n", clientName, image)
}

