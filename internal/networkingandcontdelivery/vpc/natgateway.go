package vpc

//
////Still pending the call to a subnet type to convert to string
//
//import (
//	"core-infra/config"
//	"fmt"
//	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
//	cec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//	"strings"
//)
//
//func AddNat(ctx *pulumi.Context, Subnets []*ec2.Subnet) ([]*ec2.NatGateway, error) {
//	var NAT []*ec2.NatGateway
//
//	//load initial tags
//	initialTags := config.NewInitialTags()
//
//	//Load configurations
//	cd := config.NewConfig(ctx)
//
//	configEnv := cd.Env
//
//	for _, subnet := range Subnets {
//
//		//Create NAT gateways inside the public subnets
//
//		//We want to check first the value of the tag Name of our subnets in order to check the subnet resource name
//		//as the subnet outputs is of type ec2.Subnet and its logical name is not passed to our nat function.
//		//And as the tag field from the subnets are of pulumi.StringOutput type then we can convert it using the
//		//output lifting of the underlying values. Reference: https://www.pulumi.com/docs/intro/concepts/inputs-outputs/
//		subnetNameTagValue := subnet.Tags.Index(pulumi.Int(0)).Value()
//
//		//this just return another output so any operation should be done inside the callback but the creation of new
//		//resources is discouraged as the preview could differ. And there is not possibilities to do assignmeng inside
//		//as all the vars are used inside  the block
//		//var t string
//		//_ = subnetNameTagValue.ApplyT(func(tag string) string {
//		//	fmt.Println(tag)
//		//	t = tag
//		//	return fmt.Sprintf("%s-t]", tag)
//		//
//		//}).(pulumi.StringOutput)
//
//		fmt.Println(subnetNameTagValue)
//
//		//And now we can check our condition on the tag name to look for public suffix
//		if strings.HasSuffix(subnetNameTagValue, "public") {
//
//			fmt.Println("true")
//
//			//Create an EIP first to be associated with the NATGateway
//
//			eip, err := cec2.NewEip(ctx, fmt.Sprintf("%[1]v-eip", subnetNameTagValue), &cec2.EipArgs{})
//
//			if err != nil {
//				return nil, fmt.Errorf("CreateSubnet: failed creating eip resource %[1]w", err)
//			}
//
//			nat, err := ec2.NewNatGateway(ctx, fmt.Sprintf("%[1]v-nat", subnetNameTagValue), &ec2.NatGatewayArgs{
//				SubnetId:     subnet.ID(),
//				AllocationId: eip.ID(),
//				Tags: ec2.NatGatewayTagArrayInput(
//					ec2.NatGatewayTagArray{
//						ec2.NatGatewayTagArgs{
//							Key:   pulumi.String(initialTags.Name),
//							Value: pulumi.String(fmt.Sprintf("%[1]v-nat", subnetNameTagValue)),
//						},
//						ec2.NatGatewayTagArgs{
//							Key: pulumi.String(initialTags.Env), Value: pulumi.String(configEnv),
//						},
//					}),
//			})
//			if err != nil {
//				return nil, fmt.Errorf("CreateSubnet: failed creating nat gateway resource %[1]w", err)
//			}
//
//			NAT = append(NAT, nat)
//		}
//
//	}
//
//	return NAT, nil
//
//}
