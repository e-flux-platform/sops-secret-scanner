package main

import (
	"fmt"
	"log"
	"os"

	"github.com/e-flux-platform/sops-secret-scanner/internal/file"
	"github.com/urfave/cli/v2"
)

var (
	baseDir      string
	secretRegexp string
)

func main() {
	app := &cli.App{
		Version: "0.0.1",
		Name:    "sops-ecret-scanner",
		Usage:   "sop-ssecret-scanner is a SOPS utility which will scan a directory for secret files and encrypt/decrypt them based on the .sops.yaml",
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
				Value:       "deployment/",
				Destination: &baseDir,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "encrypt-all",
				Description: "Encrypt all files in the base directory",
				Action:      encrypt,
			},
			{
				Name:        "decrypt-all",
				Description: "Decrypt all files in the base directory",
				Action:      decrypt,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Panicf("application exited with error: %s", err)
	}
}

func encrypt(c *cli.Context) error {
	secretFiles, err := file.IdentifySecretFiles(baseDir, secretRegexp)
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
				fmt.Println("failed to encrypt file:", secretFilePath, err)
			}
		}
	}

	return nil
}

func decrypt(c *cli.Context) error {
	secretFiles, err := file.IdentifySecretFiles(baseDir, secretRegexp)
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
				fmt.Println("failed to decrypt file:", secretFilePath, err)
			}
		}
	}

	return nil
}
