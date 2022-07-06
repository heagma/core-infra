package vpc

import (
	"core-infra/config"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateSubnet(ctx *pulumi.Context, VPCs []*ec2.VPC) (error, []*ec2.Subnet) {
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
	configSubnetCidr := cd.SubnetCidr
	configAz := cd.Az

	//load initial tags
	initialTags := config.NewInitialTags()

	//Loop through the exported vpc resources
	for _, vpc := range VPCs {
		//we are setting 2 subnets (private and public) per az in the 3 az available.
		for _, subnetType := range configSubnetTypes {
			for _, azName := range configAz {

				//Define the names of each subnet resource. This is for logical resources (get id, arn etc.),
				//physical names have a random suffix and also Create names (adding "-out" at the end of the vpc names)
				//for ids to later be used as keys when exporting the resource in case another project/stack needs it

				subnetLogicalNames = append(subnetLogicalNames, config.FormatName(configEnv, configRegion, azName, subnetType))
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
					CidrBlock: pulumi.String(configSubnetCidr[i]),
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
					return err, Subnets
				}

				Subnets = append(Subnets, subnet)

				ctx.Export(subnetLogicalName, subnet.ID())

				break

			}
		}

	}

	return nil, Subnets
}
