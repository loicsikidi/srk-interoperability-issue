{pkgs ? import <nixpkgs> {}}:
with pkgs;
  mkShell {
    packages = [
      go
      delve

      tpm2-tools
    ];
    # we disable the hardening due to this error: https://github.com/tpm2-software/tpm2-tools/issues/1561
    # fix found here: https://github.com/NixOS/nixpkgs/issues/18995#issuecomment-249748307
    hardeningDisable = ["fortify"];
  }
