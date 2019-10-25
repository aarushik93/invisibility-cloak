package main

import (
	"gocv.io/x/gocv"
	"time"
)

func main() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Invisibility Cloak")
	curImg := gocv.NewMatWithSize(5,5, gocv.MatTypeCV8U)
	background := gocv.NewMatWithSize(5,5, gocv.MatTypeCV8U)
	hsvImg := gocv.NewMatWithSize(5,5, gocv.MatTypeCV8U)

	// this is the mask where we do all the morphing (the magic)
	morphMask := gocv.NewMat()

	// this is the mask where will keep all the inverted image data
	clearMask := gocv.NewMat()

	// this will be our final frame
	processedFrame := gocv.NewMatWithSize(5,5, gocv.MatTypeCV8U)


	// hsv values for green
	lGreenMat := gocv.NewScalar(40,0,0, 0)
	uGreenMat := gocv.NewScalar(95,255,255, 0)


	defer func() {
		curImg.Close()
		background.Close()
	}()

	time.Sleep(5 * time.Second)
	webcam.Read(&background)
	gocv.Flip(background, &background, 1)


	for {
		webcam.Read(&curImg)
		gocv.Flip(curImg, &curImg, 1)
		gocv.CvtColor(curImg, &hsvImg, gocv.ColorBGRToHSV)


		gocv.InRangeWithScalar(hsvImg, lGreenMat, uGreenMat, &morphMask)

		kernal := gocv.NewMatWithSize(5,5, gocv.MatTypeCV8U)
		gocv.MorphologyEx(morphMask, &morphMask, gocv.MorphDilate, kernal)

		gocv.BitwiseNot(morphMask, &clearMask)

		imgSize := curImg.Size()
		imgRes := gocv.NewMatWithSize(imgSize[0], imgSize[1], gocv.MatTypeCV8U)
		backgroundRes := gocv.NewMatWithSize(background.Size()[0], background.Size()[1], gocv.MatTypeCV8U)

		curImg.CopyToWithMask(&imgRes, clearMask)
		background.CopyToWithMask(&backgroundRes, morphMask)

		gocv.AddWeighted(imgRes, 1, backgroundRes, 1, 0, &processedFrame)

		window.IMShow(processedFrame)
		window.WaitKey(1)

	}
}