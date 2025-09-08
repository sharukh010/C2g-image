package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// ImageData holds the original URL and the image data
type ImageData struct {
	URL string 
	Data image.Image 
}


// Global shared map to store processed image data
var (
	processedImages = make(map[string]image.Image)
	// Mutex to protect the shared map from concurrent writes
	mu sync.Mutex
	// WaitGroup to wait for all goroutines to complet 
	wg sync.WaitGroup 
)

func main(){
	//channels for our pipeline 
	urlsCh := make(chan string)

	//Buffered channel to hold downloaded images, allowing the downloader to work ahead 
	// The buffer size can be a tuning parameter 
	downloadedCh := make(chan *ImageData,10)
	// Unbuffered channel for processed images 
	processedCh := make(chan *ImageData)

	//Add goroutines to the WaitGroup 
	wg.Add(3)

	// Start the pipeline stages
	go downloader(urlsCh,downloadedCh)
	go processor(downloadedCh,processedCh)
	go collector(processedCh)

	// A list of example image URLs 
	images := []string{
		"https://images.unsplash.com/photo-1506748686214-e9df14d4d9d0",
		"https://images.unsplash.com/photo-1517816743773-6e0fd518b4a6",
		"https://images.unsplash.com/photo-1499084732479-de2c02d45fc4",
		"https://images.unsplash.com/photo-1522202176988-66273c2fd55f",
		"https://images.unsplash.com/photo-1507525428034-b723cf961d3e",
		"https://images.unsplash.com/photo-1519125323398-675f0ddb6308",
		"https://images.unsplash.com/photo-1470770841072-f978cf4d019e",
		"https://images.unsplash.com/photo-1494790108377-be9c29b29330",
		"https://images.unsplash.com/photo-1534528741775-53994a69daeb",
		"https://images.unsplash.com/photo-1441974231531-c6227db76b6e",
	}

	// Send URLs to the downloader. This is a crucial loop. 
	for _,url := range images {
		urlsCh <- url 
	}
	close(urlsCh)

	wg.Wait()

	fmt.Printf("\nProcessing complete! Processed %d images.\n",len(processedImages))


}

// Dowloader goroutine 
func downloader(urls <-chan string, out chan<- *ImageData){
	defer wg.Done()
	defer close(out)
	if err := os.MkdirAll("input",0755); err != nil {
		log.Fatal(err)
	}
	for url:=range urls {
		log.Printf("Downloading %s...",url)
		res,err := http.Get(url)
		if err != nil {
			log.Printf("Error downloading %s: %v",url,err)
			continue 
		}
		defer res.Body.Close() 

		
		img,_,err := image.Decode(res.Body)
		if err != nil {
			log.Printf("Error decoding image from %s: %v",url,err)
			continue 
		}
		
		//send the data to the next stage 
		out <- &ImageData{URL:url,Data: img}

		fileName := filepath.Base(url)
		inputPath := filepath.Join("input",fmt.Sprintf("%s.png",fileName))
		log.Printf("Saving processed image to %s...",inputPath)

		file,err := os.Create(inputPath)
		if err != nil {
			log.Printf("Error creating file %s: %v",inputPath,err)
			continue
		}

		defer file.Close()

		if err := png.Encode(file,img); err != nil {
			log.Printf("Error encoding image to PNG: %v",err)
		}
	}
}

// Processor goroutine 
func processor(in <-chan *ImageData, out chan<- *ImageData){
	defer wg.Done()
	defer close(out)

	for imageData := range in {
		log.Printf("Processing image from %s...",imageData.URL)

		// Simple grayscale conversion 
		bounds := imageData.Data.Bounds()
		grayImg := image.NewGray(bounds)
		for y:= bounds.Min.Y; y<bounds.Max.Y; y++ {
			for x:= bounds.Min.X; x<bounds.Max.X; x++ {
				originalColor := imageData.Data.At(x,y)
				grayColor := color.GrayModel.Convert(originalColor)
				grayImg.Set(x,y,grayColor)
			}
		}
		imageData.Data = grayImg 
		out <- imageData
	}
}

// Collector goroutine 
func collector(in <-chan *ImageData){
	defer wg.Done() 
	if err := os.MkdirAll("output",0755); err != nil {
		log.Fatal(err)
	}
	for imageData := range in {
		log.Printf("Collecting processed image from %s...",imageData.URL)
		// Protect the shared map with a Mutex 
		mu.Lock()
		processedImages[imageData.URL] = imageData.Data 
		mu.Unlock()

		fileName := filepath.Base(imageData.URL)
		outputPath := filepath.Join("output",fmt.Sprintf("%s.png",fileName))
		log.Printf("Saving processed image to %s...",outputPath)

		file,err := os.Create(outputPath)
		if err != nil {
			log.Printf("Error creating file %s: %v",outputPath,err)
			continue
		}

		defer file.Close()

		if err := png.Encode(file,imageData.Data); err != nil {
			log.Printf("Error encoding image to PNG: %v",err)
		}
	}
}