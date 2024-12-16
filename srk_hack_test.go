package main

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"os"
	"path"
	"testing"

	"github.com/loicsikidi/srk-interoperability-issue/internal/testutil"
	"go.step.sm/crypto/pemutil"

	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/google/go-tpm/tpmutil"
	"github.com/stretchr/testify/require"
)

func TestSrkWithHack(t *testing.T) {
	rwc, err := testutil.OpenTPM20()
	require.Nil(t, err)

	t.Run("ensure tpm2-tools produced rsa SRK compliant key", func(t *testing.T) {
		var nilModulus tpmutil.U16Bytes
		hackDefaultRsaSrkTemplate := defaultRsaSrkTemplate
		hackDefaultRsaSrkTemplate.RSAParameters.ModulusRaw = nilModulus
		_, wantPub, err := tpm2.CreatePrimary(rwc, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", hackDefaultRsaSrkTemplate)
		require.Nil(t, err)

		b, err := os.ReadFile(path.Join(testdataDir, "rsa_srk.pub"))
		require.Nil(t, err)

		gotPub, err := pemutil.Parse(b)
		require.Nil(t, err)

		want, ok := wantPub.(*rsa.PublicKey)
		require.True(t, ok)

		got, ok := gotPub.(*rsa.PublicKey)
		require.True(t, ok)

		require.True(t, want.Equal(got))
	})

	t.Run("ensure tpm2-tools produced ecc SRK compliant key", func(t *testing.T) {
		var nilEccPoint tpm2.ECPoint
		hackDefaultEccSrkTemplate := defaultEccSrkTemplate
		hackDefaultEccSrkTemplate.ECCParameters.Point = nilEccPoint
		_, wantPub, err := tpm2.CreatePrimary(rwc, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", hackDefaultEccSrkTemplate)
		require.Nil(t, err)

		b, err := os.ReadFile(path.Join(testdataDir, "ecc_srk.pub"))
		require.Nil(t, err)

		gotPub, err := pemutil.Parse(b)
		require.Nil(t, err)

		want, ok := wantPub.(*ecdsa.PublicKey)
		require.True(t, ok)

		got, ok := gotPub.(*ecdsa.PublicKey)
		require.True(t, ok)

		require.True(t, want.Equal(got))
	})
}
