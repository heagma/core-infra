package vpc

import (
	"core-infra/config"
	"errors"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	cec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	VPC             []*ec2.VPC
	vpcLogicalNames []string
	vpcOutputNames  []string
)

//CreateVpc function for the Vpc class that wraps the Pulumi function that creates a new VPC resource
func CreateVpc(ctx *pulumi.Context) error {

	//Load configurations
	cd := config.NewConfig(ctx)

	configEnv := cd.Env
	configRegion := cd.RegionAlias

	//Assign configurations related only to vpc
	configVpcCidr := cd.VpcCidr
	configVpcNames := cd.VpcNames

	if len(configVpcCidr) < len(configVpcNames) || len(configVpcCidr) > len(configVpcNames) {
		return errors.New("mismatch between the amount of vpc and the cidr configured")
	}

	for _, vpcName := range configVpcNames {
		//Define the names of each vpc resource. This is for logical resources (get id, arn etc.), physical names have a
		//random suffix and also Create names (adding "-out" at the end of the vpc names) for the outputs names
		vpcLogicalNames = append(vpcLogicalNames, config.FormatName(configEnv, configRegion, vpcName))
		vpcOutputNames = append(vpcOutputNames, config.FormatName(vpcName, "out"))

	}

	for index, vpcLogicalName := range vpcLogicalNames {

		//Create our VPC resources
		vpc, err := ec2.NewVPC(ctx, vpcLogicalName, &ec2.VPCArgs{
			CidrBlock:          pulumi.String(configVpcCidr[index]),
			EnableDnsHostnames: pulumi.Bool(true),
			Tags: ec2.VPCTagArrayInput(
				ec2.VPCTagArray{
					ec2.VPCTagArgs{Key: pulumi.String(config.InitialTags.Name), Value: pulumi.String(vpcLogicalName)},
					ec2.VPCTagArgs{Key: pulumi.String(config.InitialTags.Env), Value: pulumi.String(configEnv)},
				},
			),
		})

		VPC = append(VPC, vpc)

		if err != nil {
			return err
		}

		//Create our Internet Gateway for each VPC using the classic aws provider (alias cec2) as aws-native does not
		//provide attachment to vpc yet. We reassign the 'err' variable
		_, err = cec2.NewInternetGateway(ctx, fmt.Sprintf("from-%[1]s", vpcLogicalName), &cec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
			Tags: pulumi.StringMap{
				config.InitialTags.Name: pulumi.String(fmt.Sprintf("from-%[1]s", vpcLogicalName)),
				config.InitialTags.Env:  pulumi.String(configEnv),
			},
		})

		if err != nil {
			return err
		}

		ctx.Export(vpcOutputNames[index], vpc.ID())
	}

	return nil
}
