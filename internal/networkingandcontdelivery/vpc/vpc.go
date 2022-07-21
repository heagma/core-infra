package vpc

import (
	"core-infra/config"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	cec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//CreateVpc function creates a resource of type *ec2.VPC wrapping the pulumi NewVPC constructor.
func CreateVpc(ctx *pulumi.Context, cfg *config.DataConfig) (*ec2.VPC, error) {
	var (
		VPC            *ec2.VPC
		vpcLogicalName string
		vpcOutputName  string
	)

	//load initial tags
	initialTags := config.NewInitialTags()

	//Define the names of each vpc resource. This is for logical resources (get id, arn etc.), physical names have a
	//random suffix assigned by Pulumi.
	//We also create output names (adding "-out" at the end of the vpc names) for the outputs names
	vpcLogicalName = config.FormatName(cfg.Env, cfg.RegionAlias, cfg.VpcName)
	vpcOutputName = config.FormatName(vpcLogicalName, "out")

	//Create our VPC resources
	vpc, err := ec2.NewVPC(ctx, vpcLogicalName, &ec2.VPCArgs{
		CidrBlock:          pulumi.String(cfg.VpcCidr),
		EnableDnsHostnames: pulumi.Bool(true),
		Tags: ec2.VPCTagArrayInput(
			ec2.VPCTagArray{
				ec2.VPCTagArgs{Key: pulumi.String(initialTags.Name), Value: pulumi.String(vpcLogicalName)},
				ec2.VPCTagArgs{Key: pulumi.String(initialTags.Env), Value: pulumi.String(cfg.Env)},
			},
		),
	})

	if err != nil {
		return nil, fmt.Errorf("CreateVpc: failed creating vpc resource %[1]w", err)
	}

	//Create our Internet Gateway for each VPC using the classic aws provider (alias cec2) as aws-native does not
	//provide attachment to vpc yet. We reassign the 'err' variable
	_, err = cec2.NewInternetGateway(ctx, fmt.Sprintf("%[1]s-igw", vpcLogicalName), &cec2.InternetGatewayArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			initialTags.Name: pulumi.String(fmt.Sprintf("%[1]s-igw", vpcLogicalName)),
			initialTags.Env:  pulumi.String(cfg.Env),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("CreateVpc: failed creating igw resource %[1]w", err)
	}

	ctx.Export(vpcOutputName, vpc.ID())

	return VPC, nil
}
