package concourse

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
	"go/build"
)

type Concourse struct {
	FlyTarget string
	PipelineName string
}

type pipeline struct {
	Jobs []struct {
		Name string `yaml:"name"`
	} `yaml:"jobs"`
}

func (c *Concourse) Jobs() ([]string, error) {
	cmd := exec.Command("fly", "-t", c.FlyTarget, "get-pipeline", "-p", c.PipelineName)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var p pipeline
	yaml.Unmarshal(stdout.Bytes(), &p)

	jobs := []string{}
	for _, j := range p.Jobs {
		jobs = append(jobs, j.Name)
	}

	return jobs, nil
}

func (c *Concourse) Builds(jobName string) ([]Build, error) {
	//cmd := exec.Command("fly", "-t", c.FlyTarget, "watch", "-j", pipelineName + "/" + jobName)
	cmd := exec.Command("fly", "-t", c.FlyTarget, "builds", "-j", c.PipelineName+"/"+jobName)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	builds, err := c.parseBuilds(stdout.String())
	if err != nil {
		fmt.Println("could not parse builds, bailing")
		return nil, err
	}

	return builds, nil
}

type Build struct {
	Id        int
	Succeeded bool
	Log       string
}

func (c *Concourse) parseBuilds(buildsText string) ([]Build, error) {
	var builds []Build

	lines := strings.Split(buildsText, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		parts := strings.Fields(line)
		fmt.Println(parts)
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}
		success := false
		if parts[3] == "succeeded" {
			success = true
		}

		log := c.BuildLog()
		builds = append(builds, Build{Id: id, Succeeded: success})

	}

	return builds, nil
}

func (c *Concourse) BuildLog(jobName string, id int) (string, error) {
	idString := strconv.Itoa(id)
	cmd := exec.Command("fly", "-t", c.FlyTarget, "watch", "-j", c.PipelineName+"/"+jobName, "-b", idString)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}
