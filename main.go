package main

import (
	"fmt"
	"github.com/mbildner/pipecleaner/lib/lpass"
	"github.com/mbildner/pipecleaner/lib/pipelines"
	"github.com/mbildner/pipecleaner/lib/concourse"
	"os"
)

func main() {
	pipelineDefinitionPath := "/Users/moshe/workspace/concourse_demo/pipeline.yml"

	pipeline, err := pipelines.New(pipelineDefinitionPath)
	if err != nil {
		fmt.Println("could not read pipeline from file, bailing")
		os.Exit(1)
	}
	secretDefinitions := pipeline.Secrets()

	lpass := lpass.LPass{}
	notes, err := lpass.Notes(secretDefinitions...)
	if err != nil {
		fmt.Println("could not read notes from LastPass, bailing")
		os.Exit(1)
	}

	c := concourse.Concourse{FlyTarget: "pipeline-bling", PipelineName: "hello-world"}
	jobs, err := c.Jobs("hello-world")
	if err != nil {
		fmt.Println("could not get pipeline jobs, bailing")
		os.Exit(1)
	}

	for _, job := range jobs {
		fmt.Println(c.Builds("hello-world", job))
	}

	for _, note := range notes {
		for k, _ := range note {
			fmt.Println(k)
		}
	}

}
