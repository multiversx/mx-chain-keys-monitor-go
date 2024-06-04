
# Monitor CLI

The **MultiversX Keys Monitor** exposes the following Command Line Interface:

```
$ node --help

NAME:
   MultiversX keys monitor tool - This is the entry point for starting a new MultiversX keys monitor
USAGE:
   monitor [global options]
   
AUTHOR:
   The MultiversX Team <contact@multiversx.com>
   
GLOBAL OPTIONS:
   --config value        The main configuration file (default: "./config/config.toml")
   --credentials value   The credentials configuration file (default: "./config/credentials.toml")
   --log-level level(s)  This flag specifies the logger level(s). It can contain multiple comma-separated value. For example, if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG log level. (default: "*:INFO ")
   --log-save            Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.
   --test-notifiers      Boolean option testing out the notifiers. The application will send a message to all configured notifiers and will close.
   --help, -h            show help
   --version, -v         print the version
   

```

