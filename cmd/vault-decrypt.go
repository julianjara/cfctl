package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/liangrog/cfctl/pkg/utils"
	"github.com/liangrog/cfctl/pkg/utils/i18n"
	"github.com/liangrog/cfctl/pkg/utils/templates"
	"github.com/liangrog/vault"
	"github.com/spf13/cobra"
)

var (
	vaultDecryptShort = i18n.T("Decrypt vault encrypted files.")

	vaultDecryptLong = templates.LongDesc(i18n.T(`
		Decrypt vault encrypted files. 'CFCTL_VAULT_PASSWORD'
		and 'CFCTL_VAULT_PASSWORD_FILE' environment variables 
		can be used to replace '--vault-password' and 
		'--vault-password-file' flags.`))

	vaultDecryptExample = templates.Examples(i18n.T(`
		# Decrypt multiple files
		$ cfctl vault decrypt file1 file2 file3

		# Decrypt using password file
		$ cfctl vault decrypt file1 --vault-password-file path/to/file`))
)

// Register sub commands
func init() {
	cmd := getCmdVaultDecrypt()

	CmdVault.AddCommand(cmd)
}

// cmd: encrypt
func getCmdVaultDecrypt() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrypt",
		Short:   vaultDecryptShort,
		Long:    vaultDecryptLong,
		Example: fmt.Sprintf(vaultDecryptExample),
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("Missing file name in command argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := batchDecrypt(
				cmd.Flags().Lookup(CMD_VAULT_PASSWORD).Value.String(),
				cmd.Flags().Lookup(CMD_VAULT_PASSWORD_FILE).Value.String(),
				args,
			)

			silenceUsageOnError(cmd, err)

			return err

		},
	}

	return cmd
}

func batchDecrypt(pss, pssFile string, files []string) error {
	passwords, err := GetPasswords(pss, pssFile, false, false)
	if err != nil {
		return err
	}

	result := make(chan error, 10)
	for _, file := range files {
		go func(file string, pass []string, res chan<- error) {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				res <- err
				return
			}

			// Try every given password
			decrypted := false
			var output []byte
			for _, p := range pass {
				output, err = vault.Decrypt(p, data)
				if err == nil {
					decrypted = true
					break
				}
			}

			if decrypted && len(output) > 0 {
				if err := ioutil.WriteFile(file, output, 0644); err != nil {
					res <- err
				}
			} else {
				res <- errors.New(fmt.Sprintf("Failed to decrypt %s using all given password", file))
			}

			res <- nil
		}(file, passwords, result)
	}

	for j := 0; j < len(files); j++ {
		err := <-result
		if err != nil {
			if err := utils.Print("", files[j], err); err != nil {
				return err
			}
		}
	}

	return nil
}
