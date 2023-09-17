package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/gabriel-vasile/mimetype"
)

// Used for build version information
var version = "development"

// Global variables used for command line flags
var (
	sourceDirectory                      string
	targetDirectory                      string
	watermarkImageFile                   string
	watermarkScaleFactor                 float64
	watermarkOpacity                     float64
	watermarkMarginRight                 int
	watermarkMarginBottom                int
	targetWatermarkedImageMaxDimension   int
	targetWatermarkedImageWidth          int
	targetWatermarkedImageHeight         int
	targetWatermarkedImageFilename       string
	targetWatermarkedImageFilenameSuffix string
)

func init() {
	// Initialize command line flags
	flag.StringVar(&sourceDirectory, "sourceDirectory", "", "Set the source directory for original images (required).")
	flag.StringVar(&targetDirectory, "targetDirectory", "", "Set the target directory for watermarked images (required).")
	flag.StringVar(&watermarkImageFile, "watermarkImageFile", "", "Specify the path and name of the watermark PNG image file (required).")
	flag.Float64Var(&watermarkScaleFactor, "watermarkScaleFactor", 100.0, "Set the scale factor for the watermark image (in percentage).")
	flag.Float64Var(&watermarkOpacity, "watermarkOpacity", 0.5, "Set the opacity/transparency of the watermark image (0.0 to 1.0).")
	flag.IntVar(&watermarkMarginRight, "watermarkMarginRight", 20, "Set the margin from the right edge for the watermark (in pixels).")
	flag.IntVar(&watermarkMarginBottom, "watermarkMarginBottom", 20, "Set the margin from the bottom edge for the watermark (in pixels).")
	flag.IntVar(&targetWatermarkedImageMaxDimension, "targetWatermarkedImageMaxDimension", 0, "Specify the maximum dimension size for the target watermarked image. Use 0 to maintain the aspect ratio. Default is 0.")
	flag.IntVar(&targetWatermarkedImageWidth, "targetWatermarkedImageWidth", 0, "Resize the target watermarked image to the specified width (in pixels). Aspect ratio will be preserved if 'targetWatermarkedImageHeight' is empty.")
	flag.IntVar(&targetWatermarkedImageHeight, "targetWatermarkedImageHeight", 0, "Resize the target watermarked image to the specified height (in pixels). Aspect ratio will be preserved if 'targetWatermarkedImageWidth' is empty.")
	flag.StringVar(&targetWatermarkedImageFilename, "targetWatermarkedImageFilename", "", "Rename all target files to the specified filename.")
	flag.StringVar(&targetWatermarkedImageFilenameSuffix, "targetWatermarkedImageFilenameSuffix", "3DIGITSCOUNT", "Set the dynamic suffix for the filename defined in 'targetWatermarkedImageFilename'. Allowed values are '3DIGITSCOUNT' (3-digit enumeration count) or 'RAND' (random 6-digit number). Default is '3DIGITSCOUNT'.")
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

	if (targetWatermarkedImageMaxDimension != 0) && (targetWatermarkedImageWidth != 0 || targetWatermarkedImageHeight != 0) {
		log.Fatal("Error: Only one of targetWatermarkedImageMaxDimension or targetWatermarkedImageWidth/targetWatermarkedImageHeight should be set.")
	}

	// Add slash to source directory, if not exists
	if !strings.HasSuffix(sourceDirectory, "/") {
		sourceDirectory += "/"
	}

	// Add slash to target directory, if not exists
	if !strings.HasSuffix(targetDirectory, "/") {
		targetDirectory += "/"
	}

	log.Println("Processing images from directory " + sourceDirectory + "...")
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

	// Sort the files by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	// Set a count
	count := 1

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
		if err := processImage(count, f, targetDirectory, watermarkImageFile, watermarkScaleFactor, watermarkMarginRight, watermarkMarginBottom, watermarkOpacity, targetWatermarkedImageMaxDimension, targetWatermarkedImageWidth, targetWatermarkedImageHeight, targetWatermarkedImageFilename, targetWatermarkedImageFilenameSuffix); err != nil {
			log.Printf("Failed to processImage of file %+v Error: %v", sourceDirectory+file.Name(), err)
			continue
		}
		count++
	}
}

