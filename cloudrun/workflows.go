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

	operationId, err := startJobExecution(ctx, job)
	if err != nil {
		return err
	}

	err = waitForExecutinToFinish(ctx, operationId, job.Region)
	if err != nil {
		return err
	}

	return nil
}

func updateCloudRunJob(ctx workflow.Context, job CloudRun) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	})
	var a *Activities
	future := workflow.ExecuteActivity(ctx, a.UpdateJob, job)
	return future.Get(ctx, nil)
}

func startJobExecution(ctx workflow.Context, job CloudRun) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	})
	var a *Activities
	future := workflow.ExecuteActivity(ctx, a.StartJobExecution, job)
	var response StartJobExecutionResponse
	err := future.Get(ctx, &response)
	return response.OperationId, err
}

func waitForExecutinToFinish(ctx workflow.Context, operationId, region string) error {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 10,
		// In a real case, setup the correct values for the
		// RetryPolicy
		// RetryPolicy: &internal.RetryPolicy{
		// },
	})

	var a *Activities
	future := workflow.ExecuteActivity(ctx, a.WaitForExecutinToFinish, operationId, region)
	return future.Get(ctx, nil)
}
