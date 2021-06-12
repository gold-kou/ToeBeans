package gcp

import (
	"context"
	"os"

	"github.com/pkg/errors"

	vision "cloud.google.com/go/vision/apiv1"
)

// DetectLabels gets labels from the Vision API for an image at the given file path.
func DetectLabels(file string) (labels []string, err error) {
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return
	}
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return
	}

	if len(annotations) == 0 {
		return labels, errors.New("No labels found")
	} else {
		for _, annotation := range annotations {
			labels = append(labels, annotation.Description)
		}
	}

	return labels, nil
}
