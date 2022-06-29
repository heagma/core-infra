# edge-architecture
Building multiple stacks with Pulumi and Go programming language

## What is it for?
This repo contains an example of infrastructure created using the Go programming language.

It includes the "core-infra" resources that are needed to have the bare minimum to deploy right away (day 1) some initial services.

The approach is to have a core-infra folder with all the required resources but have also a separated repository for all the new components that are added on day 2. 
Those new resources can be a new S3 bucket , a new Lambda function or even a new service running on EKS but that needs to interact with some other AWS resources. 

## What is a bare minimum ?
For different companies the bare minimum is most likely completely different, but many times, specially for start ups, an initial setup will require:

+ 3 Cloud Accounts. Separated in 2 organizations units (OU) for environemnt separation (dev, prd) and 1 Master Account.
+ 2-4 Environments (dev, uat, stg, prd). This repo will use only dev and prd.
+ 2 Regions for production environment.
+ 2 AZ per Region.

For testing environments requirements are usually less but for Production environment will have at minimum:
+ 2 VPC. 1 on each Region:
    * 4 Subnets (2 Privates and 2 Publics). 2 on each AZ.
    * 4 Route Tables. 1 for each Subnet.
    * 2 NAT Gateways. 1 on each Public Subnet.
    * 2 Internet Gateways. 1 for each VPC.
+ 2 EKS Cluster. 1 per Region.
    * 3 workers nodes on each cluster.
    * 1 ELB per EKS cluster.
+ 2 EFS. 1 on each VPC. 
+ 4 S3 Buckets. 2 for each VPC.
+ 2 EC2. 1 on each VPC. Usually to act as a bastion/jump host. This does not include those EC2 required by the EKS clusters. 
+ 1 IAM admin account on each cloud account with Least Privilege to perform tasks required to initial set up.
	* 2 roles to access the other cloud accounts (with switching roles). 

### Pre-requisites
- An AWS key/secret to perform the initial tasks of creating the S3 bucket to act as a backend for our Pulumi states.
- Pulumi
- Go v1.18+
	
#### Notes: 
The approach is done creating all the resources by ourselves to have more control in a more deeper level. Instead of using this approach one could use Pulumi Croswalk that comes preconfigured with best practice to spin up a whole set of VPC (for example) but leaving some of the fine grained customization to the service itself.


