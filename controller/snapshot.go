package controller
import (
	"fmt"
  "os"
	"log"
   ffmpeg "github.com/u2takey/ffmpeg-go"
	 "bytes"
	 "strings"
   "github.com/disintegration/imaging"
)

func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	 buf := bytes.NewBuffer(nil)
	 err = ffmpeg.Input(videoPath).
			 Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
			 Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
			 WithOutput(buf, os.Stdout).
			 Run()
	 if err != nil {
			 log.Fatal("生成缩略图失败：", err)
			 return "", err
	 }
	 img, err := imaging.Decode(buf)
	 if err != nil {
			 log.Fatal("生成缩略图失败：", err)
			 return "", err
	 }
	 err = imaging.Save(img, snapshotPath+".png")
	 if err != nil {
			 log.Fatal("生成缩略图失败：", err)
			 return "", err
	 }
	 names := strings.Split(snapshotPath, "/")
	 snapshotName = names[len(names)-1] + ".png"
	 return snapshotName, nil
}