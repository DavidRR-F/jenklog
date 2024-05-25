package jenkins

import (
	"bytes"
	"fmt"
)

const (
	PIPELINE_START = "[Pipeline] { (%s)"
	PIPELINE_END   = "[Pipeline] }"
)

type Log struct {
	id    string
	stage string
	bytes []byte
}

func (l *Log) ParseByStage(stage string) error {
	startMarker := []byte(fmt.Sprintf(PIPELINE_START, stage))
	startIndex := bytes.Index(l.bytes, startMarker)
	if startIndex == -1 {
		return fmt.Errorf("stage %s not found", stage)
	}
	startIndex += len(startMarker)

	endIndex := bytes.Index(l.bytes[startIndex:], []byte(PIPELINE_END))
	if endIndex == -1 {
		return fmt.Errorf("stage %s not found", stage)
	}
	endIndex += startIndex

	l.stage = stage
	l.bytes = l.bytes[startIndex:endIndex]
	return nil
}

func (l *Log) Print() {
	fmt.Printf("\nID: %s\nStage: %s\n%s", l.id, l.stage, string(l.bytes))
}
