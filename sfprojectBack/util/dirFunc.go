package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadDirFile(path string) (string, error) {
	// 읽을 디렉토리 경로 설정
	directoryPath := path

	// 디렉토리 내의 파일 목록 가져오기
	fileInfos, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Println("디렉토리를 읽는 중 오류 발생:", err)
		return "", err
	}

	// 각 파일의 내용 읽어오기
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			// 디렉토리인 경우 무시
			continue
		}

		// 파일 경로 생성
		filePath := filepath.Join(directoryPath, fileInfo.Name())

		// 파일 내용 읽기
		_, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("파일을 읽는 중 오류 발생: %s\n", err)
			return "", err
		}
		fmt.Printf("파일 이름: %s\n", fileInfo.Name())
		//fmt.Printf("파일 내용:\n%s\n", string(fileContent))
		return fileInfo.Name(), nil

	}
	return "", nil
}

func DeleteFilesInFolder(folderPath string) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 디렉토리는 건너뜁니다.
		if info.IsDir() {
			return nil
		}

		// 파일 삭제
		err = os.Remove(path)
		if err != nil {
			return err
		}

		fmt.Printf("파일 삭제: %s\n", path)
		return nil
	})

	return err
}
