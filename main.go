package main

import (
	"core-infra/internal/networkingandcontdelivery/vpc"
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		//Create my vpc
		VPCs, err := vpc.CreateVpc(ctx)

		if err != nil {
			return fmt.Errorf("main: resource creation failed %[1]w", err)

		}

		//Create our Subnets within their respective VPC and AZ. We also create their respective NAT gateways
		_, err = vpc.CreateSubnet(ctx, VPCs)

		if err != nil {
			return fmt.Errorf("main: resource creation failed %[1]w", err)
		}

		//Attach route table to subnets
		//_, err := vpc.

		return nil

	})

}
