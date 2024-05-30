[![Go Report Card](https://goreportcard.com/badge/github.com/multiversx/mx-chain-keys-monitor-go)](https://goreportcard.com/report/github.com/multiversx/mx-chain-keys-monitor-go)
[![Codecov](https://codecov.io/gh/multiversx/mx-chain-keys-monitor-go/branch/main/graph/badge.svg)](https://codecov.io/gh/multiversx/mx-chain-keys-monitor-go)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/multiversx/mx-chain-keys-monitor-go)](https://github.com/multiversx/mx-chain-keys-monitor-go/releases)
[![GitHub](https://img.shields.io/github/license/multiversx/mx-chain-keys-monitor-go)](LICENSE)

# MultiversX keys monitor

This tool allows the monitoring of BLS keys that can participate into the consensus, regardless of the shard they are currently operating in.
This is done by continuously polling the `/validator/statistics` API endpoint route. 
It can monitor one or more networks so one instance of this tool is enough for all keys used on mainnet/testnet or devnet networks.
The monitored BLS keys will be defined as lists. Also, there is the possibility to define just an owner address that staked the BLS keys 
and the application will automatically fetch the registered BLS keys for that identity.

## Feature list
- [x] Keys monitor
    - [x] Monitor any number of networks (mainnet/testnet/devnet) with one instance
    - [x] Monitor any number of BLS keys defined in separate files
    - [x] Automatically fetch the BLS keys staked by an address
    - [x] Threshold definition on each set for the allowed rating drop
    - [x] Configurable polling time for each definition set
- [x] Notification system
    - [x] Integrated the [Pushover](https://pushover.net/) service to allow easy access to push-notifications on mobile devices
    - [x] Integrated the SMTP email service to notify thorough emails the events encountered
    - [X] Integrated the [Telegram bot](https://core.telegram.org/bots) notification service 
- [x] System self-check messages
    - [x] Integrated a self-check system that can periodically send messages on the status of the app
- [x] Scripts & installation support
    - [x] Added scripts for easy setup & upgrade

## Installation

You can choose to run this tool either in a Docker on in a systemd service.

### Docker Setup

You need to have [Docker](https://docs.docker.com/engine/install/) installed on your machine.

Create a directory to save your configs in, and copy the 3 example files required for the tool to run :
```bash
mkdir <config_dir>
cp ./cmd/monitor/config/example/config.toml <config_dir>/config.toml
cp ./cmd/monitor/config/example/credentials.toml <config_dir>/credentials.toml
cp ./cmd/monitor/config/example/network1.list <config_dir>/network1.list
```

Customize your 3 files :
- config.toml is the tool global configuration
- credentials.toml is where you save your secrets 
- network1.list is where you fill your addresses to monitor (⚠️ this file will be located in /network1.list in the container, **make sure to make this path correspond in the config.toml file** ⚠️)

Build the image, use the Dockerfile ;
```bash
sudo docker buildx build -t mx-chain-keys-monitor-go:latest -f Dockerfile .
```

Then, edit the lines 7-9 in `./docker-compose.yml` :
- `<path>/<to>/<your>/config.toml:/config.toml:ro`
- `<path>/<to>/<your>/credentials.toml:/credentials.toml:ro`
- `<path>/<to>/<your>/network1.list:/network1.list:ro`
- You can append more networkX.list files if you want to 

If you want to provide extra arguments to the tool, add them under the `command` object.

You're ready 🚀

```bash
sudo docker compose up -d
```


### Initial setup

Although it's possible, it is not recommended to run the application as `root`. For that reason, a new user is required to be created.
For example, the following script creates a new user called `ubuntu`. This script can be run as `root`.

```bash
# host update/upgrade
apt-get update
apt-get upgrade
apt autoremove

adduser ubuntu
# set a long password
usermod -aG sudo ubuntu
echo 'StrictHostKeyChecking=no' >> /etc/ssh/ssh_config

visudo   
# add this line:
ubuntu  ALL=(ALL) NOPASSWD:ALL
# save & exit

sudo su ubuntu
sudo visudo -f /etc/sudoers.d/myOverrides
# add this line:
ubuntu  ALL=(ALL) NOPASSWD:ALL
# save & exit
```

### Repo clone & scripts init

```bash
cd
git clone https://github.com/multiversx/mx-chain-keys-monitor-go
cd ~/mx-chain-keys-monitor-go/scripts
# the following init call will create ~/mx-chain-keys-monitor-go/scripts/config/local.cfg file
# and will copy the configs from ~/mx-chain-keys-monitor-go/cmd/monitor/config/example to ~/mx-chain-keys-monitor-go/cmd/monitor/config
# to avoid github pull problems
script.sh init
cd config
# edit the local.cfg file for the scripts setup
nano local.cfg
```

### local.cfg configuration

The generated local.cfg file contains the following lines:

```bash
#!/bin/bash
set -e

CUSTOM_HOME=/home/ubuntu
CUSTOM_USER=ubuntu
GITHUBTOKEN=""
MONITOR_EXTRA_FLAGS=""

#Allow user to override the current version of the monitor
OVERRIDE_VER=""
```

The `CUSTOM_HOME` and `CUSTOM_USER` will need to be changed if the current user is not `ubuntu`. 
To easily figure out the current user, the bash command `whoami` can be used.

It is strongly recommended to use a GitHub access token because the scripts consume the GitHub APIs and
throttling might occur without the access token.

The `MONITOR_EXTRA_FLAGS` can contain extra flags to be called whenever the application is started. 
The complete list of the cli command can be found [here](./cmd/monitor/CLI.md) 

The `OVERRIDE_VER` can be used during testing to manually specify an override tag/branch that will be used when building 
the application. If left empty, the upgrade process will automatically fetch and use the latest release.

### Install

After the `local.cfg` configuration step, the scripts can now install the application.
```bash
cd ~/mx-chain-keys-monitor-go/scripts
./script.sh install
```

### Application config

After the application has been installed, it is now time to configure it.
For this, you should edit the `config.toml` and `credentials.toml` files and add files containing lists of BLS keys or addresses in the 
`~/mx-chain-keys-monitor-go/cmd/monitor/config` directory.

The scripts init step already created some initial .toml and .list files to be ready to be used directly.

Configuring the **config.toml** file:

This file contains the general application configuration file.
```toml
[General]
    ApplicationName = "Keys monitoring app"
    # the application can send messages about the internal status at regular intervals
    [General.SystemSelfCheck]
        Enabled = true
        DayOfWeek = "every day" # can also be "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday" and "Sunday"
        Hour = 12 # valid interval 0-23
        Minute = 0 # valid interval 0-59
        PollingIntervalInSec = 30
    [General.Logs]
        LogFileLifeSpanInMB = 1024 # 1GB
        LogFileLifeSpanInSec = 86400 # 1 day

[OutputNotifiers]
    NumRetries = 3
    SecondsBetweenRetries = 10

    # Uses Pushover service that can notify Desktop, Android or iOS devices. Requires a valid subscription.
    # If you enable this notifier, remember to specify the credentials in credentials.toml file
    [OutputNotifiers.Pushover]
        Enabled = false
        URL = "https://api.pushover.net/1/messages.json"

    # SMTP (email) based notification
    # If you enable this notifier, remember to specify the credentials in credentials.toml file
    # If you are using gmail server, please make sure you activate the IMAP server and use App passwords instead of the account's password
    [OutputNotifiers.Smtp]
        Enabled = false
        To = "to@email.com"
        SmtpPort = 587
        SmtpHost = "smtp.gmail.com"

    # Uses Telegram service that can notify Desktop, Android or iOS devices. Requires a running bot and the chat ID for
    # the user that will be notified.
    # If you enable this notifier, remember to specify the credentials in credentials.toml file
    [OutputNotifiers.Telegram]
        Enabled = false
        URL = "https://api.telegram.org"

[[BLSKeysMonitoring]]
    AlarmDeltaRatingDrop = 1.0 # maximum Rating-TempRating value that will trigger an alarm, for the public testnet might use a higher value (2 or 3)
    Name = "network 1"
    ApiURL = "API URL 1"
    ExplorerURL = ""
    PollingIntervalInSeconds = 300  # 5 minutes
    ListFile = "./config/network1.list"
```

* The `General` section

  - The `ApplicationName` can be any string, it will be used in the notification messages.

  - The `General.SystemSelfCheck` section can enable & configure the self check monitoring system
that will output a notification message once a day or week, depending on the configuration parameters.

  - The `General.Logs` will configure the internal logging system (for debugging purposes).


* The `OutputNotifiers` is the section containing the implemented notifiers. 
There are 3 types of notifiers implemented: `Pushover`, `Smtp` and `Telegram`, each with its configuration sections.
The credentials for the notifiers are defined separately in the `credentials.toml` so that is the place where 
passwords or access tokens will be specified.


* The `BLSKeysMonitoring` defines the section used on one network. 

  The `config.toml` file accepts any number of this type of section.
 
  - The `AlarmDeltaRatingDrop` option will specify how large the difference between the "Epoch start rating" - "Current rating" is allowed before 
emitting an alert. Usually, a value of `1` will suffice in most cases. On the public testnet, this value might be increased to `2` or even `3` to 
reduce the number of false positive alarms due to the nature of the other nodes operating on the network.   

  - The `Name` defines a string for the monitored network. Can be something like `Mainnet`, `Testnet`, `Devnet` or any kind of identification string.

  - The `ApiURL` defines the API url for that network. Examples here include `https://api.multiversx.com` for the mainnet, 
`https://testnet-api.multiversx.com` for the testnet, and `https://devnet-api.multiversx.com` for the devnet.

  - The `ExplorerURL` is used whenever a BLS key alert message is emitted, to automatically include the link to that BLS key page.
If left empty, the message will still be emitted, but it will not contain any link.
Examples here include `https://explorer.multiversx.com` for the mainnet, 
`https://testnet-explorer.multiversx.com` for the testnet, and `https://devnet-explorer.multiversx.com` for the devnet.

  - The `PollingIntervalInSeconds` represents the time in seconds between the calls on the API URL. 

  - The `ListFile` will contain the name of the file containing BLS or identity keys. Refer to the example file called 
`network1.list` to check how keys/identities can be defined.

### Application start

After editing the required config files, the application can be started.
```bash
cd ~/mx-chain-keys-monitor-go/scripts
./script.sh start
```

When the application starts, it automatically emits a notification message. This is valid also for the `stop` operation 
issued by the `./script.sh stop` command.

### Backup and upgrade

It is a good practice to save the .toml, .list and the local.cfg files somewhere else just in case the application is cleaned up accidentally.
The upgrade call for the monitor app is done through this command:
```bash
cd ~/mx-chain-keys-monitor-go/scripts
./script.sh upgrade
```

### Uninstalling

The application can be removed by executing the following script:
```bash
cd ~/mx-chain-keys-monitor-go/scripts
./script.sh cleanup
```

### Troubleshooting

If the application fails to start (maybe there is a bad config in the .toml files), the following command can be issued:
```bash
sudo journalctl -f -u mx-chain-keys-monitor.service
```

Also, if the application misbehaves, the logs can be retrieved by using this command:
```bash
cd ~/mx-chain-keys-monitor-go/scripts
./script.sh get_logs
```

If the application crashes and you have followed the installation via Docker, the command to retrieve the logs is as follows:
```bash
sudo docker logs -f mx-chain-keys-monitor-go
```
