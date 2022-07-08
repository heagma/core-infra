package vpc

import (
	"core-infra/config"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

func CreateSubnet(ctx *pulumi.Context, VPCs []*ec2.VPC) ([]*ec2.Subnet, error) {
	var (
		Subnets            []*ec2.Subnet
		subnetLogicalNames []string
		subnetOutputNames  []string
	)

	//Load configurations
	cd := config.NewConfig(ctx)

	configEnv := cd.Env
	configRegion := cd.RegionAlias

	//Assign configurations related only to vpc
	configSubnetTypes := cd.SubnetTypes
	configSubnetCidrs := cd.SubnetCidrs
	configAz := cd.Az

	//load initial tags
	initialTags := config.NewInitialTags()

	//Loop through the exported vpc resources
	for _, vpc := range VPCs {
		//we are setting 2 subnets (private and public) per az in the 3 az available.
		for _, subnetType := range configSubnetTypes {
			for _, azName := range configAz {

				//Define the names of each subnet resource. This is for logical resources (get id, arn etc.),
				//Outputs names is to later be used as keys in case another project/stack needs it

				subnetLogicalNames = append(subnetLogicalNames, config.FormatName(configEnv, configRegion, azName, "sn", subnetType))
				subnetOutputNames = append(subnetOutputNames, config.FormatName(azName, subnetType, "out"))
			}
		}

		//Loop through the logical names to create named subnets within each AZ
		for i, subnetLogicalName := range subnetLogicalNames {
			for _, azName := range configAz {
				//Create our subnets taking the exported vpc id as parameters and the az taken from the config settings
				//here we assign the logical names at the same time unlike when creating our VPC
				subnet, err := ec2.NewSubnet(ctx, subnetLogicalName, &ec2.SubnetArgs{
					AvailabilityZone: pulumi.String(azName), VpcId: vpc.VpcId,
					CidrBlock: pulumi.String(configSubnetCidrs[i]),
					Tags: ec2.SubnetTagArrayInput(
						ec2.SubnetTagArray{
							ec2.SubnetTagArgs{
								Key: pulumi.String(initialTags.Name), Value: pulumi.String(subnetLogicalName),
							},
							ec2.SubnetTagArgs{
								Key: pulumi.String(initialTags.Env), Value: pulumi.String(configEnv),
							},
						},
					),
				})

				if err != nil {
					return nil, err
				}

				Subnets = append(Subnets, subnet)

				ctx.Export(subnetLogicalName, subnet.ID())

				//Create NAT gateways inside the public subnets
				if strings.HasSuffix(subnetLogicalName, "public") {
					//Create an EIP first to be associated with the NATGateway

					_, err = ec2.NewNatGateway(ctx, fmt.Sprintf("%[1]s-nat", subnetLogicalName), &ec2.NatGatewayArgs{
						SubnetId: subnet.ID(),
						Tags: ec2.NatGatewayTagArrayInput(
							ec2.NatGatewayTagArray{
								ec2.NatGatewayTagArgs{
									Key:   pulumi.String(initialTags.Name),
									Value: pulumi.String(fmt.Sprintf("%[1]s-nat", subnetLogicalName)),
								},
								ec2.NatGatewayTagArgs{
									Key: pulumi.String(initialTags.Env), Value: pulumi.String(configEnv),
								},
							}),
					})
					if err != nil {
						return nil, err
					}
				}

				break

			}
		}

	}

	return Subnets, nil
}
