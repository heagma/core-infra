package vpc

import (
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	vpcResult          ec2.LookupVPCResult
	subnetLogicalNames []string
	subnetOutputNames  []string
)

func CreateSubnet() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		//Load configurations
		cd := loadConfig(ctx)
		configSubnetNames := cd.SubnetNames
		configAz := cd.Az

		//Get availability zones
		//we are setting 2 subnets per az in the 3 az available.

		for index, subnetName := range configSubnetNames {
			for _, az := range configAz {

			}

		}
		//Loop through the exported vpcs resources
		for _, name := range vpcOutputNames {
			vpcId, err := ec2.LookupVPC(ctx, &ec2.LookupVPCArgs{
				VpcId: name,
			})

			if err != nil {
				return err
			}

			//Create our subnets
			subnet, error := ec2.NewSubnet(ctx)

		}

		return nil
	})

}
