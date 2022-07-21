package vpc

import (
	"core-infra/config"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//CustomSubnet struct is a composite (not inheritance) struct that embed the ec2.Subnet struct
//This allows us to create a new type of object that also include a name (there is no name field on ec2.Subnet) for our
//subnets to be used as input in another caller instead of relying on ApplyT for Output<T> native from Pulumi.
type CustomSubnet struct {
	logicalName string
	*ec2.Subnet
}

func CreateSubnet(ctx *pulumi.Context, cfg *config.DataConfig, vpc *CustomVPC) ([]*CustomSubnet, error) {
	var (
		Subnets            []*CustomSubnet
		subnetLogicalNames []string
		subnetOutputNames  []string
	)

	//load initial tags
	initialTags := config.NewInitialTags()

	//we are setting 2 subnets (private and public) per az in the 3 az available.
	for _, subnetType := range cfg.SubnetTypes {
		for _, azName := range cfg.Az {

			//Define the names of each subnet resource. This is for logical resources (get id, arn etc.),
			//Outputs names is to later be used as keys in case another project/stack needs it.

			subnetLogicalNames = append(subnetLogicalNames, config.FormatName(cfg.Env, cfg.RegionAlias, azName, "sn", subnetType))
			subnetOutputNames = append(subnetOutputNames, config.FormatName(azName, subnetType, "out"))
		}
	}

	//Loop through the logical names to create named subnets within each AZ
	for i, subnetLogicalName := range subnetLogicalNames {
		for _, azName := range cfg.Az {
			//Create our subnets taking the exported vpc id as parameters and the az taken from the config settings
			//here we assign the logical names at the same time unlike when creating our VPC
			subnet, err := ec2.NewSubnet(ctx, subnetLogicalName, &ec2.SubnetArgs{
				AvailabilityZone: pulumi.String(azName), VpcId: vpc.VpcId,
				CidrBlock: pulumi.String(cfg.SubnetCidrs[i]),
				Tags: ec2.SubnetTagArrayInput(
					ec2.SubnetTagArray{
						ec2.SubnetTagArgs{
							Key: pulumi.String(initialTags.Name), Value: pulumi.String(subnetLogicalName),
						},
						ec2.SubnetTagArgs{
							Key: pulumi.String(initialTags.Env), Value: pulumi.String(cfg.Env),
						},
					},
				),
			})

			if err != nil {
				return nil, fmt.Errorf("CreateSubnet: failed creating subnet resource %[1]w", err)
			}

			Subnets = append(Subnets, &CustomSubnet{subnetLogicalName, subnet})

			ctx.Export(subnetLogicalName, subnet.ID())

			break

		}
	}

	return Subnets, nil
}
