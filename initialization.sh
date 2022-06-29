#!/bin/bash

## Install Pulumi
curl -fsSL https://get.pulumi.com | sh

##Create an S3 Bucket to store the remote state
aws s3api create-bucket --bucket aws-pulumi-core-infra-state --region ap-east-1 --create-bucket-configuration LocationConstraint=ap-east-1
