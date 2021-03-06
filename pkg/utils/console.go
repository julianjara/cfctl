package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type MessageType string

// Message types
const (
	MessageTypeInfo MessageType = "info"

	MessageTypeWarn MessageType = "warn"

	MessageTypeError MessageType = "error"
)

// Print format types
type FormatType string

const (
	FormatYaml FormatType = "yaml"
	FormatJson FormatType = "json"
	FormatCmd  FormatType = "cmd"
)

// Function type to output to stdout
type StdoutStrFn func(input interface{}) (string, error)

// Convert to yaml string
func toYamlStr(input interface{}) (string, error) {
	output, err := yaml.Marshal(input)
	return string(output), err
}

// Convert to json string with indent
func toJsonStr(input interface{}) (string, error) {
	output, err := json.MarshalIndent(input, "", "    ")
	return string(output), err
}

// default
func toDefault(input interface{}) (string, error) {
	return fmt.Sprintf("%s", input), nil
}

// stdout string conversion factory
// Default to json
func StdoutStrFactory(format FormatType) StdoutStrFn {
	switch format {
	case FormatYaml:
		return toYamlStr
	case FormatJson:
		return toJsonStr
	case FormatCmd:
		return toDefault
	default:
		return toDefault
	}
}

// Print given interface to given format
func Print(format FormatType, s ...interface{}) error {
	if len(s) == 0 {
		return errors.New(MsgFormat("Printing output error: No object given", MessageTypeError))
	}

	fn := StdoutStrFactory(format)
	for _, ss := range s {
		out, err := fn(ss)
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	return nil
}

// Format error message
func MsgFormat(msg string, msgType MessageType, options ...string) string {
	return fmt.Sprintf("%s", msg)
}

// Generic Print info.
func InfoPrint(s ...interface{}) error {
	return Print(FormatCmd, s...)
}

// Print to stdout with info header.
func StdoutInfo(s ...interface{}) error {
	s = append([]interface{}{fmt.Sprintf("[ %s ] ", MessageTypeInfo)}, s...)
	_, err := fmt.Print(s...)
	return err
}

// Print to stdout with warn header.
func StdoutWarn(s ...interface{}) error {
	s = append([]interface{}{fmt.Sprintf("[ %s ] ", MessageTypeWarn)}, s...)
	_, err := fmt.Print(s...)
	return err
}

// Print to stdout with error header.
func StdoutError(s ...interface{}) error {
	s = append([]interface{}{fmt.Sprintf("[ %s ] ", MessageTypeError)}, s...)
	_, err := fmt.Print(s...)
	return err
}
