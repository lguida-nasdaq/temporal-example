package cloudrun

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Deploy client workflow
func DeployClientWorkflow(ctx workflow.Context, clientName, image string) error {
	job := CloudRun{
		Name: fmt.Sprintf("%s-migration", clientName),
		ProjectId: "pissenlit",
		Region: "us-central1",
		// For this test is a different image but it should be the same as the service
		Image: "us-docker.pkg.dev/pissenlit/metrio/sleeper:latest",
	}
	// service := CloudRun{
	// 	Name: clientName,
	// 	ProjectId: "pissenlit",
	// 	Region: "us-central1",
	// 	Image: image,
	// }

	err := updateCloudRunJob(ctx, job)
	if err != nil {
		return err
	}
	return nil
}

func updateCloudRunJob(ctx workflow.Context, job CloudRun) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		// Activity timeout
		StartToCloseTimeout: time.Minute * 5,
	})
	var a *Activities
	future := workflow.ExecuteActivity(ctx, a.UpdateJob, job)

	// If we had a struct return type, we would pass the pointer
	// instead of nil
	return future.Get(ctx, nil)
}

