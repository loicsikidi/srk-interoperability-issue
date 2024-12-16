package testutil

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/google/go-tpm/legacy/tpm2"
)

const (
	tpmRoot = "/sys/class/tpm"
)

// OpenTPM20 initializes access to the TPM 2.0.
func OpenTPM20() (io.ReadWriteCloser, error) {
	candidateTPMs, err := probeSystemTPMs20()
	if err != nil {
		return nil, err
	}

	for _, tpmPath := range candidateTPMs {
		return openTPM20Path(tpmPath)
	}

	return nil, fmt.Errorf("no TPM 2.0 found")
}

// probeSystemTPMs20 returns paths of TPM 2.0 devices available.
func probeSystemTPMs20() ([]string, error) {
	var paths []string

	tpmDevs, err := os.ReadDir(tpmRoot)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if os.IsNotExist(err) {
		return nil, nil
	}

	for _, tpmDev := range tpmDevs {
		if strings.HasPrefix(tpmDev.Name(), "tpm") {
			tpmPath := path.Join(tpmRoot, tpmDev.Name())
			// Use the "/sys/class/tpm/tpmX/device/caps" path to determine the TPM
			// version. TPM 2.0 devices do not create this path.
			if _, err := os.Stat(path.Join(tpmPath, "caps")); err != nil {
				if !os.IsNotExist(err) {
					return nil, err
				}
				// add only TPM 2.0 paths
				paths = append(paths, tpmPath)
			}
		}
	}

	return paths, nil
}

// openTPM20Path attempts to open a TPM 2.0 device.
func openTPM20Path(tpmPath string) (io.ReadWriteCloser, error) {
	// If the TPM has a kernel-provided resource manager, we should
	// use that instead of communicating directly.
	devPath := path.Join("/dev", path.Base(tpmPath))
	f, err := os.ReadDir(path.Join(tpmPath, "device", "tpmrm"))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		// no resource-managed TPM device found, use the direct TPM access
	} else if len(f) > 0 {
		devPath = path.Join("/dev", f[0].Name())
	}

	return tpm2.OpenTPM(devPath)
}
