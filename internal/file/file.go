package file

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"time"

	sops "go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/aes"
	"go.mozilla.org/sops/v3/cmd/sops/common"
	"go.mozilla.org/sops/v3/config"
	"go.mozilla.org/sops/v3/version"
)

type File struct {
	Encrypted bool
	cipher    aes.Cipher
	config    *config.Config
	store     common.Store
	tree      *sops.Tree
}

func Load(filePath string) (*File, error) {
	configFile, err := config.FindConfigFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot find config file: %w", err)
	}

	config, err := config.LoadCreationRuleForFile(configFile, filePath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot load config file: %w", err)
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read input file: %w", err)
	}

	fs := File{
		cipher: aes.NewCipher(),
		store:  common.DefaultStoreForPath(filePath),
		config: config,
	}

	tree, err := fs.store.LoadEncryptedFile(fileBytes)
	switch {
	default:
		return nil, err
	case errors.Is(err, sops.MetadataNotFound):
		branches, err := fs.store.LoadPlainFile(fileBytes)
		if err != nil {
			return nil, fmt.Errorf("cannot load plain file %s: %w", filePath, err)
		}

		tree = sops.Tree{
			Metadata: sops.Metadata{
				KeyGroups:         config.KeyGroups,
				EncryptedSuffix:   config.EncryptedSuffix,
				EncryptedRegex:    config.EncryptedRegex,
				UnencryptedRegex:  config.UnencryptedRegex,
				UnencryptedSuffix: config.UnencryptedSuffix,
				Version:           version.Version,
				ShamirThreshold:   config.ShamirThreshold,
				LastModified:      time.Now().UTC(),
			},
			Branches: branches,
			FilePath: filePath,
		}

		fs.Encrypted = false
	case errors.Is(err, nil):
		fs.Encrypted = true
	}

	fs.tree = &tree
	fs.tree.FilePath = filePath

	return &fs, nil
}

// IdentifySecretFiles returns a list of all files in the given directory that
func IdentifySecretFiles(directory string, secretRegexp string) ([]string, error) {
	fileMatcher, err := regexp.Compile(secretRegexp)
	if err != nil {
		return nil, fmt.Errorf("invalid secret-regexp (%q): %w", secretRegexp, err)
	}

	var secretFiles []string

	err = filepath.WalkDir(directory, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !fileMatcher.MatchString(filePath) {
			return nil
		}

		secretFiles = append(secretFiles, filePath)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return secretFiles, nil
}
