package vpc

import (
	"core-infra/config"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	//subnetLogicalNames []string
	subnetOutputNames []string
)

func CreateSubnet(ctx *pulumi.Context) error {
	//Load configurations
	cd := config.NewConfig(ctx)
	configSubnetNames := cd.SubnetNames
	configAz := cd.Az

	//we are setting 2 subnets (private and public) per az in the 3 az available.

	//Loop through the exported vpc resources
	for _, name := range vpcOutputNames {
		vpcResult, err := ec2.LookupVPC(ctx, &ec2.LookupVPCArgs{
			VpcId: name,
		})

		if err != nil {
			return err
		}

		for index, subnetName := range configSubnetNames {
			//Define the names of each subnet resource. This is for logical resources (get id, arn etc.),
			//physical names have a random suffix and also Create names (adding "-out" at the end of the vpc names)
			//for ids to later be used as keys when exporting the resource

			//subnetLogicalNames = append(subnetLogicalNames, config.FormatName(configEnv, configRegion, vpcName))
			subnetOutputNames = append(subnetOutputNames, config.FormatName(subnetName, "out"))

			for _, azName := range configAz {

				//Create our subnets taking the exported vpc id as parameters and the az taken from the config settings
				subnet, err := ec2.NewSubnet(ctx, subnetName, &ec2.SubnetArgs{
					AvailabilityZone: pulumi.String(azName), VpcId: pulumi.String(*vpcResult.VpcId),
				})

				if err != nil {
					return err
				}

				ctx.Export(subnetOutputNames[index], subnet.ID())

			}

		}

	}

	return nil
}
