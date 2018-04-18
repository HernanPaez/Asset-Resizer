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
	"encoding/json"
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

type AssetImage struct {
	Idiom    string `json:"idiom"`
	Filename string `json:"filename"`
	Scale    string `json:"scale"`
}

type AssetInfo struct {
	Version int    `json:"version"`
	Author  string `json:"author"`
}

type AssetCatalog struct {
	Images []AssetImage `json:"images"`
	Info   AssetInfo    `json:"author"`
}

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
		scanDir(path, true, true)
	} else {
		resizeFileAtPath(path)
	}

	fmt.Fprintf(os.Stdout, "\nDone!\n")
	os.Exit(0)
}

func getCleanName(filename string) string {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	atIdx := strings.Index(name, "@")
	if atIdx != -1 {
		name = name[0:atIdx]
	}
	return name
}

func getFileNameFromPath(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx != -1 && len(path) > idx+1 {
		return path[idx+1 : len(path)]
	}
	return path
}

func getDirFromPath(path string) string {
	idx := strings.LastIndex(path, "/")
	if idx != -1 {
		return path[0:idx]
	}
	return path
}

func scanDir(path string, recursive bool, createAssetCatalog bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		fpath := path + "/" + f.Name()
		if f.IsDir() == false && strings.ToLower(filepath.Ext(f.Name())) == ".png" {
			//Get a clean filename, remove extension and any scale reference ie: @2x @3x

			//Generate scaled files
			assets, err := resizeFileAtPath(fpath)
			if err != nil {
				return err
			}

			//create json file
			if createAssetCatalog {
				info := AssetInfo{Version: 1, Author: "asset-resizer"}
				catalog := AssetCatalog{Images: assets, Info: info}
				saveJSON(catalog, getDirFromPath(fpath))
			}

		} else if f.IsDir() && recursive {
			if err := scanDir(fpath, recursive, createAssetCatalog); err != nil {
				return err
			}
		}
	}

	return nil
}

//Scales an image from @2x -> @3x @1x
func scaleFromRetina(img1 image.Image, path string) ([]AssetImage, error) {
	destinationPath := getCleanName(path)

	originalWidth := float64(img1.Bounds().Max.X)
	originalHeight := float64(img1.Bounds().Max.Y)

	assets := make([]AssetImage, 3)

	fmt.Fprintf(os.Stdout, "\tCreating @3x file\n")
	path, err := scaleImage(img1, uint(originalWidth*1.5), uint(originalHeight*1.5), destinationPath, "@3x")
	if err != nil {
		return nil, errors.New("Resize Error: " + err.Error())
	}
	assets[0] = AssetImage{Idiom: "universal", Filename: getFileNameFromPath(path), Scale: "3x"}

	fmt.Fprintf(os.Stdout, "\tRenaming @2x file\n")
	path, err = scaleImage(img1, uint(originalWidth), uint(originalHeight), destinationPath, "@2x")
	if err != nil {
		return nil, errors.New("Resize Error: %s" + err.Error())
	}
	assets[1] = AssetImage{Idiom: "universal", Filename: getFileNameFromPath(path), Scale: "2x"}

	fmt.Fprintf(os.Stdout, "\tCreating @1x file\n")
	path, err = scaleImage(img1, uint(originalWidth*0.5), uint(originalHeight*0.5), destinationPath, "")
	if err != nil {
		return nil, errors.New("Resize Error: %s" + err.Error())
	}
	assets[2] = AssetImage{Idiom: "universal", Filename: getFileNameFromPath(path), Scale: "1x"}

	return nil, nil
}

func resizeFileAtPath(path string) ([]AssetImage, error) {
	fmt.Fprintf(os.Stdout, "\n- Resizing file: %s", path)

	fImg, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Open Error: %s", err.Error())
		return nil, err
	}
	defer fImg.Close()

	img1, _, err := image.Decode(fImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Decode Error: %s", err.Error())
		return nil, err
	}

	if strings.Contains(strings.ToLower(path), "@2x") {
		fmt.Fprintf(os.Stdout, "\n- Resizing from @2x...\n")
		assets, err := scaleFromRetina(img1, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			return nil, err
		}

		return assets, nil

	} else {
		fmt.Fprintf(os.Stdout, "\n- Resizing from @3x...\n")
		assets, err := scaleFromSuperRetina(img1, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err.Error())
			return nil, err
		}

		return assets, nil
	}
}

//Scales an image from @3x -> @2x @1x
func scaleFromSuperRetina(img1 image.Image, path string) ([]AssetImage, error) {
	destinationPath := getCleanName(path)
	originalWidth := float64(img1.Bounds().Max.X)
	originalHeight := float64(img1.Bounds().Max.Y)

	assets := make([]AssetImage, 3)

	fmt.Fprintf(os.Stdout, "\tRenaming @3x file\n")
	path, err := scaleImage(img1, uint(originalWidth), uint(originalHeight), destinationPath, "@3x")
	if err != nil {
		return nil, errors.New("Resize Error: " + err.Error())
	}
	assets[0] = AssetImage{Idiom: "universal", Filename: getFileNameFromPath(path), Scale: "3x"}

	fmt.Fprintf(os.Stdout, "\tCreating @2x file\n")
	path, err = scaleImage(img1, uint(originalWidth*0.66), uint(originalHeight*0.66), destinationPath, "@2x")
	if err != nil {
		return nil, errors.New("Resize Error: %s" + err.Error())
	}
	assets[1] = AssetImage{Idiom: "universal", Filename: getFileNameFromPath(path), Scale: "2x"}

	fmt.Fprintf(os.Stdout, "\tCreating @1x file\n")
	path, err = scaleImage(img1, uint(originalWidth*0.33), uint(originalHeight*0.33), destinationPath, "")
	if err != nil {
		return nil, errors.New("Resize Error: %s" + err.Error())
	}
	assets[2] = AssetImage{Idiom: "universal", Filename: getFileNameFromPath(path), Scale: "1x"}

	return assets, nil
}

func saveJSON(catalog AssetCatalog, path string) error {
	raw, err := json.Marshal(catalog)
	if err != nil {
		return err
	}

	jsonPath := path + "/Contents.json"
	file, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(raw)

	return nil
}

func scaleImage(img image.Image, width uint, height uint, path string, suffix string) (string, error) {

	newImage := resize.Resize(width, height, img, resize.Lanczos3)

	newPath := path + suffix + ".png"
	file, err := os.Create(newPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = png.Encode(file, newImage)
	return newPath, err
}
