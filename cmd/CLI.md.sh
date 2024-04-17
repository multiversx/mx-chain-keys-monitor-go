#!/usr/bin/env bash

generate() {
    generateForMonitor
}

generateForMonitor() {
    HELP="
# Monitor CLI

The **MultiversX Keys Monitor** exposes the following Command Line Interface:
$(code)
\$ node --help

$(./monitor/monitor --help | head -n -3)
$(code)
"
    echo "$HELP" > ./monitor/CLI.md
}

code() {
    printf "\n\`\`\`\n"
}

generate
