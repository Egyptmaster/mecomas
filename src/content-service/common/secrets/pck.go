package secrets

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Secret struct {
	User     string
	Password string
}

type Secrets map[string]Secret

func Load(yml string) (Secrets, error) {
	file, err := os.Open(yml)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()
	dec := yaml.NewDecoder(file)

	var secrets Secrets
	return secrets, dec.Decode(&secrets)
}

func (secrets Secrets) Get(key string) (Secret, error) {
	if s, ok := secrets[key]; ok {
		return s, nil
	}
	return Secret{}, fmt.Errorf("the key %s does not exist", key)
}
