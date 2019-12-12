package github

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Label represents GitHub label.
type Label struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Color       string `yaml:"color"`
}

// FromManifestToLabels loads a YAML file and umarshal it to []Label.
func FromManifestToLabels(path string) ([]Label, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var labels []Label
	err = yaml.Unmarshal(buf, &labels)
	return labels, err
}
