package vpc

import (
	"core-infra/config"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	subnetLogicalNames []string
	subnetOutputNames  []string
)

func CreateSubnet(ctx *pulumi.Context) error {

	//Load configurations
	cd := config.NewConfig(ctx)

	configSubnetNames := cd.SubnetNames
	configAz := cd.Az

	configEnv := cd.Env
	configRegion := cd.RegionAlias

	//Loop through the exported vpc resources

	for _, vpc := range VPC {
		//we are setting 2 subnets (private and public) per az in the 3 az available.
		for _, subnetName := range configSubnetNames {

			for _, azName := range configAz {

				//Define the names of each subnet resource. This is for logical resources (get id, arn etc.),
				//physical names have a random suffix and also Create names (adding "-out" at the end of the vpc names)
				//for ids to later be used as keys when exporting the resource

				subnetLogicalNames = append(subnetLogicalNames, config.FormatName(configEnv, configRegion, azName, subnetName))
				subnetOutputNames = append(subnetOutputNames, config.FormatName(azName, subnetName, "out"))
			}
		}

		//Loop through the logical names to create named subnets within each AZ
		for _, subnetLogicalName := range subnetLogicalNames {
			for _, azName := range configAz {
				//Create our subnets taking the exported vpc id as parameters and the az taken from the config settings
				//here we assign the logical names at the same time unlike when creating our VPC
				subnet, err := ec2.NewSubnet(ctx, subnetLogicalName, &ec2.SubnetArgs{
					AvailabilityZone: pulumi.String(azName), VpcId: vpc.VpcId,
				})

				if err != nil {
					return err
				}

				ctx.Export(subnetLogicalName, subnet.ID())

				break

			}
		}

	}

	return nil
}
