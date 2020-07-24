package searchers

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/tests"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSearchElasticBeanstalkEnvironments(t *testing.T) {
	wf := aw.New()

	session, r := tests.NewAWSRecorderSession(util.GetCurrentFilename())
	defer tests.PanicOnError(r.Stop)
	err := SearchElasticBeanstalkEnvironments(wf, "elasticbeanstalk", session, true, "")
	if err != nil {
		t.Errorf("got error from search: %v", err)
	}

	cupaloy.SnapshotT(t, wf.Feedback.Items)
}