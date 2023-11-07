package cloudrun

import (
	"context"

	"go.temporal.io/sdk/temporal"
	"google.golang.org/api/run/v1"
)

type ClourRunServicFactory interface {
	GetClient(ctx context.Context, region string) (*run.APIService, error)
}

type Activities struct {
	runFactory ClourRunServicFactory
}

func NewActivities(runFactory ClourRunServicFactory) *Activities {
	return &Activities{runFactory: runFactory}
}

func (a *Activities) UpdateJob(ctx context.Context, job CloudRun) error {
	runClient, err := a.runFactory.GetClient(ctx, job.Region)
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