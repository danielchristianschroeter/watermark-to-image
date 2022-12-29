# watermark-to-image
This application inserts a watermark in the lower right corner into all images of an specified images folder.

# Examples
Check out examples/original_images and examples/watermarked_images to get an idea what the application is doing.

## Usage
1. Build or download prebuild executeable
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
