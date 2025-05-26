# 🎬 Cloud Video Editor

![Image](https://github.com/user-attachments/assets/1dc04ef0-30f7-4c33-8646-ec0af4811773)

**Cloud Video Editor**는 Vue 3 기반의 웹 비디오 편집 플랫폼으로,
브라우저 상에서 동영상을 직접 편집하고 원본 대비 60%로 요약된 영상을 생성할 수 있습니다.
사용자는 동영상을 업로드한 뒤, 재생 화면과 프레임 단위 타임라인을 통해 원하는 시점을 탐색하고 편집할 수 있습니다.

## 📌 주요 기능 요약

* 🔄 **동영상 업로드 및 미리보기**
* 🖼 **프레임 타임라인 기반 탐색**
* 🎞️ **영상 요약 모드 및 편집 조작**
* ⌛ **영상 프레임 추출 및 시각화**
* 🔧 **Drag & Drop + Timeline 이미지 클릭 재생 기능**
* 🎛 **로딩 애니메이션 및 키보드 단축키 지원**

---

## 📁 프로젝트 구조

```
sfproject/
├── public/                  # 정적 자산 (favicon 등)
├── src/
│   ├── main.js              # Vue 앱 진입점
│   ├── App.vue              # 루트 컴포넌트 (VideoEditor 포함)
│   ├── assets/              # 편집 도구 관련 이미지 아이콘
│   └── components/
│       ├── VideoEditor.vue     # 영상 업로드, 재생, 추출, 요약 등 핵심 기능
│       └── TimelinePanel.vue   # 프레임 기반 타임라인 표시 및 탐색
├── package.json             # 의존성 및 npm 스크립트
├── vue.config.js            # Vue 설정 (포트 등)
└── jsconfig.json            # IDE 자동완성 지원용 설정
```

---

## 📦 설치 및 실행

```bash
npm install      # 의존성 설치
npm run serve    # 개발 서버 실행 (기본 포트 80)
npm run build    # 프로덕션 빌드
```

---

## 🔧 사용된 주요 라이브러리

| 라이브러리                    | 설명                  |
| ------------------------ | ------------------- |
| `vue@3`                  | Vue.js 3 기반 프론트엔드   |
| `axios`                  | 비동기 요청              |
