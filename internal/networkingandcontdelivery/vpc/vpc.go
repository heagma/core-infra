package vpc

import (
	"errors"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	vpcLogicalNames []string
	vpcOutputNames  []string
)

//CreateVpc function for the Vpc class that wraps the Pulumi function that creates a new VPC resource
func CreateVpc() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//Load configurations
		cd := loadConfig(ctx)

		//Assign configurations related only to vpc
		configVpcCidr := cd.VpcCidr
		configVpcNames := cd.VpcNames
		configEnv := cd.Env
		configRegion := cd.RegionAlias

		if len(configVpcCidr) < len(configVpcNames) || len(configVpcCidr) > len(configVpcNames) {
			return errors.New("mismatch between the amount of vpc and the cidr configured")
		}

		for index, vpcName := range configVpcNames {
			//Define the names of each vpc resource. This is for logical resources (get id, arn etc.), physical names have a
			//random suffix and also Create names (adding "-out" at the end of the vpc names) for ids to later be used as
			//keys when exporting the resource
			vpcLogicalNames = append(vpcLogicalNames, formatName(configEnv, configRegion, vpcName))
			vpcOutputNames = append(vpcOutputNames, formatName(vpcName, "out"))

			//Create our VPC resources
			vpc, err := ec2.NewVPC(ctx, vpcName, &ec2.VPCArgs{
				CidrBlock:          pulumi.String(configVpcCidr[index]),
				EnableDnsHostnames: pulumi.Bool(true),
				Tags: ec2.VPCTagArrayInput(
					ec2.VPCTagArray{
						ec2.VPCTagArgs{Key: pulumi.String(initialTags.Name), Value: pulumi.String(vpcName)},
						ec2.VPCTagArgs{Key: pulumi.String(initialTags.Env), Value: pulumi.String(configEnv)},
					},
				),
			})
			if err != nil {
				return err
			}

			ctx.Export(fmt.Sprintf(vpcOutputNames[index], vpcName), vpc.ID())
		}
		return nil
	})
}
