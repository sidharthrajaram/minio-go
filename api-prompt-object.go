package minio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"io"
	"net/http"
)

// PromptObject wrapper function that accepts a request context
func (c *Client) PromptObject(ctx context.Context, bucketName, objectName string, prompt string, opts PromptObjectOptions) (string, error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return "FAIL", ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Code:       "InvalidBucketName",
			Message:    err.Error(),
		}
	}
	if err := s3utils.CheckValidObjectName(objectName); err != nil {
		return "FAIL", ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Code:       "XMinioInvalidObjectName",
			Message:    err.Error(),
		}
	}

	opts.AddReqParam("lambdaArn", opts.lambdaArn)

	opts.SetHeader("Content-Type", "application/json")

	opts.AddKwarg("prompt", prompt)
	promptReqBytes, err := json.Marshal(opts.kwargs)

	// Execute POST on bucket/object.
	resp, err := c.executeMethod(ctx, http.MethodPost, requestMetadata{
		bucketName:       bucketName,
		objectName:       objectName,
		queryValues:      opts.toQueryValues(),
		customHeader:     opts.Header(),
		contentSHA256Hex: sum256Hex(promptReqBytes),
		contentBody:      bytes.NewReader(promptReqBytes),
		contentLength:    int64(len(promptReqBytes)),
	})
	if err != nil {
		return "FAIL", err
	}
	// Return string that is returned as HttpResponse
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body")
			return
		}
	}(resp.Body)

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Convert to string and return
	responseString := string(bodyBytes)
	return responseString, nil
}
