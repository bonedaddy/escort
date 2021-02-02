package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "escort"
	app.Commands = cli.Commands{
		&cli.Command{
			Name:  "compress",
			Usage: "compress data using DEFLATE and base64 encode it",
			Action: func(c *cli.Context) error {
				var data string
				if c.String("input") != "" {
					data = c.String("input")
				} else if c.String("input.file") != "" {
					dBytes, err := ioutil.ReadFile(c.String("input.file"))
					if err != nil {
						return err
					}
					data = string(dBytes)
				} else {
					return errors.New("input.file and input are nil")
				}
				buffer := new(bytes.Buffer)
				writer, err := flate.NewWriter(buffer, flate.BestCompression)
				if err != nil {
					return err
				}
				if _, err := writer.Write([]byte(data)); err != nil {
					return err
				}

				if err := writer.Close(); err != nil {
					return err
				}
				parts := Chunks(base64.StdEncoding.EncodeToString(buffer.Bytes()), 250)
				for i, part := range parts {
					fmt.Printf("%v|%s\n", i, part)
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "base64",
			Usage: "base64 encode/decode util commands",
			Subcommands: cli.Commands{
				&cli.Command{
					Name:  "encode",
					Usage: "base64 encode input data",
					Action: func(c *cli.Context) error {
						fmt.Println(base64.StdEncoding.EncodeToString([]byte(c.String("input"))))
						return nil
					},
				},
				&cli.Command{
					Name:  "decode",
					Usage: "base64 decode input data",
					Action: func(c *cli.Context) error {
						data, err := base64.StdEncoding.DecodeString(c.String("input"))
						if err != nil {
							return err
						}
						fmt.Println(string(data))
						return nil
					},
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "input",
			Usage: "input data to compress",
		},
		&cli.StringFlag{
			Name:  "input.file",
			Usage: "input data to compress from a file",
		},
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

// Chunks is used to split a string into segments of chunkSize
// https://stackoverflow.com/questions/25686109/split-string-by-length-in-golang
func Chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}
