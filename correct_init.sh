#!/usr/bin/env bash
rm -rf ./testdata
mkdir -p ./testdata

## Create TCG rsa SRK compliant key
echo -n -e $(printf '\\x00%.0s' {1..256}) | tpm2 createprimary --hierarchy o --hash-algorithm sha256 --key-algorithm rsa2048:null:aes128cfb \
    --attributes 'decrypt|restricted|fixedtpm|fixedparent|sensitivedataorigin|userwithauth|noda' \
    --key-context ./testdata/rsa_srk.ctx --format pem --output ./testdata/rsa_srk.pub \
    --unique-data -

## Create TCG ecc SRK compliant key
echo -n -e $(printf '\\x00%.0s' {1..64}) | tpm2 createprimary --hierarchy o --hash-algorithm sha256 --key-algorithm ecc256:aes128cfb \
    --attributes 'decrypt|restricted|fixedtpm|fixedparent|sensitivedataorigin|userwithauth|noda' \
    --key-context ./testdata/ecc_srk.ctx --format pem --output ./testdata/ecc_srk.pub \
    --unique-data -
