package main

import (
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/image/draw"
)

// Used for build version information
var version = "development"

// Global variables used for command line flags
var (
	sourceDirectory       string
	targetDirectory       string
	watermarkImageFile    string
	watermarkIncreasement float64
	watermarkOpacity      float64
	watermarkMarginRight  int
	watermarkMarginBottom int
)

func init() {
	// Initialize command line flags
	flag.StringVar(&sourceDirectory, "sourceDirectory", "", "Source directory of original images")
	flag.StringVar(&targetDirectory, "targetDirectory", "", "Target directory for watermarked images")
	flag.StringVar(&watermarkImageFile, "watermarkImageFile", "", "Path and name of the watermark png image file")
	flag.Float64Var(&watermarkIncreasement, "watermarkIncreasement", 100, "Scale factor of the watermark image file")
	flag.Float64Var(&watermarkOpacity, "watermarkOpacity", 0.5, "Opacity/Transparency of the watermark image file")
	flag.IntVar(&watermarkMarginRight, "watermarkMarginRight", 20, "Margin to the right edge for the watermark")
	flag.IntVar(&watermarkMarginBottom, "watermarkMarginBottom", 20, "Margin to the bottom edge for the watermark")
}

func processImage(file *os.File, targetDirectory string, watermarkImageFile string, watermarkIncreasement float64, watermarkMarginRight int, watermarkMarginBottom int, watermarkOpacity float64) error {
	var orientationInt uint16
	// Load the watermark png image file using image.Decode
	watermarkFile, err := os.Open(watermarkImageFile)
	if err != nil {
		log.Fatalf("Failed to os.Open of file %+v Error: %+v", watermarkImageFile, err)
	}
	defer watermarkFile.Close()
	watermark, _, err := image.Decode(watermarkFile)
	if err != nil {
		log.Fatalf("Failed to image.Decode of file %+v Error: %+v", watermarkImageFile, err)
	}

	// Read exif metadata of image
	exifData, _ := exif.SearchAndExtractExifWithReader(file)

	// Check if exif metadata available in image
	if len(exifData) > 0 {
		im, err := exifcommon.NewIfdMappingWithStandard()
		if err != nil {
			log.Fatalf("Failed to exifcommon.NewIfdMappingWithStandard in file %+v Error: %+v", file.Name(), err)
		}
		ti := exif.NewTagIndex()

		// Read Orientation-Tag from exif metadata
		_, index, err := exif.Collect(im, ti, exifData)
		if err != nil {
			log.Fatalf("Failed to exif.Collect in file %+v Error: %+v", file.Name(), err)
		}
		results, err := index.RootIfd.FindTagWithName("Orientation")
		if err != nil {
			log.Fatalf("Failed to index.RootIfd.FindTagWithName in file %+v Error: %+v", file.Name(), err)
		}
		orientation, err := results[0].Value()
		if err != nil {
			log.Fatalf("Failed to extract value of results from exif data in file %+v Error: %+v", file.Name(), err)
		}
		orientationSlice := orientation.([]uint16)
		orientationInt = orientationSlice[0]
	} else {
		log.Printf("Notice: No exif data found for file %+v", file.Name())
	}

	// Seek to the beginning of the file before calling image.Decode
	_, err = file.Seek(0, 0)
	if err != nil {
		log.Fatalf("Failed to file.Seek of file %+v Error: %+v", file.Name(), err)
	}

	// Decode the image from the buffer
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to image.Decode of file %+v Error: %+v", file.Name(), err)
	}

	// Rotate and/or flip image based on Orientation-Tag from EXIF-Metadata
	switch orientationInt {
	case 0:
		// no adjustment needed
	case 1:
		// no adjustment needed
	case 2:
		img = imaging.FlipH(img)
	case 3:
		img = imaging.Rotate(img, 180, color.Transparent)
	case 4:
		img = imaging.FlipH(imaging.Rotate(img, 180, color.Transparent))
	case 5:
		img = imaging.FlipV(imaging.Rotate(img, 90, color.Transparent))
	case 6:
		img = imaging.Rotate(img, 270, color.Transparent)
	case 7:
		img = imaging.FlipV(imaging.Rotate(img, 270, color.Transparent))
	case 8:
		img = imaging.Rotate(img, 90, color.Transparent)
	default:
		// no adjustment needed
	}

	// Get the width and the height of our image.
	photoWidth := img.Bounds().Dx()
	photoHeight := img.Bounds().Dy()

	// 	Get the width and the height of our watermark.
	watermarkWidth := watermark.Bounds().Dx()
	watermarkHeight := watermark.Bounds().Dy()

	// Increase the dimensions of the watermark by a given percentage.
	resizedWidth := int(float64(watermarkWidth) * (float64(watermarkIncreasement) / 100.0))
	resizedHeight := int(float64(watermarkHeight) * (float64(watermarkIncreasement) / 100.0))
	resizedWatermark := image.NewRGBA(image.Rect(0, 0, resizedWidth, resizedHeight))
	draw.NearestNeighbor.Scale(resizedWatermark, resizedWatermark.Bounds(), watermark, watermark.Bounds(), draw.Over, nil)

	// Figure out the dstX value.
	dstX := photoWidth - resizedWidth - watermarkMarginRight

	// Figure out the dstY value.
	dstY := photoHeight - resizedHeight - watermarkMarginBottom

	// Overlay the watermark onto the photo image with a specified opacity/transparency
	dst := imaging.Overlay(img, resizedWatermark, image.Pt(dstX, dstY), watermarkOpacity)

	// Save the new image
	outFile, err := os.Create(targetDirectory + filepath.Base(file.Name()))
	if err != nil {
		log.Fatalf("Failed to os.Create of file %+v Error: %+v", targetDirectory+filepath.Base(file.Name()), err)
	}
	defer outFile.Close()
	jpeg.Encode(outFile, dst, nil)

	log.Println("Saved watermarked image to " + targetDirectory + filepath.Base(file.Name()))

	return nil
}

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Println("watermark-to-image. Version: " + version)
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(sourceDirectory) == 0 || len(targetDirectory) == 0 || len(watermarkImageFile) == 0 {
		log.Fatal("Usage: -sourceDirectory <sourceDirectory> -targetDirectory <targetDirectory> -watermarkImageFile <watermarkImageFile>")
	}

	// Add slash to source directory, if not exists
	if !strings.HasSuffix(sourceDirectory, "/") {
		sourceDirectory += "/"
	}

	// Add slash to target directory, if not exists
	if !strings.HasSuffix(targetDirectory, "/") {
		targetDirectory += "/"
	}

	log.Println("Processing imges from directory " + sourceDirectory + "...")
	// Open the source directory
	dir, err := os.Open(sourceDirectory)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	// Get a list of all the files in the source directory
	files, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the images in the source directory
	for _, file := range files {
		// Skip files that start with a dot (e.g. .DS_Store file)
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		// Open the image
		f, err := os.Open(sourceDirectory + file.Name())
		if err != nil {
			log.Printf("Failed to os.Open of file %+v Error: %v", sourceDirectory+file.Name(), err)
			continue
		}
		defer f.Close()

		// Skip image/heic
		mimeType, err := mimetype.DetectFile(sourceDirectory + file.Name())
		if err != nil {
			log.Printf("Failed to mimetype.DetectFile of file %+v Error: %v", sourceDirectory+file.Name(), err)
		}
		switch mimeType.String() {
		case "image/heic":
			log.Printf("Skipping %v Reason: heic images are not supported yet", sourceDirectory+file.Name())
			continue
		}

		// Process the image
		if err := processImage(f, targetDirectory, watermarkImageFile, watermarkIncreasement, watermarkMarginRight, watermarkMarginBottom, watermarkOpacity); err != nil {
			log.Printf("Failed to processImage of file %+v Error: %v", sourceDirectory+file.Name(), err)
			continue
		}
	}
}
