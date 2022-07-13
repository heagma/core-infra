package vpc

//Still pending the call to a subnet type to convert to string

import (
	"core-infra/config"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	cec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

func AddNat(ctx *pulumi.Context, cfg *config.DataConfig, Subnets []*CustomSubnet) ([]*ec2.NatGateway, error) {
	var NAT []*ec2.NatGateway

	configEnv := cfg.Env

	//load initial tags
	initialTags := config.NewInitialTags()

	for _, subnet := range Subnets {

		//Attach a Nat gateway to each public subnet
		if strings.HasSuffix(subnet.logicalName, "public") {

			//Create an EIP first to be associated with the NATGateway
			//Lets give our eip a temp name to be used as a logical name
			eipTempName := config.FormatName(subnet.logicalName, "nat", "eip")

			eip, err := cec2.NewEip(ctx, eipTempName, &cec2.EipArgs{
				Tags: pulumi.StringMap{
					initialTags.Name: pulumi.String(eipTempName),
					initialTags.Env:  pulumi.String(configEnv),
				},
			})

			if err != nil {
				return nil, fmt.Errorf("CreateSubnet: failed creating eip resource %[1]w", err)
			}

			nat, err := ec2.NewNatGateway(ctx, fmt.Sprintf("%[1]s-nat", subnet.logicalName), &ec2.NatGatewayArgs{
				SubnetId:     subnet.ID(),
				AllocationId: eip.ID(),
				Tags: ec2.NatGatewayTagArrayInput(
					ec2.NatGatewayTagArray{
						ec2.NatGatewayTagArgs{
							Key:   pulumi.String(initialTags.Name),
							Value: pulumi.String(fmt.Sprintf("%[1]s-nat", subnet.logicalName)),
						},
						ec2.NatGatewayTagArgs{
							Key: pulumi.String(initialTags.Env), Value: pulumi.String(configEnv),
						},
					}),
			})
			if err != nil {
				return nil, fmt.Errorf("CreateSubnet: failed creating nat gateway resource %[1]w", err)
			}

			NAT = append(NAT, nat)
		}

	}

	return NAT, nil

}
