package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/recolude/imgconv/imgconv"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func shrinkCommand(curDir string, filenameMatcher string, recurse bool, maxSize int, out io.Writer) error {
	files, err := ioutil.ReadDir(curDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && recurse {
			err := shrinkCommand(filepath.Join(curDir, file.Name()), filenameMatcher, recurse, maxSize, out)
			if err != nil {
				return err
			}
		}
		match, err := filepath.Match(filenameMatcher, file.Name())
		if err != nil {
			return err
		}
		if match {
			pngPath := filepath.Join(curDir, file.Name())
			fmt.Fprintf(out, "%s\n", pngPath)

			pngFile, err := os.OpenFile(pngPath, os.O_RDWR, 0o666)
			if err != nil {
				return err
			}

			pngInData, err := ioutil.ReadAll(pngFile)
			if err != nil {
				return err
			}

			pngOutData := &bytes.Buffer{}

			err = imgconv.ResizePNG(bytes.NewBuffer(pngInData), pngOutData, maxSize)
			if err != nil {
				return err
			}

			_, err = pngFile.Seek(0, 0)
			if err != nil {
				return err
			}

			finalData := pngOutData.Bytes()
			written, err := pngFile.Write(finalData)
			if err != nil {
				return err
			}
			pngFile.Truncate(int64(written))

			pngFile.Close()

		}
	}
	return nil
}

func convertCommand(curDir string, filenameMatcher string, recurse bool, delete bool, maxSize int, out io.Writer) error {
	files, err := ioutil.ReadDir(curDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() && recurse {
			err := convertCommand(filepath.Join(curDir, file.Name()), filenameMatcher, recurse, delete, maxSize, out)
			if err != nil {
				return err
			}
		}
		match, err := filepath.Match(filenameMatcher, file.Name())
		if err != nil {
			return err
		}
		if match {
			imgPath := filepath.Join(curDir, file.Name())
			fmt.Fprintf(out, "%s\n", imgPath)
			imgIn, err := os.Open(imgPath)
			if err != nil {
				return err
			}

			pngPath := strings.TrimSuffix(imgPath, filepath.Ext(imgPath)) + ".png"
			pngFile, err := os.Create(pngPath)
			if err != nil {
				return err
			}

			switch filepath.Ext(strings.ToLower(imgPath)) {
			case ".tga":
				err = imgconv.ConvertTGA(imgIn, pngFile, maxSize)
				if err != nil {
					return err
				}

			case ".tiff", ".tif":
				err = imgconv.ConvertTIFF(imgIn, pngFile, maxSize)
				if err != nil {
					return err
				}
			}

			pngFile.Close()
			imgIn.Close()

			if delete {
				err = os.Remove(imgPath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func buildCLI(out io.Writer, errOut io.Writer, onErrFunc func(error), fs afero.Fs) *cli.App {
	return &cli.App{
		Name:    "texmap",
		Usage:   "ingests multiple maps from different textures and builds you new images from the desired maps",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name: "Eli C. Davis",
			},
		},
		Writer:    out,
		ErrWriter: errOut,
		ExitErrHandler: cli.ExitErrHandlerFunc(func(ctx *cli.Context, err error) {
			onErrFunc(err)
		}),
		Commands: []*cli.Command{
			{
				Name:      "to-png",
				Aliases:   []string{"tp"},
				Usage:     "create pngs from tgas or tiffs",
				ArgsUsage: "[filename matcher]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "recurse",
						Value:    false,
						Usage:    "search all sub directories",
						Required: false,
						Aliases:  []string{"r"},
					},
					&cli.BoolFlag{
						Name:     "delete",
						Value:    false,
						Usage:    "delete the tgas that get converted",
						Required: false,
						Aliases:  []string{"d"},
					},
					&cli.IntFlag{
						Name:     "max-size",
						Value:    0,
						Usage:    "max size of the resulting image",
						Required: false,
						Aliases:  []string{"s"},
					},
				},
				Action: func(c *cli.Context) error {
					resursive := c.Bool("recurse")
					delete := c.Bool("delete")
					size := c.Int("max-size")

					if c.Args().Len() == 0 {
						return convertCommand("./", "*.tga", resursive, delete, size, c.App.Writer)
					}

					return convertCommand("./", c.Args().First(), resursive, delete, size, c.App.Writer)
				},
			},
			{
				Name:      "resize",
				Usage:     "Shinks images if they exceed a given threshold specified",
				ArgsUsage: "[filename matcher]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:     "recurse",
						Value:    false,
						Usage:    "search all sub directories",
						Required: false,
						Aliases:  []string{"r"},
					},
					&cli.IntFlag{
						Name:     "max-size",
						Value:    2048,
						Usage:    "max size of the resulting image",
						Required: true,
						Aliases:  []string{"m"},
					},
				},
				Action: func(c *cli.Context) error {
					resursive := c.Bool("recurse")
					size := c.Int("max-size")

					if c.Args().Len() == 0 {
						return shrinkCommand("./", "*.png", resursive, size, c.App.Writer)
					}

					err := shrinkCommand("./", c.Args().First(), resursive, size, c.App.Writer)

					return err
				},
			},
		},
	}
}

func main() {
	if err := buildCLI(os.Stdout, os.Stderr, cli.HandleExitCoder, afero.NewOsFs()).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
