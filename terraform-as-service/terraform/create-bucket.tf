provider "aws" {
  profile = "default"
  region  = "us-east-1"
}

resource "aws_s3_bucket" "arshad-bucket-from-tf-service" {
  bucket = "arshad-bucket-from-tf-service"
  acl    = "private"
}