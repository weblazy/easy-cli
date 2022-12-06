package conf

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestGetConfig(t *testing.T) {
	c := GetConfig()
	s, err := yaml.Marshal(&c)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(s))
}
