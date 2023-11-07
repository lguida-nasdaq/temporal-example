package cloudrun

import (
	"context"
	"fmt"

	"google.golang.org/api/option"
	runv1 "google.golang.org/api/run/v1"
	runv2 "google.golang.org/api/run/v2"
)

type ErrorType string

func (e ErrorType) String() string {
	return string(e)
}

const (
	jobExecutionFailedError ErrorType = "jobExecutionFailed"
	jobIsNotFinished        ErrorType = "jobIsNotFinished"
)

type CloudRun struct {
	Name 	  string
	ProjectId string
	Image     string
	Region    string
}

// func (c CloudRun) FullName() string {
// 	return
// }

// func (c CloudRun) FullPath() string {
// 	return fmt.Sprintf("projects/%s/locations/%s/jobs/%s", c.ProjectId, c.Region, c.Name)
// }

type StartJobExecutionResponse struct {
	OperationId string
}


type GoogleCloudRunServiceFactory struct {}

func (f *GoogleCloudRunServiceFactory) GetV1Client(ctx context.Context, region string) (*runv1.APIService, error) {
	return runv1.NewService(
		ctx,
		option.WithEndpoint(fmt.Sprintf("https://%s-run.googleapis.com", region)),
	)
}

func (f *GoogleCloudRunServiceFactory) GetV2Client(ctx context.Context, region string) (*runv2.Service, error) {
	return runv2.NewService(
		ctx,
		option.WithEndpoint(fmt.Sprintf("https://%s-run.googleapis.com", region)),
	)
}