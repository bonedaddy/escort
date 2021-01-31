package main

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
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
				buffer := new(bytes.Buffer)
				writer, err := flate.NewWriter(buffer, flate.BestCompression)
				if err != nil {
					return err
				}
				if _, err := writer.Write([]byte(c.String("input"))); err != nil {
					return err
				}

				if err := writer.Close(); err != nil {
					return err
				}

				fmt.Println(base64.StdEncoding.EncodeToString(buffer.Bytes()))
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
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
