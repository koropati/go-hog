package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
)

// HOGFeature represents a HOG feature.
type HOGFeature struct {
    Orientation int
    Magnitude float64
}

// ExtractHOGFeatures extracts HOG features from an image.
func ExtractHOGFeatures(image image.Image, cellSize int, blockStride int, numBins int) []HOGFeature {
    // Convert the image to grayscale.
    grayImage := image.NewGray(image.Rect(0, 0, image.Width(image), image.Height(image)))
    for y := 0; y < image.Height(); y++ {
        for x := 0; x < image.Width(); x++ {
            grayImage.SetGray(x, y, color.Gray{R: image.At(x, y).R, G: image.At(x, y).G, B: image.At(x, y).B})
        }
    }

    // Calculate the gradients of the image.
    gradients := make([]float64, image.Width(image)*image.Height(image))
    for y := 0; y < image.Height()-1; y++ {
        for x := 0; x < image.Width()-1; x++ {
            gradients[y*image.Width(image)+x] = math.Abs(grayImage.At(x+1, y).Y-grayImage.At(x, y).Y) + math.Abs(grayImage.At(x, y+1).X-grayImage.At(x, y).X)
        }
    }

    // Calculate the HOG features.
    features := make([]HOGFeature, 0)
    for y := 0; y < image.Height()/cellSize; y++ {
        for x := 0; x < image.Width()/cellSize; x++ {
            for o := 0; o < numBins; o++ {
                magnitude := 0.0
                for i := 0; i < cellSize; i++ {
                    for j := 0; j < cellSize; j++ {
                        if (y*cellSize+i)*image.Width(image)+x*cellSize+j < len(gradients) {
                            magnitude += gradients[y*cellSize+i]*cellSize*gradients[x*cellSize+j]*cellSize * math.Cos((float64(o)*2.0*math.Pi)/float64(numBins))
                        }
                    }
                }
                features = append(features, HOGFeature{Orientation: o, Magnitude: magnitude})
            }
        }
    }

    return features
}

func main() {
    // Load the image.
    image, err := image.Open("image.png")
    if err != nil {
        fmt.Println(err)
        return
    }

    // Extract the HOG features.
    features := ExtractHOGFeatures(image, 8, 4, 9)

    // Print the features.
    for _, feature := range features {
        fmt.Println(feature)
    }
}
