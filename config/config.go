package config

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"strings"
)

//DataConfig struct to store our configuration data from the stack settings file
type DataConfig struct {
	Env         string
	RegionAlias string
	VpcCidr     []string
	VpcNames    []string
	SubnetTypes []string
	SubnetCidr  []string
	Az          []string
}

//NewConfig function loads the config values from a settings file and
//return a named parameter, in this case confData of type ConfigData
func NewConfig(ctx *pulumi.Context) (confData *DataConfig) {

	//Load configurations
	cfg := config.New(ctx, "")
	cfg.RequireObject("data", &confData)

	return
}

type commonTags struct {
	Name string
	Env  string
}

//NewInitialTags constructor creates some initial tags name that is used very often
func NewInitialTags() *commonTags {
	return &commonTags{Name: "Name", Env: "Env"}
}

var InitialTags = NewInitialTags()

//FormatName function format the name with n parameters
func FormatName(values ...string) string {
	var name string
	for _, value := range values {
		name += value + "-"
	}
	return strings.TrimRight(name, "-")
}
