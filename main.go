package main

import (
	"flag"
	"fmt"
	"os"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	input := flag.String("i", "", "Input location for photos.")
	startNumber := flag.String("s", "", "Start number for photos.")
	outputName := flag.String("o", "", "Output name.")
	flag.Parse()
	if *input == "" || *startNumber == "" || *outputName == "" {
		fmt.Println("Need to input required flags. -i <input_dir> -s <start_number_index> -o <output_file_name>")
		flag.Usage()
		os.Exit(1)
	}
	err := ffmpeg.Input(*input+"/DSC%05d.JPG", ffmpeg.KwArgs{"framerate": "24", "start_number": *startNumber}).
		Output("./"+*outputName+".mp4", ffmpeg.KwArgs{"c:v": "libx264", "pix_fmt": "yuv420p"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		fmt.Println("Error with ffmpeg:", err)
	}
}