func processImage(count int, file *os.File, targetDirectory string, watermarkImageFile string, watermarkScaleFactor float64, watermarkMarginRight int, watermarkMarginBottom int, watermarkOpacity float64, targetWatermarkedImageMaxDimension int, targetWatermarkedImageWidth int, targetWatermarkedImageHeight int, targetWatermarkedImageFilename string, targetWatermarkedImageFilenameSuffix string) error {
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
			log.Printf("Notice: Orientation tag not found in file %+v", file.Name())
			// Handle this case as needed, or continue with your code
		} else {
			orientation, err := results[0].Value()
			if err != nil {
				log.Fatalf("Failed to extract value of results from exif data in file %+v Error: %+v", file.Name(), err)
			}
			orientationSlice := orientation.([]uint16)
			orientationInt = orientationSlice[0]
		}
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

	// 	Get the width and the height of our watermark
	watermarkWidth := watermark.Bounds().Dx()
	watermarkHeight := watermark.Bounds().Dy()

	// Increase the dimensions of the watermark by a given percentage
	resizedWidth := int(float64(watermarkWidth) * (float64(watermarkScaleFactor) / 100.0))
	resizedHeight := int(float64(watermarkHeight) * (float64(watermarkScaleFactor) / 100.0))

	// Resize the watermark using Lanczos interpolation
	resizedWatermark := imaging.Resize(watermark, resizedWidth, resizedHeight, imaging.Lanczos)

	// Resize the watermarked image if specified
	if targetWatermarkedImageMaxDimension > 0 {
		img = resizeImageProportionally(img, targetWatermarkedImageMaxDimension)
	} else if targetWatermarkedImageWidth > 0 || targetWatermarkedImageHeight > 0 {
		img = resizeImage(img, targetWatermarkedImageWidth, targetWatermarkedImageHeight)
	}

	// Figure out the dstX value
	dstX := img.Bounds().Dx() - resizedWidth - watermarkMarginRight

	// Figure out the dstY value
	dstY := img.Bounds().Dy() - resizedHeight - watermarkMarginBottom

	// Overlay the watermark onto the photo image with a specified opacity/transparency
	dst := imaging.Overlay(img, resizedWatermark, image.Pt(dstX, dstY), watermarkOpacity)

	if targetWatermarkedImageFilename != "" && targetWatermarkedImageFilenameSuffix != "" {
		// Save the new image with a custom filename and a dynamic suffix
		var Suffix string
		switch targetWatermarkedImageFilenameSuffix {
		case "3DIGITSCOUNT":
			Suffix = fmt.Sprintf("%03d", count)
		case "RAND":
			Suffix = strconv.Itoa(rand.Intn(999999))
		}
		outFile, err := os.Create(targetDirectory + targetWatermarkedImageFilename + Suffix + filepath.Ext(file.Name()))
		if err != nil {
			log.Fatalf("Failed to os.Create of file %+v Error: %+v", targetDirectory+targetWatermarkedImageFilename+Suffix+filepath.Ext(file.Name()), err)
		}
		defer outFile.Close()
		jpeg.Encode(outFile, dst, nil)

		log.Println("Saved watermarked image to " + targetDirectory + targetWatermarkedImageFilename + Suffix + filepath.Ext(file.Name()))
	} else {
		// Save the new image with the same filename
		outFile, err := os.Create(targetDirectory + filepath.Base(file.Name()))
		if err != nil {
			log.Fatalf("Failed to os.Create of file %+v Error: %+v", targetDirectory+filepath.Base(file.Name()), err)
		}
		defer outFile.Close()
		jpeg.Encode(outFile, dst, nil)

		log.Println("Saved watermarked image to " + targetDirectory + filepath.Base(file.Name()))
	}

	return nil
}

// Resize an image proportionally to the specified maximum dimension.
func resizeImageProportionally(img image.Image, maxDimension int) *image.NRGBA {
	// Determine image dimensions and aspect ratio
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// Calculate the new dimensions while maintaining the aspect ratio
	var newWidth, newHeight int
	if height > width {
		newHeight = maxDimension
		newWidth = (maxDimension * width) / height
	} else {
		newWidth = maxDimension
		newHeight = (maxDimension * height) / width
	}

	// Resize the image
	resizedImg := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	return resizedImg
}

// Resize an image to the specified width and height.
func resizeImage(img image.Image, width, height int) *image.NRGBA {
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	return resizedImg
}
