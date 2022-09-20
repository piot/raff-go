/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"errors"
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	raff "github.com/piot/raff-go/src"
	"io"
	"log"
	"os"
)

type Options struct {
	View ViewCmd `cmd:"" help:"views a raff file"`
}

type ViewCmd struct {
	Path      string `help:"path to file" arg:"" default:"." type:"file"`
	Verbosity int    `help:"verbose output" type:"counter" short:"v"`
}

func (c *ViewCmd) Run() error {
	file, err := os.Open(c.Path)
	if err != nil {
		return err
	}

	if err := raff.ReadHeader(file); err != nil {
		return err
	}

	nameColor := color.New(color.FgBlue)
	octetCountColor := color.New(color.FgHiMagenta)

	for {
		header, err := raff.ReadChunkHeader(file)
		if errors.Is(err, io.EOF) {
			log.Printf("Encountered end of file")

			break
		}

		if err != nil {
			return fmt.Errorf("could not read chunk header %w", err)
		}

		if header.Icon == 0 && header.Name == 0 {
			break
		}

		nameWithColor := nameColor.Sprint(raff.NameToString(header.Name))
		octetCountWithColor := octetCountColor.Sprintf("%d", header.OctetCount)

		fmt.Printf("%2s %4s %s\n", raff.IconToString(header.Icon), nameWithColor, octetCountWithColor)

		if _, err := file.Seek(int64(header.OctetCount), os.SEEK_CUR); err != nil {
			return fmt.Errorf("could not seek to next chunk %w", err)
		}
	}

	return nil
}

func main() {
	color.NoColor = false
	ctx := kong.Parse(&Options{})

	if err := ctx.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
