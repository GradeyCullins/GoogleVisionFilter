package vision

import (
	"context"
	"log"
	"os"

	vision "cloud.google.com/go/vision/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

type BatchAnnotateResponse map[string]*pb.AnnotateImageResponse

// GetImgAnnotations sends a batch request to the Google Vision API to retrieve tags and likelihood values
// for the URI list in the batch request.
func GetImgAnnotations(contentURIMap map[string]string) (BatchAnnotateResponse, error) {
	if len(contentURIMap) == 0 {
		return nil, nil
	}

	res := make(BatchAnnotateResponse, 0)
	ctx := context.Background()

	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Loop over the files and create annotate requests.
	for uri, path := range contentURIMap {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		image, err := vision.NewImageFromReader(f)
		if err != nil {
			return nil, err
		}

		req := &pb.AnnotateImageRequest{
			Image: image,
			Features: []*pb.Feature{
				{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: 5},
			},
		}

		res[uri], err = client.AnnotateImage(ctx, req)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func passSafeSearch(file string) error {
	// [START init]
	ctx := context.Background()

	// Create the client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	// [END init]

	// [START request]
	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return err
	}

	res, err := client.AnnotateImage(ctx, &pb.AnnotateImageRequest{
		Image: image,
		Features: []*pb.Feature{
			{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: 5},
		},
	})
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}

// findLabels gets labels from the Vision API for an image at the given file path.
func findLabels(file string) ([]string, error) {
	// [START init]
	ctx := context.Background()

	// Create the client.
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	// [END init]

	// [START request]
	// Open the file.
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	image, err := vision.NewImageFromReader(f)
	if err != nil {
		return nil, err
	}

	// Perform the request.
	annotations, err := client.DetectLabels(ctx, image, nil, 10)
	if err != nil {
		return nil, err
	}
	// [END request]
	// [START transform]
	var labels []string
	for _, annotation := range annotations {
		labels = append(labels, annotation.Description)
	}
	return labels, nil
	// [END transform]
}
