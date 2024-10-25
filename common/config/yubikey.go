//go:build yubikey

package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strconv"

	"github.com/go-piv/piv-go/v2/piv"
	"github.com/knadh/koanf/v2"
	"github.com/ovh/okms-cli/common/utils/exit"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	RegisterAuthMethod("yubikey", newYubikeyAuth)
}

type yubikeyAuth struct {
	cert *x509.Certificate
	slot piv.Slot
}

func newYubikeyAuth(cmd *cobra.Command, k *koanf.Koanf, envPrefix string) EndpointAuth {
	yk := &yubikeyAuth{
		slot: piv.SlotAuthentication, // Default slot is 9a
	}
	if slotId := GetString(k, "slot", envPrefix+"PIV_SLOT", nil); slotId != "" {
		var ok bool
		if yk.slot, ok = parseSlotID(slotId); !ok {
			exit.OnErr(errors.New("Invalid yubikey slot ID"))
		}
	}

	if certFile := GetString(k, "cert", envPrefix+"_CERT", cmd.Flags().Lookup("cert")); certFile != "" {
		pemBlock, _ := pem.Decode(exit.OnErr2(os.ReadFile(certFile)))
		if pemBlock.Type != "CERTIFICATE" {
			exit.OnErr(errors.New("Invalid certificate"))
		}
		yk.cert = exit.OnErr2(x509.ParseCertificate(pemBlock.Bytes))
	}
	return yk
}

func (epCfg *yubikeyAuth) TlsCertificates() []tls.Certificate {
	cards := exit.OnErr2(piv.Cards())
	if len(cards) == 0 {
		exit.OnErr(errors.New("No yubi key found"))
	}
	// Let's assume the first key is the yubi key
	yubiKey := exit.OnErr2(piv.Open(cards[0]))

	// FIXME: Where / How to close the connection ?
	// The issue preventing us from deferring the close is that we
	// return a tls.Certificate thta have a remote private key / signer that requires
	// the connection to stay open. Anyway, the connection will be closed on program exit, so we may not
	// need to close it by ourself. However it would be cleaner.
	// defer yubiKey.Close()
	cert := epCfg.cert
	if cert == nil {
		// Retrieve the x509 certificate from the authentication slot
		// unless it is provided externally by the user
		cert = exit.OnErr2(yubiKey.Certificate(epCfg.slot))
	}
	pub := cert.PublicKey

	// Get the crypto.PrivateKey associated with the certificate. The private key is always in the yubi key
	// and can never be extracted.
	priv := exit.OnErr2(yubiKey.PrivateKey(epCfg.slot, pub, piv.KeyAuth{
		// Will ask for the pin when needed
		PINPrompt: func() (pin string, err error) {
			password, err := pterm.DefaultInteractiveTextInput.WithMask("*").WithOnInterruptFunc(func() {
				exit.OnErr(errors.New("Interrupted"))
			}).Show("Enter HW token PIN")
			println("Press your HW token security button if needed")
			return password, err
		},
	}))

	return []tls.Certificate{
		{
			Certificate: [][]byte{cert.Raw},
			PrivateKey:  priv,
			Leaf:        cert,
		},
	}
}

func parseSlotID(slotID string) (piv.Slot, bool) {
	key, err := strconv.ParseUint(slotID, 16, 32)
	if err != nil {
		return piv.Slot{}, false
	}

	switch uint32(key) {
	case piv.SlotAuthentication.Key:
		return piv.SlotAuthentication, true
	case piv.SlotSignature.Key:
		return piv.SlotSignature, true
	case piv.SlotCardAuthentication.Key:
		return piv.SlotCardAuthentication, true
	case piv.SlotKeyManagement.Key:
		return piv.SlotKeyManagement, true
	}

	return piv.RetiredKeyManagementSlot(uint32(key))
}
