# gelfcat
A command-line utility to send files to Graylog. Usage:

    gelfcat [options] [file1 file2 ...]
        -connection="wan": The Graylog connection type.
        -facility="gelfcat": The name of the facility sending the message
        -format: [gonx format](https://github.com/satyrius/gonx#format) or "json"
        -host="localhost": The name of the host sending the message
        -max-chunk-lan=8154: The maximum chunk size for LAN connections.
        -max-chunk-wan=1420: The maximum chunk size for WAN connections.
        -port=12201: The Graylog server port.
        -server="localhost": The Graylog server name.

If no files are specified, input is read from stdin.
