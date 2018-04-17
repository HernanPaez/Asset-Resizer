/*
MIT License

Copyright (c) 2018 Luis Hernan Paez

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "You need to provide a file!\n\nSyntax: asset-resizer path/to/file.png\n")
		os.Exit(1)
		return
	}
	if len(os.Args) > 2 {
		fmt.Fprintf(os.Stderr, "You can't resize more than one file at the same time!\n\nSyntax: asset-resizer path/to/file.png\n")
		os.Exit(1)
		return
	}

	arguments := os.Args[1:]
	path := arguments[0]

	stat, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Stat Error: %s", err.Error())
		os.Exit(1)
		return
	}

	if stat.IsDir() {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Directory Reading Error: %s", err.Error())
			os.Exit(1)
			return
		}

		for _, f := range files {
			if f.IsDir() == false && strings.ToLower(filepath.Ext(f.Name())) == ".png" {

				err := resizeFileAtPath(path + "/" + f.Name())
				if err != nil {
					os.Exit(1)
					return
				}

			}
		}

	} else {
		resizeFileAtPath(path)
	}

	fmt.Fprintf(os.Stdout, "\nDone!\n")
	os.Exit(0)
}

//Scales an image from @2x -> @3x @1x
func scaleFromRetina(img1 image.Image, name string) error {
	originalWidth := float64(img1.Bounds().Max.X)
	originalHeight := float64(img1.Bounds().Max.Y)

	fmt.Fprintf(os.Stdout, "\tCreating @3x file\n")
	err := scaleImage(img1, uint(originalWidth*1.5), uint(originalHeight*1.5), name, "@3x")
	if err != nil {
		return errors.New("Resize Error: " + err.Error())
	}

	fmt.Fprintf(os.Stdout, "\tRenaming @2x file\n")
	err = scaleImage(img1, uint(originalWidth), uint(originalHeight), name, "@2x")
	if err != nil {
		return errors.New("Resize Error: %s" + err.Error())
	}

	fmt.Fprintf(os.Stdout, "\tCreating @1x file\n")
	err = scaleImage(img1, uint(originalWidth*0.5), uint(originalHeight*0.5), name, "")
	if err != nil {
		return errors.New("Resize Error: %s" + err.Error())
	}

	return nil
}

func resizeFileAtPath(path string) error {
	name := strings.TrimSuffix(path, filepath.Ext(path))
	fmt.Fprintf(os.Stdout, "\n- Resizing file: %s", path)

	fImg, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open Error: %s", err.Error())
		return err
	}
	defer fImg.Close()

	img1, _, err := image.Decode(fImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Decode Error: %s", err.Error())
		return err
	}

	if strings.Contains(name, "@2x") {
		fmt.Fprintf(os.Stdout, "\n- Resizing from @2x...\n")
		name = strings.TrimSuffix(name, "@2x")
		err = scaleFromRetina(img1, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			return err
		}
	} else {
		fmt.Fprintf(os.Stdout, "\n- Resizing from @3x...\n")
		name = strings.TrimSuffix(name, "@3x")
		err = scaleFromSuperRetina(img1, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			return err
		}
	}

	return nil

}

//Scales an image from @3x -> @2x @1x
func scaleFromSuperRetina(img1 image.Image, name string) error {
	originalWidth := float64(img1.Bounds().Max.X)
	originalHeight := float64(img1.Bounds().Max.Y)

	fmt.Fprintf(os.Stdout, "\tRenaming @3x file\n")
	err := scaleImage(img1, uint(originalWidth), uint(originalHeight), name, "@3x")
	if err != nil {
		return errors.New("Resize Error: " + err.Error())
	}

	fmt.Fprintf(os.Stdout, "\tCreating @2x file\n")
	err = scaleImage(img1, uint(originalWidth*0.66), uint(originalHeight*0.66), name, "@2x")
	if err != nil {
		return errors.New("Resize Error: %s" + err.Error())
	}

	fmt.Fprintf(os.Stdout, "\tCreating @1x file\n")
	err = scaleImage(img1, uint(originalWidth*0.33), uint(originalHeight*0.33), name, "")
	if err != nil {
		return errors.New("Resize Error: %s" + err.Error())
	}

	return nil
}

func scaleImage(img image.Image, width uint, height uint, filename string, suffix string) error {

	newImage := resize.Resize(width, height, img, resize.Lanczos3)

	file, err := os.Create(filename + suffix + ".png")
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, newImage)
	return err
}
