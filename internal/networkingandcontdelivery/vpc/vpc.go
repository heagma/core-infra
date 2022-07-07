package vpc

import (
	"core-infra/config"
	"errors"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	cec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//CreateVpc function creates a resource of type *ec2.VPC wraping the pulumi NewVPC constructor.
func CreateVpc(ctx *pulumi.Context) (error, []*ec2.VPC) {
	var (
		VPC             []*ec2.VPC
		vpcLogicalNames []string
		vpcOutputNames  []string
	)

	//Load configurations
	cd := config.NewConfig(ctx)

	configEnv := cd.Env
	configRegion := cd.RegionAlias

	//Assign configurations related only to vpc
	configVpcCidrs := cd.VpcCidrs
	configVpcNames := cd.VpcNames

	//load initial tags
	initialTags := config.NewInitialTags()

	if len(configVpcCidrs) < len(configVpcNames) || len(configVpcCidrs) > len(configVpcNames) {
		return errors.New("mismatch between the amount of vpc and the cidr configured"), VPC
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
			CidrBlock:          pulumi.String(configVpcCidrs[index]),
			EnableDnsHostnames: pulumi.Bool(true),
			Tags: ec2.VPCTagArrayInput(
				ec2.VPCTagArray{
					ec2.VPCTagArgs{Key: pulumi.String(initialTags.Name), Value: pulumi.String(vpcLogicalName)},
					ec2.VPCTagArgs{Key: pulumi.String(initialTags.Env), Value: pulumi.String(configEnv)},
				},
			),
		})

		if err != nil {
			return err, VPC
		}

		VPC = append(VPC, vpc)

		//Create our Internet Gateway for each VPC using the classic aws provider (alias cec2) as aws-native does not
		//provide attachment to vpc yet. We reassign the 'err' variable
		_, err = cec2.NewInternetGateway(ctx, fmt.Sprintf("%[1]s-igw", vpcLogicalName), &cec2.InternetGatewayArgs{
			VpcId: vpc.ID(),
			Tags: pulumi.StringMap{
				initialTags.Name: pulumi.String(fmt.Sprintf("%[1]s-igw", vpcLogicalName)),
				initialTags.Env:  pulumi.String(configEnv),
			},
		})

		if err != nil {
			return err, VPC
		}

		ctx.Export(vpcOutputNames[index], vpc.ID())
	}

	return nil, VPC
}
