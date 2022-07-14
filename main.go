package main

import (
	"core-infra/config"
	"core-infra/internal/networkingandcontdelivery/vpc"
	"fmt"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		//load 'global' configurations to pass into each resource call
		cfg := config.NewConfig(ctx)

		//Create my vpc
		VPCs, err := vpc.CreateVpc(ctx, cfg)

		if err != nil {
			return fmt.Errorf("main: resource creation failed %[1]w", err)

		}

		//Create our Subnets within their respective VPC and AZ. We also create their respective NAT gateways
		subnets, err := vpc.CreateSubnet(ctx, cfg, VPCs)

		if err != nil {
			return fmt.Errorf("main: resource creation failed %[1]w", err)
		}

		//Create Nat gateways inside our subnets
		_, err = vpc.AddNat(ctx, cfg, subnets)

		if err != nil {
			return fmt.Errorf("main: resource creation failed %[1]w", err)
		}

		//Attach route table to subnets
		//_, err := vpc.

		return nil

	})

}
