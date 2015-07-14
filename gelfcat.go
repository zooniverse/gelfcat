package main

import (
    "github.com/robertkowalski/graylog-golang"
    "github.com/zooniverse/gonx"
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "os"
)

type GenericLogMessage struct {
    Host string `json:"host"`
    Facility string `json:"facility"`
    Message string `json:"message"`
}

func scan_file_unformatted (f *os.File, g *gelf.Gelf, host *string,
        facility *string) {
    scanner := bufio.NewScanner(f)

    for scanner.Scan() {
        m := GenericLogMessage{*host, *facility, scanner.Text()}
        b, _ := json.Marshal(m)
        g.Log(string(b))
    }
}

func scan_file_formatted (f *os.File, g *gelf.Gelf, host *string,
        facility *string, format *string) {
    r := gonx.NewReader(f, *format)
    default_entry := gonx.NewEntry(gonx.Fields{
        "host": *host,
        "facility": *facility,
    })

    for {
        rec, err := r.Read()
        if err == io.EOF {
            break
        }
        rec.Merge(default_entry)
        rec.SetField("message", rec.Raw())
        b, _ := rec.ToJson()
        g.Log(string(b))
    }
}

func scan_file (f *os.File, g *gelf.Gelf, host *string,
        facility *string, format *string) {
    if *format == "" {
        scan_file_unformatted(f, g, host, facility)
    } else {
        scan_file_formatted(f, g, host, facility, format)
    }
}

func main () {
    host := flag.String("host", "localhost",
                        "The name of the host sending the message.")
    facility := flag.String("facility", "gelfcat",
                            "The name of the facility sending the message.")
    server := flag.String("server", "localhost", "The Graylog server name.")
    port := flag.Int("port", 12201, "The Graylog server port.")
    connection := flag.String("connection", "wan",
                              "The Graylog connection type.")
    max_chunk_wan := flag.Int("max-chunk-wan", 1420,
                              "The maximum chunk size for WAN connections.")
    max_chunk_lan := flag.Int("max-chunk-lan", 8154,
                              "The maximum chunk size for LAN connections.")
    format := flag.String("format", "",
                          "The format to use when parsing entries for fields.")

    flag.Parse()
    file_names := flag.Args()

    g := gelf.New(gelf.Config{
        GraylogHostname: *server,
        GraylogPort: *port,
        Connection: *connection,
        MaxChunkSizeWan: *max_chunk_wan,
        MaxChunkSizeLan: *max_chunk_lan,
    })

    if len(file_names) > 0 {
        for _, file_name := range file_names {
            f, err := os.Open(file_name)
            defer f.Close()

            if err != nil {
                fmt.Fprintf(os.Stderr, "Error: Could not open file %s\n",
                            file_name)
                continue
            }
            scan_file(f, g, host, facility, format)
        }
    } else {
        scan_file(os.Stdin, g, host, facility, format)
    }
}
