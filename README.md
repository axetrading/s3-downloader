# S3 Downloader

This utility is a simple command-line tool written in Go for downloading files from Amazon S3. It takes an S3 URI and a local file path as input and saves the S3 object to the specified local path.

## Requirements

- Go 1.x or later.
- AWS credentials configured (either via environment variables, AWS credentials file, or IAM roles if running on AWS services like EC2).

## Installation

To install this utility, clone the repository and build the binary:

```bash
go build -o s3downloader
```

## Usage
Run the tool with the following command:
```bash
./s3downloader -s3uri s3://bucket-name/path/to/item -output /path/to/local/file
```
## Arguments
- s3uri: The S3 URI of the file to download (format: s3://bucket-name/path/to/item).
- output: The local path where the downloaded file should be saved.

## Example
```bash
./s3downloader -s3uri s3://mybucket/myfolder/myfile.txt -output ./myfile.txt
```
This command will download myfile.txt from s3://mybucket/myfolder/ in S3 to the current directory.

## Licence 
This project is licensed under the MIT License - see the LICENSE file for details