package pipelines

import (
	"io/ioutil"
	"regexp"
	"strings"
	"github.com/mbildner/pipecleaner/lib/lpass"
)

type Pipeline struct {
	Definition string
}

func New (definitionPath string) (*Pipeline, error) {
	pipelineDefinition, err := ioutil.ReadFile(definitionPath)
	if err != nil {
		return nil, err
	}

	pipeline := Pipeline{
		Definition: string(pipelineDefinition),
	}

	return &pipeline, nil
}

func (p *Pipeline) Secrets () []lpass.NoteIdentifier {
	re := regexp.MustCompile(`\({2}(.+)\){2}`)
	definitions := re.FindAllString(p.Definition, -1)

	secretInjections := []lpass.NoteIdentifier{}
	for _, noteDefinition := range definitions {
		secret := buildSecret(clean(noteDefinition))
		secretInjections = append(secretInjections, secret)
	}

	return secretInjections
}

func buildSecret(noteDefinition string) lpass.NoteIdentifier {
	parts := strings.Split(noteDefinition, "/Notes/")
	note := parts[0]
	key := parts[1]

	secret := lpass.NoteIdentifier{
		NoteName: note,
		KeyName:  key,
	}

	return secret
}

func clean (secret string) string {
	return strings.TrimRight(strings.TrimLeft(secret, "(("), "))")
}