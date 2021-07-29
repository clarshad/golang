provider "aws" {
  profile = "default"
  region  = "us-east-1"
}

resource "aws_s3_bucket" "test-arshad-bucket" {
  bucket = "test-arshad-bucket"
  acl    = "private"
}