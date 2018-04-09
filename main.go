/*
MIT License

Copyright (c) [year] [fullname]

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
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	arguments := os.Args[1:]
	path := arguments[0]
	name := strings.TrimSuffix(path, filepath.Ext(path))

	fmt.Fprintf(os.Stdout, "Opening file: %s\n", path)
	fmt.Fprintf(os.Stdout, "Name: %s\n", name)

	fImg, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open Error: %s", err.Error())
		return
	}
	defer fImg.Close()

	img1, _, err := image.Decode(fImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Decode Error: %s", err.Error())
		return
	}

	if strings.Contains(name, "@2x") {
		name = strings.TrimSuffix(name, "@2x")
		err = scaleFromRetina(img1, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			return
		}
	} else {
		name = strings.TrimSuffix(name, "@3x")
		err = scaleFromSuperRetina(img1, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			return
		}
	}
}

func scaleFromRetina(img1 image.Image, name string) error {
	originalWidth := float64(img1.Bounds().Max.X)
	originalHeight := float64(img1.Bounds().Max.Y)

	err := scaleImage(img1, uint(originalWidth*1.5), uint(originalHeight*1.5), name, "@3x")
	if err != nil {
		return errors.New("Resize Error: " + err.Error())
	}

	err = scaleImage(img1, uint(originalWidth), uint(originalHeight), name, "@2x")
	if err != nil {
		return errors.New("Resize Error: %s" + err.Error())
	}

	err = scaleImage(img1, uint(originalWidth*0.5), uint(originalHeight*0.5), name, "")
	if err != nil {
		return errors.New("Resize Error: %s" + err.Error())
	}

	return nil
}

func scaleFromSuperRetina(img1 image.Image, name string) error {
	originalWidth := float64(img1.Bounds().Max.X)
	originalHeight := float64(img1.Bounds().Max.Y)

	err := scaleImage(img1, uint(originalWidth), uint(originalHeight), name, "@3x")
	if err != nil {
		return errors.New("Resize Error: " + err.Error())
	}

	err = scaleImage(img1, uint(originalWidth*0.66), uint(originalHeight*0.66), name, "@2x")
	if err != nil {
		return errors.New("Resize Error: %s" + err.Error())
	}

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
