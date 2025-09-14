package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/groot314/timelapse/internal/helpers"
	"github.com/groot314/timelapse/internal/progress"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	input := flag.String("i", "", "Input location for photos.")
	startNumber := flag.String("s", "", "Start number for photos.")
	outputName := flag.String("o", "", "Output name.")
	flag.Parse()
	if *input == "" || *outputName == "" {
		fmt.Println("Need to input required flags. -i <input_dir> -s <start_number_index> -o <output_file_name>")
		flag.Usage()
		os.Exit(1)
	}
	if *startNumber == "" {
		*startNumber, _ = helpers.FindPhotoStart(*input)
	}

	ffmpegInput := *input + "/DSC%05d.JPG"
	ffmpegInputArgs := ffmpeg.KwArgs{"framerate": "24", "start_number": *startNumber}
	a, err := ffmpeg.Probe(ffmpegInput, ffmpegInputArgs)
	if err != nil {
		panic(err)
	}
	totalDuration, err := progress.ProbeDuration(a)
	if err != nil {
		panic(err)
	}

	err = ffmpeg.Input(ffmpegInput, ffmpegInputArgs).
		Output("./"+*outputName+".mp4", ffmpeg.KwArgs{"c:v": "libx264", "pix_fmt": "yuv420p"}).
		GlobalArgs("-progress", "unix://"+progress.ProgressTempSock(totalDuration)).
		OverWriteOutput().Run()
	// OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		fmt.Println("Error with ffmpeg:", err)
	}
}
