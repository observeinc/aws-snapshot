package cloudwatchlogs

import (
	"context"
	"fmt"

	"github.com/observeinc/aws-snapshot/pkg/api"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ListBucketsOutput struct {
	*s3.Bucket
	*s3.GetBucketTaggingOutput
	*s3.GetBucketLocationOutput
	*s3.GetBucketPolicyOutput
	*s3.GetBucketAclOutput
}

func (o *ListBucketsOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		Data: o,
	})
	return
}

type ListBuckets struct {
	API
	Region *string
}

type CountBucketsOutput struct {
	Count int `json:"Count"`
}

func (o *CountBucketsOutput) Records() (records []*api.Record) {
	records = append(records, &api.Record{
		Data: o,
	})
	return
}

var _ api.RequestBuilder = &ListBuckets{}

// New implements api.RequestBuilder
func (fn *ListBuckets) New(name string, config interface{}) ([]api.Request, error) {
	var input s3.ListBucketsInput
	// get our config into input
	if err := api.DecodeConfig(config, &input); err != nil {
		return nil, err
	}
	// define the call function
	call := func(ctx context.Context, ch chan<- *api.Record) error {
		output, err := fn.ListBucketsWithContext(ctx, &input)
		if err != nil {
			return err
		}
		r, _ := ctx.Value("runner_config").(api.Runner)

		if r.Stats {
			countBucketOutput := &CountBucketsOutput{Count: len(output.Buckets)}
			// Send it
			if err := api.SendRecords(ctx, ch, name, countBucketOutput); err != nil {
				return err
			}
		} else {
			// for each bucket
			for _, b := range output.Buckets {
				/// get bucket locations
				locationOutput, err := fn.GetBucketLocationWithContext(ctx, &s3.GetBucketLocationInput{
					Bucket: b.Name,
				})
				if err != nil {
					// look for access denied and NoSuchBucket
					aerr, ok := err.(awserr.Error)
					switch {
					case ok && aerr.Code() == "AccessDenied" || aerr.Code() == "NoSuchBucket":
						// Eat AccessDenied
						continue
					default:
						return fmt.Errorf("failed to get bucket location for %s: %w", *b.Name, err)
					}
				}

				listBucketsOutput := &ListBucketsOutput{
					Bucket:                  b,
					GetBucketLocationOutput: locationOutput,
				}

				// wonderful legacy quirks
				if locationOutput.LocationConstraint == nil {
					locationOutput.LocationConstraint = aws.String("us-east-1")
				}

				if *locationOutput.LocationConstraint == *fn.Region {
					// We can only request tags from bucket in our region.
					listBucketsOutput.GetBucketTaggingOutput, err = fn.GetBucketTaggingWithContext(ctx, &s3.GetBucketTaggingInput{
						Bucket: b.Name,
					})
					if err != nil {
						aerr, ok := err.(awserr.Error)
						switch {
						case ok && aerr.Code() == "NoSuchTagSet":
							// null tagset, keep going
						case ok && aerr.Code() == "AccessDenied":
							continue
						default:
							return fmt.Errorf("failed to get bucket tags for %s: %w", *b.Name, err)
						}
					}
				}

				if *locationOutput.LocationConstraint == *fn.Region {
					// We can only request policies from bucket in our region.
					listBucketsOutput.GetBucketPolicyOutput, err = fn.GetBucketPolicyWithContext(ctx, &s3.GetBucketPolicyInput{
						Bucket: b.Name,
					})
					if err != nil {
						aerr, ok := err.(awserr.Error)
						switch {
						case ok && aerr.Code() == "NoSuchBucketPolicy":
							continue
						case ok && aerr.Code() == "AccessDenied":
							continue
						default:
							return fmt.Errorf("failed to get bucket policy for %s: %w", *b.Name, err)
						}
					}
				}

				if *locationOutput.LocationConstraint == *fn.Region {
					// We can only request acls from bucket in our region.
					listBucketsOutput.GetBucketAclOutput, err = fn.GetBucketAclWithContext(ctx, &s3.GetBucketAclInput{
						Bucket: b.Name,
					})
					if err != nil {
						aerr, ok := err.(awserr.Error)
						switch {
						case ok && aerr.Code() == "AccessDenied":
							continue
						default:
							return fmt.Errorf("failed to get bucket acl for %s: %w", *b.Name, err)
						}
					}
				}
				if err := api.SendRecords(ctx, ch, name, listBucketsOutput); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return []api.Request{call}, nil
}
