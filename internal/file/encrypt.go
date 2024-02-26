package file

import (
	"fmt"
	"log"
	"os"

	"github.com/getsops/sops/v3/cmd/sops/common"
)

func (fs *File) Encrypt() error {
	if fs.Encrypted {
		return nil
	}

	dataKey, errs := fs.tree.GenerateDataKey()
	if len(errs) > 0 {
		return fmt.Errorf("Could not generate data key: %s", errs)
	}

	err := common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    fs.tree,
		Cipher:  fs.cipher,
	})
	if err != nil {
		return err
	}

	encryptedFile, err := fs.store.EmitEncryptedFile(*fs.tree)
	if err != nil {
		return err
	}

	f, err := os.Create(fs.tree.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Encrypting file: %s\n", fs.tree.FilePath)

	_, err = f.Write(encryptedFile)

	return err
}
