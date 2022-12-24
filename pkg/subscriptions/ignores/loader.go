package ignores

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type Config struct {
	SubscriptionsToIgnore     map[string][]string `yaml:"subscriptionsToIgnore"`
	OnlySubscriptionsToLookup map[string][]string `yaml:"onlySubscriptionsToLookup"`
}

type LoaderFromFile struct {
	Path *string
}

func GetSubIgnores(l Loader) (*Config, error) {
	dataInbytes, err := l.getIgnoresFromSource()
	if err != nil {
		log.Fatalf("Unable to read ignores from selected source: %v", err)
	}

	filterConfig := &Config{}
	if err := yaml.Unmarshal(dataInbytes, filterConfig); err != nil {
		log.Fatalf("Unable to deserialize read data: %v", err)
	}

	log.Printf("Deserialized: %v", filterConfig)
	return filterConfig, nil
}

func (i *LoaderFromFile) getIgnoresFromSource() ([]byte, error) {
	file, err := os.Open(*i.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	return io.ReadAll(file)
}

type Loader interface {
	getIgnoresFromSource() ([]byte, error)
}
