package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/archway-network/augusta-testnet-signer/types"
	cosmosFlag "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var signIDCmd = &cobra.Command{
	Use:   "sign_id [key_name]",
	Short: "Sign an id message with your key to prove that you own the key",
	Long: `This command asks you some personal details like legal name, email address and github handle, 
           and uses the key associated with [key_name] to sign the message containing that data and print it`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyringBackend, err := cmd.Flags().GetString(cosmosFlag.FlagKeyringBackend)
		if err != nil {
			return err
		}

		backendDir, err := cmd.Flags().GetString(cosmosFlag.FlagKeyringDir)
		if err != nil {
			return err
		}

		kr, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, backendDir, os.Stdin)
		if err != nil {
			return err
		}

		keyName := args[0]
		_, err = kr.Key(keyName)
		if err != nil {
			return err
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Please Enter your Full Legal Name:")
		fullLegalName, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fullLegalName = strings.TrimSpace(fullLegalName)

		fmt.Printf("Please enter your github handle:")
		githubHandle, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		githubHandle = strings.TrimSpace(githubHandle)

		fmt.Printf("Please enter your email address (Use same email address, you would use in kyc form):")
		emailAddress, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		emailAddress = strings.TrimSpace(emailAddress)

		container, err := types.CreateContainer(fullLegalName, githubHandle, emailAddress, keyName, kr)
		if err != nil {
			return err
		}

		fmt.Println("Your Augusta incentivized testnet address is: ", container.ID.AccountAddress)
		fmt.Println("Amino encoded Public key is:", container.ID.PubKey)

		marshalledContainer, err := json.MarshalIndent(container, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println("Please submit json below the line in the form")
		fmt.Println("-----------------------------")
		fmt.Println(string(marshalledContainer))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(signIDCmd)
}
