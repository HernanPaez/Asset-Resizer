# asset-resizer
## About
An iOS Asset resizer made with Go. This application generates all assets sizes from one single image.

* If you provide an @2x image it will generate @3x and @1x images
* If you provide an @3x image it will generate @2x and @1x images


## Installing 
Unless you're a developer, you probably want a binary you can run. [Download](https://github.com/HernanPaez/asset-resizer/releases) the latest release for macOS.
It should work on Linux and Windows platforms by building it from source. I will try to generate build for these platforms soon.

1. Change file permissions and make it executable:
```bash
chmod +x asset-resizer
```

2. Make a symlink to your bin folder:
```bash
ls /usr/local/bin/ > /dev/null || sudo mkdir /usr/local/bin/ ; sudo ln -s $(pwd)/asset-resizer /usr/local/bin/asset-resizer
```

## Usage

If the filename contains @2x as a suffix it will take this as an normal retina image
else the application will process the image as a Super-Retina (@3x)

this will generate @3x, @2x, @1x by scaling image.png to @2x and @1x

```bash
asset-generator path/to/my/image.png 
```

this will generate @3x, @2x, @1x by scaling image.png to @3x and @1x

```bash
asset-generator path/to/my/retina-image@2x.png
```

this will generate @3x, @2x, @1x by scaling image.png to @3x and @1x

```bash
./asset-generator path/to/my/super-retina-image@3x.png
```
