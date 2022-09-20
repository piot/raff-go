/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
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
		if err == io.EOF {
			log.Printf("Encountered end of file")
			break
		}
		if err != nil {
			return err
		}

		if header.Icon == 0 && header.Name == 0 {
			log.Printf("Found proper end of file chunk")
			break
		}

		nameWithColor := nameColor.Sprint(raff.NameToString(header.Name))
		octetCountWithColor := octetCountColor.Sprintf("%d", header.OctetCount)
		fmt.Printf("%2s %4s %s\n", raff.IconToString(header.Icon), nameWithColor, octetCountWithColor)
		file.Seek(int64(header.OctetCount), os.SEEK_CUR)
	}

	return nil
}

func main() {
	ctx := kong.Parse(&Options{})

	err := ctx.Run()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
