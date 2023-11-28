package gcs

import (
	"context"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const OpenshiCiBucketName = "origin-ci-test"

type GCSBucket struct {
	// retrieval mechanisms
	bkt *storage.BucketHandle
}

func BucketHandleClient() *GCSBucket {
	d, _ := storage.NewClient(context.Background(), option.WithoutAuthentication())

	return &GCSBucket{
		bkt: d.Bucket(OpenshiCiBucketName),
	}
}

func (b *GCSBucket) GetJobJunitContent(orgName string, repoName string, pullNumber string, jobId string, jobType string, jobName string, junitName string) []byte {
	query := &storage.Query{}

	if jobType == "periodic" {
		query.Prefix = fmt.Sprintf("logs/%s/%s", jobName, jobId) // logs/periodic-ci-redhat-appstudio-infra-deployments-main-appstudio-e2e-tests-periodic
	} else if jobType == "presubmit" {
		query.Prefix = fmt.Sprintf("pr-logs/pull/%s_%s/%s/%s/%s/artifacts", orgName, repoName, pullNumber, jobName, jobId)
	}

	it := b.bkt.Objects(context.Background(), query)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}

		if strings.HasSuffix(obj.Name, junitName) {
			if b.ContentExists(context.Background(), obj.Name) {
				content, _ := b.GetContent(context.Background(), obj.Name)

				return content
			}
		}
	}
	return nil
}

func (b *GCSBucket) ContentExists(ctx context.Context, path string) bool {
	// Get an Object handle for the path
	obj := b.bkt.Object(path)

	// if we can get the attrs then presume the object exists
	// otherwise presume it doesn't
	_, err := obj.Attrs(ctx)
	return err == nil
}

func (b *GCSBucket) GetContent(ctx context.Context, path string) ([]byte, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("missing path to GCS content for jobrun")
	}

	// Get an Object handle for the path
	obj := b.bkt.Object(path)

	// use the object attributes to try to get the latest generation to try to retrieve the data without getting a cached
	// version of data that does not match the latest content.  I don't know if this will work, but in the easy case
	// it doesn't seem to fail.
	objAttrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading GCS attributes for jobrun: %w", err)
	}
	obj = obj.Generation(objAttrs.Generation)

	// Get an io.Reader for the object.
	gcsReader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("error reading GCS content for jobrun: %w", err)
	}

	defer gcsReader.Close()

	return io.ReadAll(gcsReader)
}
