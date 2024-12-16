# [*Proof Of Concept*] SRK interoperability issues between `tpm2-tools` and `go-tpm-tools`

*Context: when I create a SRK key with [tpm2-tools](https://github.com/tpm2-software/tpm2-tools/tree/master) (cf. `tpm2 createPrimary`) and [go-tpm](https://github.com/google/go-tpm) there are differents while I'm expecting to get the same...*

## Prequisites

1. Install [tpm2-tools](https://tpm2-tools.readthedocs.io/en/latest/INSTALL/)
1. Install [golang](https://go.dev/doc/install#) >= 1.22.0

> [!TIP]
> If you are using [Nix](https://nixos.org/), just run `nix-shell` command in the root directory

## How-to reproduce

1. Produce ECC & RSA SKR with `tpm2-tools`

    ```bash
    ./init.sh
    ```

> [!NOTE]
> The script assume that you don't require `sudo` to interact with the TPM

> [!IMPORTANT]
> The command can take up to **~15sec** to finished (due to RSA 2048 creation)

2. Compare with SKR produced by `go-tpm`

    ```bash
    go test -timeout 30s -run ^TestSrk$ github.com/loicsikidi/srk-interoperability-issue
    ```

> [!IMPORTANT]
> The command can take up to **~15sec** to finished (due to RSA 2048 creation)


You should see the message below demonstrating the mismatch:

```
--- FAIL: TestSrk (2.54s)
    --- FAIL: TestSrk/ensure_tpm2-tools_produced_rsa_SRK_compliant_key (2.38s)
        srk_test.go:47: 
                Error Trace:    redacted.../srk-interoperability-issue/srk_test.go:47
                Error:          Should be true
                Test:           TestSrk/ensure_tpm2-tools_produced_rsa_SRK_compliant_key
                Messages:       SRK should match
    --- FAIL: TestSrk/ensure_tpm2-tools_produced_ecc_SRK_compliant_key (0.15s)
        srk_test.go:66: 
                Error Trace:    redacted.../srk-interoperability-issue/srk_test.go:66
                Error:          Should be true
                Test:           TestSrk/ensure_tpm2-tools_produced_ecc_SRK_compliant_key
                Messages:       SRK should match
FAIL
FAIL    github.com/loicsikidi/srk-interoperability-issue        2.548s
FAIL
```

## Root cause?

After a thorough investigation, it seems that the mismatch is due to a difference between templates used in `tpm2-tools` and `go-tpm-tools`.

According to *TCG Credential Profile EK 2.0* `go-tpm-tools` filled `public.unique.buffer` with an array of 0.


| Key type | TCG Spec | *go-tpm-tools* implementation |
| -------- | -------- | ------------------------------|
| RSA 2048 | [here](https://trustedcomputinggroup.org/wp-content/uploads/TCG-EK-Credential-Profile-V-2.5-R2_published.pdf#page=38) | [code](https://github.com/google/go-tpm-tools/blob/main/client/template.go#L46) |
| ECC P256 | [here](https://trustedcomputinggroup.org/wp-content/uploads/TCG-EK-Credential-Profile-V-2.5-R2_published.pdf#page=39) | [code](https://github.com/google/go-tpm-tools/blob/main/client/template.go#L54-L57) 

I don't know how to do the same with `tpm2-tools`...

### Demonstration

If I remove `public.unique.buffer` from `go-tpm-tools` it matches the SRK produceds by `tpm2-tools`.

```bash
go test -timeout 30s -run ^TestSrkWithHack$ github.com/loicsikidi/srk-interoperability-issue
```

> [!IMPORTANT]
> The command can take up to **~15sec** to finished (due to RSA 2048 creation)

You should see that the test case works: 

```
ok      github.com/loicsikidi/srk-interoperability-issue
```