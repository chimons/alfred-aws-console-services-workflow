package main

import (
	"io/ioutil"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/core"
	"github.com/rkoval/alfred-aws-console-services-workflow/filters"
	"gopkg.in/yaml.v2"
)

var wf *aw.Workflow

func init() {
	wf = aw.New()
}

func parseYaml() []core.AwsService {
	awsServices := []core.AwsService{}
	yamlFile, err := ioutil.ReadFile("console-services.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &awsServices)
	if err != nil {
		log.Fatal(err)
	}
	return awsServices
}

func populateItems(awsServices []core.AwsService, query string) (string, error) {
	awsServicesById := make(map[string]*core.AwsService)
	for i, awsService := range awsServices {
		awsServicesById[awsService.Id] = &awsServices[i]
	}

	// TODO add better lexing here to route filters

	splitQuery := strings.Split(query, " ")
	if len(splitQuery) > 1 && awsServicesById[splitQuery[0]] != nil {
		id := splitQuery[0]
		query = strings.Join(splitQuery[1:], " ")
		awsService := awsServicesById[id]
		searcher := filters.SearchersByServiceId[id]
		if strings.HasPrefix(query, "$") && searcher != nil {
			query = query[1:]
			log.Printf("using searcher associated with %s", id)
			err := searcher(wf, query)
			if err != nil {
				return "", err
			}
			return query, nil
		} else if len(awsServicesById[splitQuery[0]].Sections) > 0 {
			log.Printf("filtering on sections for %s", id)
			filters.ServiceSections(wf, *awsService, query)
			return query, nil
		}
	}

	filters.Services(wf, awsServices, query)
	return query, nil
}

func run() {
	var query string
	if args := wf.Args(); len(args) > 0 {
		query = strings.TrimSpace(args[0])
	}

	awsServices := parseYaml()

	query, err := populateItems(awsServices, query)

	if err != nil {
		log.Printf("error: %v", err)
	} else if query != "" {
		log.Printf("filtering with query %s", query)
		res := wf.Filter(query)

		log.Printf("%d results match %q", len(res), query)

		for i, r := range res {
			log.Printf("%02d. score=%0.1f sortkey=%s", i+1, r.Score, wf.Feedback.Keywords(i))
		}
	}

	wf.WarnEmpty("No matching services found", "Try a different query?")

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
