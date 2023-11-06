package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"sfprojectBack/app"
	"strconv"
)

func futureFuntion() {
	//구글 로그인 기능 구현
	//slack을 통한 알림 기능 서버에 문제가 생겼을 때 나에게 알려주기
	//log파일 작성
	//서버에 배포하기
	//요약된 동영상 다운로드 하기
	//
}

func getFPS(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "0", "-of", "csv=p=0", "-select_streams", "v:0", "-show_entries", "stream=r_frame_rate", videoPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// r_frame_rate는 "num/den" 형태로 출력됩니다. 예: "30000/1001"
	re := regexp.MustCompile(`(\d+)/(\d+)`)
	matches := re.FindStringSubmatch(string(output))
	if matches == nil {
		return 0, fmt.Errorf("failed to parse FPS from output")
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	den, err := strconv.Atoi(matches[2])
	if err != nil {
		return 0, err
	}

	return float64(num) / float64(den), nil
}

func main() {
	fmt.Println("Hello world!")

	router := app.SetupRouter()

	log.Fatal(http.ListenAndServe(":8000", router))

	//rootVideoPath := "./originalVideo"
	// err := util.SegmentVideo(rootVideoPath)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//util.EditVideo(rootVideoPath)

	//서버 측 코드에서는 Content-Disposition 헤더로 파일 이름을 추출하고, 파일을 저장하는 데 사용

	// f, err := getFPS("originalVideo/original.mp4")
	// fmt.Println(err)
	// fmt.Println(f)

}
