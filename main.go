package main

import (
	"core-infra/internal/networkingandcontdelivery/vpc"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		//Create my vpc
		vpc.CreateVpc(ctx)

		//Create our Subnets withing their respective VPC and AZ
		vpc.CreateSubnet(ctx)

		return nil

	})

}
