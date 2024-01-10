package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	s3Uri := flag.String("s3uri", "", "The S3 URI of the item to download (s3://bucket-name/path/to/item)")
	outputPath := flag.String("output", "", "The output path for the downloaded item")
	flag.Parse()

	if *s3Uri == "" || *outputPath == "" {
		fmt.Println("Missing required arguments: s3uri or output")
		os.Exit(1)
	}

	// Parse the S3 URI
	parsedUrl, err := url.Parse(*s3Uri)
	if err != nil {
		fmt.Printf("Invalid S3 URI: %v\n", err)
		os.Exit(1)
	}

	if parsedUrl.Scheme != "s3" {
		fmt.Println("Invalid S3 URI: Must start with s3://")
		os.Exit(1)
	}

	bucket := parsedUrl.Host
	item := strings.TrimPrefix(parsedUrl.Path, "/")

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Create(*outputPath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &item,
	})
	if err != nil {
		fmt.Printf("Failed to download item: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully downloaded %q to %q\n", *s3Uri, *outputPath)
}
