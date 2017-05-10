package lpass

import (
	"os/exec"
	"fmt"
	"bytes"
	"gopkg.in/yaml.v2"
)

type Note map[string]string

type LPass struct {}

type NoteIdentifier struct {
	NoteName string
	KeyName string
}

func (l *LPass) Notes(identifiers ...NoteIdentifier) ([]Note, error) {
	var notes []Note

	for _, identifier := range identifiers {
		note, err := l.Note(identifier)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (l *LPass) Note(identifier NoteIdentifier) (Note, error) {
	cmd := exec.Command("lpass", "show", "--notes", fmt.Sprintf(`%s`, identifier.NoteName))

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var note Note
	err = yaml.Unmarshal(stdout.Bytes(), &note)
	if err != nil {
		return nil, err
	}

	return note, nil
}
