# core-infra
Building Cloud resources with Pulumi and the Go programming language.

## What is it for?
This repo contains the core-infra example of infrastructure created using the Go programming language together with Pulumi IaC.
It creates resources that are needed to have the bare minimum to deploy quickly (day 1) some initial services.

A day 2 operation is to add some more resources/services as per required on the business logic and those should live in a different repository.
Those new deployments can be a new S3 bucket , a new Lambda function or even new services running on EKS but that needs to interact with some other AWS resources.

## What is a bare minimum ?
For different companies the bare minimum is most likely completely different, but many times, specially for companies just starting up in the cloud, an initial setup will require some minimum resources as follow:

+ 3 Accounts in 2 organizations units (OU) for environment separation (dev, prd) and 1 in Master Account.
+ 1 IAM admin account on each cloud account with the Least Privilege to perform tasks required for the initial set up.
  * 2 roles to access the other cloud accounts (role switching).
+ 2-4 Environments (dev, uat, stg, prd). This repo will use only 2 (dev and prd).
+ 2 Regions for prd environment. Just 1 for dev.
+ 3 Availability Zones per region.
+ 1 VPC on each Region:
  * 2 Subnets on each of the 3 AZ (3 Privates and 3 Publics).
  * 1 Route Table for each Subnet.
  * 1 NAT Gateway on each Public Subnet.
  * 1 Internet Gateways for each VPC.
+ 1 EKS Cluster per Region.
  * 3 workers nodes on each cluster.
  * 1 ELB per EKS cluster.
+ 1 EFS on each VPC.
+ 2 S3 Buckets for each VPC.
+ 1 EC2 on each VPC. To act as a bastion/jump host. This does not include those EC2 required by the EKS clusters.
+ 1 RDS on each VPC.


### Pre-requisites
- An AWS key/secret to perform the initial tasks of creating the S3 bucket to act as a backend for our Pulumi states.
- Pulumi
- Go v1.18+

#### Notes:
The approach is done creating all the resources by ourselves to have more control in a deeper level (vpc/subnets/routes/).
Instead of using this approach one could use Pulumi Crosswalk that comes preconfigured with best practice to spin up a whole set of VPC (for example) but leaving some of the fine grained customization to the service itself.