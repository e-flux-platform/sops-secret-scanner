package file

import (
	"fmt"
	"log"
	"os"

	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/keyservice"
)

func (fs *File) Decrypt() error {
	if !fs.Encrypted {
		return nil
	}

	_, err := common.DecryptTree(common.DecryptTreeOpts{
		Cipher: fs.cipher,
		KeyServices: []keyservice.KeyServiceClient{
			keyservice.NewLocalClient(),
		},
		Tree: fs.tree,
	})
	if err != nil {
		return fmt.Errorf("cannot decrypt file %s: %w", fs.tree.FilePath, err)
	}

	decryptedFile, err := fs.store.EmitPlainFile(fs.tree.Branches)
	if err != nil {
		return err
	}

	f, err := os.Create(fs.tree.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Decrypting file: %s\n", fs.tree.FilePath)

	_, err = f.Write(decryptedFile)

	return err
}
