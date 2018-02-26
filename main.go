package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

var chanPN = make(chan string, 20)

func main() {

	path := flag.String("p", "./", "要改变尺寸的图片所在的文件夹")
	pwidth := flag.Uint("w", 60, "要改变的图片宽度")
	pHight := flag.Uint("h", 50, "要改变的图片高度")
	flag.Parse()

	pictureName(*path, *pwidth, *pHight)
	return
}

//遍历指定目录下的图片
func pictureName(path string, pWidth, pHight uint) error {

	var err error
	f, err := os.Stat(path)
	if err != nil {
		return fmtERR(err)
	}
	if !f.IsDir() {
		path, name := filepath.Split(path)

		if path == "" {
			path = "./"
		}

		PrictureSize(path, name, pWidth, pHight)
		return err
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return fmtERR(err)
	}

	for _, fi := range files {
		if fi.IsDir() {
		} else {
			str := filepath.Ext(fi.Name())
			if str != ".jpg" && str != ".JPG" {
				continue
			}
			PrictureSize(path, fi.Name(), pWidth, pHight)
		}
	}
	return err
}

//PrictureSize 改变传入的图片地址的图片的尺寸
func PrictureSize(savepath, pricturePath string, pWidth, pHight uint) error {

	//打开图片
	file, err := os.Open(pricturePath)
	if err != nil {
		return fmtERR(err)
	}

	//解码图片
	img, err := jpeg.Decode(file)
	if err != nil {
		return fmtERR(err)
	}
	file.Close()

	// 改变图片尺寸
	m := resize.Resize(pWidth, pHight, img, resize.Lanczos3)

	if err = testdir(savepath); err != nil {
		return err
	}
	//创建新图片
	out, err := os.Create(savepath + "/change/" + pricturePath)
	if err != nil {
		return fmtERR(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

	return err
}

func testdir(path string) error {
	path = path + "/change"
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0666)
		if err != nil {
			return fmtERR(err)
		}
	}
	return err
}

func fmtERR(err error) error {
	fmt.Println(err)
	return err
}
