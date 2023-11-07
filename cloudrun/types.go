package cloudrun

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

type CloudRun struct {
	Name 	  string
	ProjectId string
	Image     string
	Region    string
}

func (c CloudRun) FullName() string {
	return fmt.Sprintf("namespaces/%s/jobs/%s", c.ProjectId, c.Name)
}

type GoogleCloudRunServiceFactory struct {}

func (f *GoogleCloudRunServiceFactory) GetClient(ctx context.Context, region string) (*run.APIService, error) {
	return run.NewService(
		ctx,
		option.WithEndpoint(fmt.Sprintf("https://%s-run.googleapis.com", region)),
	)
}