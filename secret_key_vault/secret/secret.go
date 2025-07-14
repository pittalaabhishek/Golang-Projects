package secret

import (
	"encoding/json"
	"errors"
	"os"

	"secret_key_vault/encrypt"
)

type FileVault struct {
	Key  string
	Path string
}

func NewFileVault(key, path string) *FileVault {
	return &FileVault{Key: key, Path: path}
}

func (fv *FileVault) readSecrets() (map[string]string, error) {
	data, err := os.ReadFile(fv.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]string), nil // file not found = empty map
		}
		return nil, err
	}
	decrypted, err := encrypt.Decrypt(string(data), fv.Key)
	if err != nil {
		return nil, err
	}
	var secrets map[string]string
	err = json.Unmarshal([]byte(decrypted), &secrets)
	return secrets, err
}

func (fv *FileVault) writeSecrets(secrets map[string]string) error {
	jsonData, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	encrypted, err := encrypt.Encrypt(string(jsonData), fv.Key)
	if err != nil {
		return err
	}
	return os.WriteFile(fv.Path, []byte(encrypted), 0600)
}

func (fv *FileVault) Set(key, value string) error {
	secrets, err := fv.readSecrets()
	if err != nil {
		return err
	}
	secrets[key] = value
	return fv.writeSecrets(secrets)
}

func (fv *FileVault) Get(key string) (string, error) {
	secrets, err := fv.readSecrets()
	if err != nil {
		return "", err
	}
	val, ok := secrets[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return val, nil
}