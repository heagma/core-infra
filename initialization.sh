#!/bin/bash

## Install Pulumi
curl -fsSL https://get.pulumi.com | sh

##Create an S3 Bucket to store the remote state
aws s3api create-bucket --bucket <add-bucket-name-here> --region <add-region-here> --create-bucket-configuration LocationConstraint=<add-same-region-here>
