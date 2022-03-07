package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bonedaddy/escort/pkg"
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
				var (
					data []byte
					err  error
				)
				if c.String("input") != "" {
					data = []byte(c.String("input"))
				} else if c.String("input.file") != "" {
					data, err = ioutil.ReadFile(c.String("input.file"))
					if err != nil {
						return err
					}
				} else {
					return errors.New("input.file and input are nil")
				}
				core := pkg.NewCore(nil, "1", 250)
				parts, err := core.Trick(data)
				if err != nil {
					return err
				}
				for _, part := range parts {
					fmt.Println(part)
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
