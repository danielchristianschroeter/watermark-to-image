# watermark-to-image

Watermark-to-Image is a command-line application for inserting watermarks into the lower right corner of images located in a specified folder. It supports various customization options.

## Table of Contents

- [Examples](#examples)
- [Usage](#usage)
- [Clone and Build](#clone-and-build)
- [Command Line Parameters](#command-line-parameters)
- [Download Executables](#download-executables)
- [License](#license)

# Examples

You can see the application in action by examining the following images:

|                                                                  Original image                                                                  |                                                                   Watermarked image                                                                    |
| :----------------------------------------------------------------------------------------------------------------------------------------------: | :----------------------------------------------------------------------------------------------------------------------------------------------------: |
| ![landscape_original](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/original_images/landscape.jpg) | ![landscape_watermarked](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/watermarked_images/landscape.jpg) |
|  ![portrait_original](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/original_images/portrait.jpg)  |  ![portrait_watermarked](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/watermarked_images/portrait.jpg)  |

## Usage

1. Build or download the prebuilt executable.
2. Execute the following command in your terminal:

```shell
./watermark-to-image -sourceDirectory ./examples/original_images -targetDirectory ./examples/watermarked_images -watermarkImageFile ./examples/watermark.png
```

## Clone and Build

You can also clone the project and build it locally. Here are the steps:

```shell
$ git clone https://github.com/danielchristianschroeter/watermark-to-image
$ cd watermark-to-image
$ go build .
```

## Command Line Parameters

The application supports various command-line parameters for customization. Here is a list of available parameters:

```shell
./watermark-to-image --help
  -sourceDirectory string
        Set the source directory for original images (required).
  -targetDirectory string
        Set the target directory for watermarked images (required).
  -targetWatermarkedImageFilename string
        Rename all target files to the specified filename. Requires 'targetWatermarkedImageExtension' to be set.
  -targetWatermarkedImageFilenameSuffix string
        Set the dynamic suffix for the filename defined in 'targetWatermarkedImageFilename'. Allowed values are '3DIGITSCOUNT' (3-digit enumeration count) or 'RAND' (random 6-digit number). Default is '3DIGITSCOUNT'. (default "3DIGITSCOUNT")
  -targetWatermarkedImageHeight int
        Resize the target watermarked image to the specified height (in pixels). Aspect ratio will be preserved if 'targetWatermarkedImageWidth' is empty.
  -targetWatermarkedImageMaxDimension int
        Specify the maximum dimension size for the target watermarked image. Use 0 to maintain the aspect ratio. Default is 0.
  -targetWatermarkedImageWidth int
        Resize the target watermarked image to the specified width (in pixels). Aspect ratio will be preserved if 'targetWatermarkedImageHeight' is empty.
  -watermarkImageFile string
        Specify the path and name of the watermark PNG image file (required).
  -watermarkIncreasement float
        Set the scale factor for the watermark image (in percentage). (default 100)
  -watermarkMarginBottom int
        Set the margin from the bottom edge for the watermark (in pixels). (default 20)
  -watermarkMarginRight int
        Set the margin from the right edge for the watermark (in pixels). (default 20)
  -watermarkOpacity float
        Set the opacity/transparency of the watermark image (0.0 to 1.0). (default 0.5)
```

## Download Executables

You can download prebuilt executables for various operating systems from the [Releases](https://github.com/danielchristianschroeter/watermark-to-image/releases) page.
_Note for macOS users: Due to the application not being officially signed, it may need to be manually approved under your system's security settings._

## License

This project is licensed under the [Apache License 2.0](LICENSE). The Apache License 2.0 is a permissive open-source license that provides comprehensive legal protections for contributors and users of the software. You are free to use, modify, and distribute this software according to the terms of the Apache License 2.0.
