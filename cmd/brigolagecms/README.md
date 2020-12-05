# BrigolageCMS
BrigolageCMS is a web server.

# Configuration

Note that the values with type of `duration` are parsed by Go's [`time.ParseDuration`](https://golang.org/pkg/time/#ParseDuration).
They should have values like `5s` for five seconds.

## Command Line Options
All command line options override values from the environment and any configuration file values. 

## Environment Variables
The server will load values from the environment.
Any environment value will override the same value from a configuration file.

Generally the variable names are prefixed with `BRIGOLAGECMS_`.

## Configuration File
The server will load a configuration file if requested.
The file must contain a single valid JSON object.

Note that command line options and environment variables will override the
values in the configuration file.

### Example File
The `config.json` file is a example configuration file.

## Default values
We recommend running the program with the `--help` flag on the command line to discover the default configuration values. 