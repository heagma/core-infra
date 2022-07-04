package vpc

import (
	"core-infra/config"
	"fmt"
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

	//we are setting 2 subnets (private and public) per az in the 3 az available.

	//Loop through the exported vpc resources

	for _, vpc := range VPC {

		//vpcResult := ec2.LookupVPCOutput(ctx, ec2.LookupVPCOutputArgs{
		//	VpcId: vpc.VpcId,
		//})

		for index, azName := range configAz {
			fmt.Printf("the AZ are: %[1]s\n", azName)

			for _, subnetName := range configSubnetNames {
				fmt.Println(subnetName)

				//Define the names of each subnet resource. This is for logical resources (get id, arn etc.),
				//physical names have a random suffix and also Create names (adding "-out" at the end of the vpc names)
				//for ids to later be used as keys when exporting the resource

				subnetLogicalNames = append(subnetLogicalNames, config.FormatName(configEnv, configRegion, azName, subnetName))
				subnetOutputNames = append(subnetOutputNames, config.FormatName(subnetName, "out"))

				fmt.Printf("logical names %[1]v\n", subnetLogicalNames[index])

				//Create our subnets taking the exported vpc id as parameters and the az taken from the config settings
				subnet, err := ec2.NewSubnet(ctx, subnetLogicalNames[index], &ec2.SubnetArgs{
					AvailabilityZone: pulumi.String(azName), VpcId: vpc.VpcId,
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
