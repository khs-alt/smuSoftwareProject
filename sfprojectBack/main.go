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

func fadein() {
	// 비디오 클립 페이드인, 페이드아웃 효과 추가 및 자르기
	clip1 := exec.Command("ffmpeg", "-i", "./output/output_01.mp4", "-ss", "10", "-to", "20", "-vf", "fade=t=in:st=0:d=1,fade=t=out:st=9:d=1", "clip1.mp4")
	clip2 := exec.Command("ffmpeg", "-i", "./output/output_02.mp4", "-ss", "10", "-to", "20", "-vf", "fade=t=in:st=0:d=1,fade=t=out:st=9:d=1", "clip2.mp4")
	clip3 := exec.Command("ffmpeg", "-i", "./output/output_03.mp4", "-ss", "20", "-to", "30", "-vf", "fade=t=in:st=0:d=1,fade=t=out:st=9:d=1", "clip3.mp4")

	// 각 클립을 처리합니다.
	if err := clip1.Run(); err != nil {
		log.Fatalf("Error processing clip1: %v", err)
	}
	if err := clip2.Run(); err != nil {
		log.Fatalf("Error processing clip2: %v", err)
	}
	if err := clip3.Run(); err != nil {
		log.Fatalf("Error processing clip3: %v", err)
	}
	// if err := clip4.Run(); err != nil {
	// 	log.Fatalf("Error processing clip4: %v", err)
	// }

	// 비디오 클립 합치기
	concatCmd := exec.Command("ffmpeg", "-y", "-f", "concat", "-safe", "0", "-i", "concat_list.txt", "-c", "copy", "combined_temp.mp4")

	// concat_list.txt는 clip1.mp4, clip2.mp4, clip3.mp4, clip4.mp4 파일명을 포함해야 합니다.
	if err := concatCmd.Run(); err != nil {
		log.Fatalf("Error concatenating clips: %v", err)
	}

	// 오디오 파일 처리 및 볼륨 조절
	audioCmd := exec.Command("ffmpeg", "-i", "./originalVideo/original.mp4", "-af", "afade=t=in:st=0:d=1,volume=0.1", "audio.mp3")

	if err := audioCmd.Run(); err != nil {
		log.Fatalf("Error processing audio: %v", err)
	}

	// 최종 비디오에 오디오 합치기
	finalCmd := exec.Command("ffmpeg", "-i", "combined_temp.mp4", "-i", "audio.mp3", "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "combined.mp4")

	if err := finalCmd.Run(); err != nil {
		log.Fatalf("Error combining video and audio: %v", err)
	}
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

	//fadein()
}
