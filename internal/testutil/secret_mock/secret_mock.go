package secret_mock

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/secret"

func New() secret.Provider {
	return &mockSecretProvider{}
}

type mockSecretProvider struct {
}

func (m *mockSecretProvider) Secret() ([]byte, error) {
	return []byte("secret"), nil
}

func (m *mockSecretProvider) Close() {

}
