package main

import (
	"core-infra/internal/networkingandcontdelivery/vpc"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		//Create my vpc
		err, VPCs := vpc.CreateVpc(ctx)

		if err != nil {
			return err

		}

		//Create our Subnets within their respective VPC and AZ
		err, _ = vpc.CreateSubnet(ctx, VPCs)

		if err != nil {
			return err
		}

		return nil

	})

}
