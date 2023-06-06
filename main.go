package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/e-flux-platform/sops-secret-scanner/internal/file"
	"github.com/urfave/cli/v2"
)

var (
	baseDir        string
	secretRegexp   string
	secretFilePath string
)

func main() {
	app := &cli.App{
		Version: "0.0.1",
		Name:    "sops-secret-scanner",
		Usage:   "sops-secret-scanner is a SOPS utility which will scan a directory for secret files and encrypt/decrypt them based on the closest .sops.yaml configuration",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "secret-regexp",
				Usage:       "Regular expression to match secret files",
				Value:       `^.+\/secrets?\/.+$`,
				Destination: &secretRegexp,
			},
			&cli.StringFlag{
				Name:        "base-dir",
				Usage:       "Base directory to scan for secret files",
				Value:       ".",
				Destination: &baseDir,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "list-secrets",
				Usage:  "List all files which match the secret-regexp",
				Action: listSecrets,
			},
			{
				Name:   "encrypt-all",
				Usage:  "Encrypt all files in the base directory",
				Action: encryptMany,
			},
			{
				Name:   "decrypt-all",
				Usage:  "Decrypt all files in the base directory",
				Action: decryptMany,
			},
			{
				Name:   "encrypt",
				Usage:  "Encrypt a single file",
				Action: encryptOne,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "secret-file",
						Usage:       "Path to the secret file",
						Aliases:     []string{"f"},
						Required:    true,
						Destination: &secretFilePath,
					},
				},
			},
			{
				Name:   "decrypt",
				Usage:  "Decrypt a single file",
				Action: decryptOne,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "secret-file",
						Usage:       "Path to the secret file",
						Aliases:     []string{"f"},
						Required:    true,
						Destination: &secretFilePath,
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Print(err)
	}
}

func listSecrets(c *cli.Context) error {
	secretFiles, err := file.IdentifySecretFiles(mustResolveAbsoluteFilePath(baseDir), secretRegexp)
	if err != nil {
		return fmt.Errorf("cannot identify secret files in %q: %w", baseDir, err)
	}

	log.Printf("Found %d secret files", len(secretFiles))
	for _, secretFilePath := range secretFiles {
		log.Println(secretFilePath)
	}

	return nil
}

func encryptOne(c *cli.Context) error {
	fileStatus, err := file.Load(mustResolveAbsoluteFilePath(secretFilePath))
	if err != nil {
		return err
	}

	if !fileStatus.Encrypted {
		if err := fileStatus.Encrypt(); err != nil {
			log.Println("failed to encrypt file:", secretFilePath, err)
		}
	} else {
		log.Println("file is not encrypted, skipping...")
	}

	return nil
}

func decryptOne(c *cli.Context) error {
	fileStatus, err := file.Load(mustResolveAbsoluteFilePath(secretFilePath))
	if err != nil {
		return err
	}

	if fileStatus.Encrypted {
		if err := fileStatus.Decrypt(); err != nil {
			log.Println("failed to decrypt file:", secretFilePath, err)
		}
	} else {
		log.Println("file is not encrypted, skipping...")
	}

	return nil
}

func encryptMany(c *cli.Context) error {
	secretFiles, err := file.IdentifySecretFiles(mustResolveAbsoluteFilePath(baseDir), secretRegexp)
	if err != nil {
		return fmt.Errorf("cannot identify secret files in %q: %w", baseDir, err)
	}

	for _, secretFilePath := range secretFiles {
		fileStatus, err := file.Load(secretFilePath)
		if err != nil {
			return err
		}

		if !fileStatus.Encrypted {
			if err := fileStatus.Encrypt(); err != nil {
				log.Println("failed to encrypt file:", secretFilePath, err)
			}
		}
	}

	return nil
}

func decryptMany(c *cli.Context) error {
	secretFiles, err := file.IdentifySecretFiles(mustResolveAbsoluteFilePath(baseDir), secretRegexp)
	if err != nil {
		return fmt.Errorf("cannot identify secret files in %q: %w", baseDir, err)
	}

	for _, secretFilePath := range secretFiles {
		fileStatus, err := file.Load(secretFilePath)
		if err != nil {
			return err
		}

		if fileStatus.Encrypted {
			if err := fileStatus.Decrypt(); err != nil {
				log.Println("failed to decrypt file:", secretFilePath, err)
			}
		}
	}

	return nil
}

func mustResolveAbsoluteFilePath(filePath string) string {
	absoluteFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("failed to resolve absolute file path for %q: %v", filePath, err)
	}

	return absoluteFilePath
}
