package ignores

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"log"
	"testing"
)

type MockedLoader struct {
	dataInBytes []byte
}

func (l *MockedLoader) LoadIgnores(in map[string]interface{}) {
	yamlBytes, err := yaml.Marshal(in)
	if err != nil {
		log.Fatalf("Unable to load provided Ignores: %v", in)
	}

	l.dataInBytes = yamlBytes
}

func (l *MockedLoader) getIgnoresFromSource() ([]byte, error) {
	return l.dataInBytes, nil
}

func TestLoaderFromFileEmptyIgnores(t *testing.T) {
	loader := MockedLoader{}

	emptyMap := map[string]interface{}{}
	loader.LoadIgnores(emptyMap)

	ignores, err := GetSubIgnores(&loader)
	assert.NoError(t, err)

	expected := Config{}
	assert.EqualValues(t, expected, *ignores)
}

func TestLoaderFromFileIgnoresCorrect(t *testing.T) {
	loader := MockedLoader{}

	subsToIgnore := map[string]interface{}{
		"subscriptionsToIgnore": map[string][]string{
			"12345": nil,
			"54321": nil,
		},
	}
	loader.LoadIgnores(subsToIgnore)

	ignores, err := GetSubIgnores(&loader)
	assert.NoError(t, err)

	expected := Config{
		SubscriptionsToIgnore: map[string][]string{
			"12345": []string{},
			"54321": []string{},
		},
	}
	assert.EqualValues(t, expected, *ignores)
}

func TestLoaderFromFileOnlyLookup(t *testing.T) {
	loader := MockedLoader{}

	subsToSearch := map[string]interface{}{
		"onlySubscriptionsToLookup": map[string][]string{
			"12345": nil,
			"54321": nil,
		},
	}
	loader.LoadIgnores(subsToSearch)

	ignores, err := GetSubIgnores(&loader)
	assert.NoError(t, err)

	expected := Config{
		OnlySubscriptionsToLookup: map[string][]string{
			"12345": []string{},
			"54321": []string{},
		},
	}
	assert.EqualValues(t, expected, *ignores)
}
