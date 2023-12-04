# 파일을 읽을 경로를 지정합니다. 이 경로를 적절하게 수정하세요.
file_path = '/path/to/your/text/file.txt'

try:
    # 파일을 읽기 모드('r')로 엽니다.
    with open(file_path, 'r', encoding='utf-8') as file:
        # 파일 내용을 읽어옵니다.
        file_contents = file.read()
        
        # 파일 내용을 화면에 출력합니다.
        print("파일 내용:")
        print(file_contents)

except FileNotFoundError:
    print(f"경로에 파일이 없습니다: {file_path}")
except Exception as e:
    print(f"파일을 읽는 동안 오류가 발생했습니다: {str(e)}")
