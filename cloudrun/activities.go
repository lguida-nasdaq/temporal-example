package cloudrun

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/temporal"
	"google.golang.org/api/run/v1"
	runv1 "google.golang.org/api/run/v1"
	runv2 "google.golang.org/api/run/v2"
)

type ClourRunServicFactory interface {
	GetV1Client(ctx context.Context, region string) (*runv1.APIService, error)
	GetV2Client(ctx context.Context, region string) (*runv2.Service, error)
}

type Activities struct {
	runFactory ClourRunServicFactory
}

func NewActivities(runFactory ClourRunServicFactory) *Activities {
	return &Activities{runFactory: runFactory}
}

func (a *Activities) UpdateJob(ctx context.Context, job CloudRun) error {
	runClient, err := a.runFactory.GetV1Client(ctx, job.Region)
	if err != nil {
		return temporal.NewApplicationError(err.Error(), "Google Error")
	}
	spec := &run.Job{
		ApiVersion: "run.googleapis.com/v1",
		Kind: "Job",
		Metadata: &run.ObjectMeta{
			Name: job.Name,
		},
		Spec: &run.JobSpec{
			Template: &run.ExecutionTemplateSpec{
				Spec: &run.ExecutionSpec{
					Template: &run.TaskTemplateSpec{
						Spec: &run.TaskSpec{
							Containers: []*run.Container{
								{
									Image: job.Image,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = runClient.Namespaces.Jobs.ReplaceJob(job.FullName(), spec).Do()
	return err
}

func (a *Activities) StartJobExecution(ctx context.Context, job CloudRun) (*StartJobExecutionResponse, error) {
	executionReponse := &StartJobExecutionResponse{}
	runClient, err := a.runFactory.GetV2Client(ctx, job.Region)
	if err != nil {
		return executionReponse, temporal.NewApplicationError(err.Error(), "Google Error")
	}
	operation, err := runClient.Projects.Locations.Jobs.Run(job.FullPath(), &runv2.GoogleCloudRunV2RunJobRequest{}).Do()
	if err != nil {
		return executionReponse, fmt.Errorf("Failed to start job execution: %w", err)
	}
	// Return the operation name so we can fetch it.
	executionReponse.OperationId = operation.Name
	return executionReponse, nil
}

// Reads the job executin operation current state. Return a retryable error if the
// operation is not finished, or a non-retryable one if the operation has an error.
func (a *Activities) WaitForExecutinToFinish(ctx context.Context, operationId, region string) error {
	runClient, err := a.runFactory.GetV2Client(ctx, region)
	if err != nil {
		return temporal.NewApplicationError(err.Error(), "Google Error")
	}
	operation, err := runClient.Projects.Locations.Operations.Get(operationId).Do()
	if err != nil {
		return err
	}

	if !operation.Done {
		// We raise an error if not done, so the retry policy will try
		// to get the operation again
		return temporal.NewApplicationError("Operation not done yet", jobIsNotFinished.String())
	}

	if operation.Error != nil {
		// The operation finished with an error, this should stop the retry policy
		 return temporal.NewNonRetryableApplicationError(
			"Failed to run cloud run job",
			jobExecutionFailedError.String(),
			fmt.Errorf("[%d]: %s", operation.Error.Code, operation.Error.Message),
		 )
	}

	// Operation sucessuly finished
	return nil
}