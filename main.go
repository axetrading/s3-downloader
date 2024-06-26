package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	s3Uri := flag.String("s3uri", "", "The S3 URI of the item to download (s3://bucket-name/path/to/item)")
	outputFolder := flag.String("output", "", "Optional: The folder to save the downloaded item")
	flag.Parse()

	if *s3Uri == "" {
		fmt.Println("Missing required argument: s3uri")
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
	filename := filepath.Base(item) // Extract file name from the item path

	outputPath := filename
	if *outputFolder != "" {
		outputPath = filepath.Join(*outputFolder, filename)

		// Create the output directory if it doesn't exist
		if err := os.MkdirAll(*outputFolder, os.ModePerm); err != nil {
			fmt.Printf("Failed to create output directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Determine the region to use for the S3 client
	region := "eu-west-2" // Default region
	if customRegion := os.Getenv("S3_CONFIG_BUCKET"); customRegion != "" {
		region = customRegion
	}

	// Load AWS configuration with dynamic region setting
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &item,
	})
	if err != nil {
		fmt.Printf("Failed to download item: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Copy the response body to the file
	if _, err := io.Copy(file, resp.Body); err != nil {
		fmt.Printf("Failed to write to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully downloaded %q to %q\n", *s3Uri, outputPath)
}
