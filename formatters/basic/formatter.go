// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package basic

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

const BasicFormatterType string = "basic"

type BasicFormatter struct {
	Config *Config
}

func (f *BasicFormatter) Type() string {
	return BasicFormatterType
}

func (f *BasicFormatter) Format(yamlContent []byte) ([]byte, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(yamlContent))
	documents := []yaml.Node{}
	for {
		var docNode yaml.Node
		err := decoder.Decode(&docNode)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		documents = append(documents, docNode)
	}

	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	e.SetIndent(f.Config.Indent)
	for _, doc := range documents {
		err := e.Encode(&doc)
		if err != nil {
			return nil, err
		}
	}

	if f.Config.IncludeDocumentStart {
		return withDocumentStart(b.Bytes()), nil
	}
	return b.Bytes(), nil
}

func withDocumentStart(document []byte) []byte {
	documentStart := "---\n"
	return append([]byte(documentStart), document...)
}
