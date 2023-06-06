# sops-secret-scanner

This is a command line utility built on top of [Mozilla SOPS](https://github.com/mozilla/sops) which which encrypts all files in a `secrets` directory that is a child of the `base-dir`.

The motivation for this was the ability to create a `pre-commit` hook which can capture and encrypt potential secrets before they are pushed to a remote repository.
```
NAME:
   sops-secret-scanner - sops-secret-scanner is a SOPS utility which will scan a directory for secret files and encrypt/decrypt them based on the closest .sops.yaml configuration

USAGE:
   sops-secret-scanner [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   list-secrets  List all files which match the secret-regexp
   encrypt-all   Encrypt all files in the base directory
   decrypt-all   Decrypt all files in the base directory
   encrypt       Encrypt a single file
   decrypt       Decrypt a single file
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --secret-regexp value  Regular expression to match secret files (default: "^.+\\/secrets?\\/.+$")
   --base-dir value       Base directory to scan for secret files (default: ".")
   --help, -h             show help
   --version, -v          print the version
```

### List all secret files in currenct directory

The default configuration will use the currenct directory as the `base-dir` meaning you can exclude it if you're only interested in finding secrets that exist in the currenct directory and its children.

```
sops-secret-scanner list-secrets
```

### Encrypt a file

You can provide a relative or an absolute path.

```
sops-secret-scanner encrypt -f {path_to_file}
```


### Decrypt a file

You can provide a relative or an absolute path.

```
sops-secret-scanner decrypt -f {path_to_file}
```

### Decrypt all files

```
sops-secret-scanner --base-dir {dir} decrypt-all
```

### Encrypt all files

```
sops-secret-scanner --base-dir {dir} encrypt-all
```
