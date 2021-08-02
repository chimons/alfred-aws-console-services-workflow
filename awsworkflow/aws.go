package awsworkflow

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func InitAWS(transport http.RoundTripper, profile *awsconfig.Profile, region *awsconfig.Region) aws.Config {
	profileLoadOptionsFunc := func(o *config.LoadOptions) error {
		return nil
	}
	regionLoadOptionsFunc := func(o *config.LoadOptions) error {
		return nil
	}
	if profile != nil {
		profileLoadOptionsFunc = config.WithSharedConfigProfile(profile.Name)
		if profile.Region != "" {
			regionLoadOptionsFunc = config.WithRegion(profile.Region)
		}
	}
	if region != nil {
		regionLoadOptionsFunc = config.WithRegion(region.Name)
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		profileLoadOptionsFunc,
		regionLoadOptionsFunc,
	)
	if err != nil {
		panic(err)
	}

	if transport != nil {
		cfg.HTTPClient = &http.Client{
			Transport: transport,
		}
	}

	initAWSConsoleDomain(cfg.Region)

	return cfg
}

var defaultAwsConsoleDomain string = "console.aws.amazon.com"
var defaultAwsConsoleDomainChina string = "console.amazonaws.cn"

func initAWSConsoleDomain(region string) {
	awsConsoleDomain, _ := os.LookupEnv("ALRED_AWS_CONSOLE_SERVICES_WORKFLOW_AWS_CONSOLE_DOMAIN")
	if awsConsoleDomain == "" {
		if strings.HasPrefix(region, "cn-") {
			awsConsoleDomain = defaultAwsConsoleDomainChina
		} else {
			awsConsoleDomain = defaultAwsConsoleDomain
		}
	}
	util.AWSConsoleDomain = awsConsoleDomain
}

func GetImageIcon(id string) *aw.Icon {
	icon := &aw.Icon{Value: "images/" + id + ".png"}
	return icon
}
