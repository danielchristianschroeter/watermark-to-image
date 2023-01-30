# watermark-to-image
This application inserts a watermark in the lower right corner into all images of a specified images folder.

# Examples
Check out examples/original_images and examples/watermarked_images directory to get an idea what the application is doing.

| Original image | Watermarked image
|:--:|:--:
| ![landscape_original](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/original_images/landscape.jpg) | ![landscape_watermarked](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/watermarked_images/landscape.jpg) |
| ![portrait_original](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/original_images/portrait.jpg) | ![portrait_watermarked](https://raw.githubusercontent.com/danielchristianschroeter/watermark-to-image/main/examples/watermarked_images/portrait.jpg) |

## Usage
1. Build or download prebuild executable
2. Execute following command:
```
./watermark-to-image -sourceDirectory ./examples/original_images -targetDirectory ./examples/watermarked_images -watermarkImageFile ./examples/watermark.png
```

### Available command line parameter
```
./watermark-to-image --help
  -sourceDirectory string
        Source directory of original images
  -targetDirectory string
        Target directory for watermarked images
  -targetWaterkaredImageFilename string
        Rename all target files to the specified filename, if set targetWaterkaredImageExtension is required
  -targetWaterkaredImageFilenameSuffix string
        Dynamic Suffix for the filename definied in targetWaterkaredImageFilename added to every target file. Allowed values are 3DIGITSCOUNT (3 digits enumeration count) or RAND (random 6 digits number) (default "3DIGITSCOUNT")
  -targetWaterkaredImageHight int
        Resize target watermarked image hight, if targetWaterkaredImageWidth is empty aspect ratio will be presvered
  -targetWaterkaredImageWidth int
        Resize target watermarked image width, if targetWaterkaredImageHight is empty aspect ratio will be presvered
  -watermarkImageFile string
        Path and name of the watermark png image file
  -watermarkIncreasement float
        Scale factor of the watermark image file (default 100)
  -watermarkMarginBottom int
        Margin to the bottom edge for the watermark (default 20)
  -watermarkMarginRight int
        Margin to the right edge for the watermark (default 20)
  -watermarkOpacity float
        Opacity/Transparency of the watermark image file (default 0.5)
```

## Clone and build the project
```
$ git clone https://github.com/danielchristianschroeter/watermark-to-image
$ cd watermark-to-image
$ go build .
```
