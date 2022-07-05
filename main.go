package main

import (
	"core-infra/internal/networkingandcontdelivery/vpc"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		//Create my vpc
		vpcErr := vpc.CreateVpc(ctx)

		if vpcErr != nil {
			return vpcErr
		}

		//Create our Subnets withing their respective VPC and AZ
		subnetErr := vpc.CreateSubnet(ctx)

		if subnetErr != nil {
			return subnetErr
		}

		return nil

	})

}
