package vpc

import (
	"core-infra/config"
	"fmt"
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws/ec2"
	cec2 "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
	"strings"
)

//load initial tags
var initialTags = config.NewInitialTags()

func CreateRouteTable(ctx *pulumi.Context, cfg *config.DataConfig, vpc *CustomVPC, subnets []*CustomSubnet) ([]*cec2.RouteTable, error) {
	var rts []*cec2.RouteTable
	const (
		pubRtSuffix = "rt-public"
		prvRtSuffix = "rt-private"
	)
	//Create the Route Table for all the public subnets
	pubRtTempName := config.FormatName(vpc.logicalName, pubRtSuffix)
	pubRt, err := cec2.NewRouteTable(ctx, pubRtTempName, &cec2.RouteTableArgs{
		VpcId: vpc.ID(),
		Tags: pulumi.StringMap{
			initialTags.Name: pulumi.String(pubRtTempName),
			initialTags.Env:  pulumi.String(cfg.Env),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("RouteTable: failed creating public route table resource %[1]w", err)
	}

	rts = append(rts, pubRt)

	//Create Routes for the public route table
	_, err = ec2.NewRoute(ctx, "route", &ec2.RouteArgs{
		RouteTableId:         pubRt.ID(),
		DestinationCidrBlock: pulumi.String(cfg.VpcCidr),
		GatewayId:            vpc.IGW.ID(),
	})
	if err != nil {
		return nil, err
	}

	//Create Route Tables for each Private Subnet
	for i, subnet := range subnets {
		if strings.HasSuffix(subnet.logicalName, "private") {
			prvRtTempName := config.FormatName(vpc.logicalName, prvRtSuffix, strconv.Itoa(i))
			prvRt, err := cec2.NewRouteTable(ctx, prvRtTempName, &cec2.RouteTableArgs{
				VpcId: vpc.ID(),
				Tags: pulumi.StringMap{
					initialTags.Name: pulumi.String(prvRtTempName),
					initialTags.Env:  pulumi.String(cfg.Env),
				},
			})
			if err != nil {
				return nil, fmt.Errorf("RouteTable: failed creating private route table resource %[1]w", err)
			}
			rts = append(rts, prvRt)
		}

		//Create Routes for each private subnet

	}

	return rts, nil
}

//
////AssociateRouteTable associate a route table with an entity
//func AssociateRouteTable() (ec2.RouteTable, error) {
//	pass
//}
