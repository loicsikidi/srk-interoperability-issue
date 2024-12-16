package main

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"os"
	"path"
	"testing"

	"github.com/loicsikidi/srk-interoperability-issue/internal/testutil"
	"go.step.sm/crypto/pemutil"

	"github.com/google/go-tpm-tools/client"
	"github.com/google/go-tpm/legacy/tpm2"
	"github.com/stretchr/testify/require"
)

const (
	testdataDir = "./testdata"
)

var (
	defaultRsaSrkTemplate = client.SRKTemplateRSA()
	defaultEccSrkTemplate = client.SRKTemplateECC()
)

func TestSrk(t *testing.T) {
	rwc, err := testutil.OpenTPM20()
	require.Nil(t, err)

	t.Run("ensure tpm2-tools produced rsa SRK compliant key", func(t *testing.T) {
		_, wantPub, err := tpm2.CreatePrimary(rwc, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", defaultRsaSrkTemplate)
		require.Nil(t, err)

		b, err := os.ReadFile(path.Join(testdataDir, "rsa_srk.pub"))
		require.Nil(t, err)

		gotPub, err := pemutil.Parse(b)
		require.Nil(t, err)

		want, ok := wantPub.(*rsa.PublicKey)
		require.True(t, ok)

		got, ok := gotPub.(*rsa.PublicKey)
		require.True(t, ok)

		require.True(t, want.Equal(got), "SRK should match")
	})

	t.Run("ensure tpm2-tools produced ecc SRK compliant key", func(t *testing.T) {
		_, wantPub, err := tpm2.CreatePrimary(rwc, tpm2.HandleOwner, tpm2.PCRSelection{}, "", "", defaultEccSrkTemplate)
		require.Nil(t, err)

		b, err := os.ReadFile(path.Join(testdataDir, "ecc_srk.pub"))
		require.Nil(t, err)

		gotPub, err := pemutil.Parse(b)
		require.Nil(t, err)

		want, ok := wantPub.(*ecdsa.PublicKey)
		require.True(t, ok)

		got, ok := gotPub.(*ecdsa.PublicKey)
		require.True(t, ok)

		require.True(t, want.Equal(got), "SRK should match")
	})
}
