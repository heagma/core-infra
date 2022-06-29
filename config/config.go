package vpc

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	"strings"
)

//ConfigData struct to store our configuration data from the stack settings file
type ConfigData struct {
	Env         string
	RegionAlias string
	VpcCidr     []string
	VpcNames    []string
	SubnetNames []string
	Az          []string
}

//loadConfig function loads the config values from a settings file and
//return a named parameter, in this case confData of type ConfigData
func loadConfig(ctx *pulumi.Context) (confData ConfigData) {

	//Load configurations
	cfg := config.New(ctx, "")
	cfg.RequireObject("data", &confData)

	return
}

type commonTags struct {
	Name string
	Env  string
}

//newInitialTags constructor creates some initial tags name that is used very often
func newInitialTags() *commonTags {
	return &commonTags{Name: "Name", Env: "Env"}
}

var initialTags = newInitialTags()

//formatName function format the name with n parameters
func formatName(values ...string) string {
	var name string
	for _, value := range values {
		name += value + "-"
	}
	return strings.TrimRight(name, "-")
}
