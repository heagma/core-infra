package vpc

//
//import (
//	"core-infra/config"
//	"fmt"
//	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
//	cec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
//	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
//)
//
////load initial tags
//var initialTags = config.NewInitialTags()
//
//func CreateRouteTable(ctx *pulumi.Context, vpc *ec2.VPC, routeTableName string) (*cec2.RouteTable, error) {
//	rt, err := cec2.NewRouteTable(ctx, routeTableName, &cec2.RouteTableArgs{
//		VpcId: vpc.ID(),
//		Tags: pulumi.StringMap{
//			initialTags.Name: pulumi.String(routeTableName),
//			initialTags.Env:  pulumi.String(configEnv),
//		},
//	})
//	if err != nil {
//		return nil, fmt.Errorf("RouteTable: failed creating route table resource %[1]w", err)
//	}
//	return rt, nil
//}
//
//func CreateRoute() {
//
//}
//
////AttachRoute creates routes and associations
//func AssociateRouteTable() (ec2.RouteTable, error) {
//
//}
