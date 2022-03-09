
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
  "log"
	"regexp"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var derivationPath = hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")

var (
	alphabet   = regexp.MustCompile("^[0-9a-f]*$")
	numWorkers = runtime.NumCPU()
)

// Wallet stores private key and address containing desired substring at Index
type Wallet struct {
	Address    string
	Mnemonic string
}

func main() {
	var one bool
	var prefix, suffix string
	flag.BoolVar(&one, "one", false, "Stop after finding first address")
	flag.StringVar(&prefix, "p", "", "Public address prefix")
	flag.StringVar(&suffix, "s", "", "Public address suffix")
	flag.Parse()
	if prefix == "" && suffix == "" {
		fmt.Printf(`
This tool generates Ethereum public and private keypair until it finds address
which contains required prefix and/or suffix.
Address part can contain only digits and letters from A to F.
For fast results suggested length of sum of preffix and suffix is 4-6 characters.
If you want more, be patient.
Usage:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}
	if !alphabet.MatchString(prefix) {
		fmt.Println("Prefix must match the alphabet:", alphabet.String())
		os.Exit(2)
	}
	if !alphabet.MatchString(suffix) {
		fmt.Println("Suffix must match the alphabet:", alphabet.String())
		os.Exit(3)
	}
	walletChan := make(chan Wallet)
	for i := 0; i < numWorkers; i++ {
		go generateWallet(prefix, suffix, walletChan)
	}
	for w := range walletChan {
		fmt.Printf(
			"Address: %s Mnemonic: %s\n",
			w.Address,
			w.Mnemonic)
		if one {
			break
		}
	}
}

func generateWallet(prefix, suffix string, walletChan chan Wallet) {
	for {

    mnemonic, err := hdwallet.NewMnemonic(256)
	if err != nil {
	log.Fatal(err)
	}
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
	log.Fatal(err)
	}
	account, err := wallet.Derive(derivationPath, false)
	if err != nil {
	log.Fatal(err)
	}

	addressHex := account.Address.Hex()
 
		if prefix != "" && !strings.HasPrefix(addressHex[2:], prefix) {
			continue
		}
		//if suffix != "" && !strings.HasSuffix(addressHex, suffix) {
		//	continue
		//}
		walletChan <- Wallet{
			Address:  addressHex,
			Mnemonic: mnemonic,
		}
	}
}