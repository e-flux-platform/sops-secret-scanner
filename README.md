# sops-secret-scanner

This is a command line utility built on top of [Mozilla SOPS](https://github.com/mozilla/sops) which which encrypts all files in a `secrets` directory that is a child of the `base-dir`.

The motivation for this was the ability to create a `pre-commit` hook which can capture and encrypt potential secrets before they are pushed to a remote repository.

```
NAME:
   ssops-ecret-scanner - sop-ssecret-scanner is a SOPS utility which will scan a directory for secret files and encrypt/decrypt them based on the .sops.yaml.

USAGE:
   ssops-ecret-scanner [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   encrypt-all
   decrypt-all
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --secret-regexp value  (default: "^.+\\/secrets?\\/.+$")
   --base-dir value       (default: ".")
   --help, -h             show help
   --version, -v          print the version
```
