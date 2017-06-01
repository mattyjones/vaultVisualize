## Configuration

If the variable is not set on the commandline then read it in via viper. In some cases Vault will also be able to
read in the value directly. The order of prescendence is:

  1. --token flag will be used
  1. Viper will read in the environment variable PREFIX_TOKEN (PREFIX == VAULT, see root.go)
  1. Viper will read in the configuration from *./vaultVisualize.yaml*
  1. Vault will read in the environment variable VAULT_TOKEN (hardcoded)

The order is top down with each item taking precedence over the item below it.

### Debuging

*--debug* will print out various useful data points if things are not working. This flag will cause the app to exit 
after the debug statement is printed, no data will be fetched.
- the full url being used
- critical variable values at execution time
- if any env variables or the automatic config are being used (command line variables may be used in conjunction)

### Logging

All errors will be logged to syslog in json format and to stdout in text format for easy reading


