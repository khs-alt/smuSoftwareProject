<template>
  <div class="videoEditor">
    <!--상단부-->
    <div class="top">
      <!-- 프로젝트(파일 업로드) 패널 -->
      <div class="project-panel">
        <!--파일 업로드 영역-->
        <div @drop="handleFileSelect" @dragover.prevent class="drop-area">
          <input type="file" ref="fileInput" @change="handleFileSelect" accept="video/*" style="display: none;">
          <p>동영상 파일을 여기로 드래그 앤 드롭하세요.</p>
          <p>또는</p>
          <button @click="openFileInput">파일 선택</button>
        </div>
      </div>
      <!--비디오 뷰어 패널-->
      <div class="viewer-panel">
        <video class="video" ref="videoPlayer" @ended="hanldeVideoEnded"></video>
      </div>
    </div>
    <!--하단부-->
    <div class="bottom">
      <!--재생시간-->
      <div class="video-time">
        <span>{{ currentTime }}</span> / <span>{{ Duration }}</span>
      </div>
      <!--컨트롤 패널-->
      <div class="control-panel">
        <!--컨트롤 버튼 영억-->
        <img src="@/assets/back.png" @click.prevent="skipTime(-5)" />
        <img :src="isPlaying ? require('@/assets/pause.png') : require('@/assets/play.png')"
          @click.prevent="togglePlayPause" />
        <img src="@/assets/forward.png" @click.prevent="skipTime(5)" />
      </div>
      <!--타임라인 패널-->
      <div class="timeline-panel">
        <!--타임라인 영역-->
        <div class="timeline">
          <div class="timeline-label" v-for="(img, i) in imgArr" :key="i"
            :style="{ left: (i * timelineImageWidth) + 'px' }">
            <span class="time-label">{{ img.time }}s</span>
          </div>
          <div class="timeline-image" v-for="(img, i) in imgArr" :key="i"
            :style="{ left: (i * timelineImageWidth) + 'px' }"  @click="seekToTime(img.time)">
            <img :src="img.url" alt="Image">
          </div>
          <div class="timeline-marker">
            <img src="@/assets/playhead.jpeg"
              :style="{ left: markerPosition + 'px' }" 
              draggable="true" 
              @dragstart="startDragPlayhead" 
              @drag="dragPlayhead"
              @dragover="stopDrag"/>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

//import axios from 'axios';

export default {
  data() {
    return {
      selectedFile: null,   // 파일 선택 여부
      imgArr: [], //추출한 이미지 저장할 배열
      videoLoaded: false,
      isPlaying: false, //동영상 재생 여부
      markerPosition: 0, //마커 위치
      timelineImageWidth: 100, // 타임라인 이미지 간격 조절
      lastImageTime: 0,
      playheadInterval: null, // 재생 헤드 업데이트를 위한 인터벌 변수 추가
      intervalId: null,
      currentTime: 0,
      Duration: 0,
    };
  },

  methods: {
    //drag & drop
    handleFileSelect(event) {
      if (event.type === 'drop') {
        event.preventDefault();
        this.selectedFile = event.dataTransfer.files[0];
      } else if (event.type === 'change') {
        this.selectedFile = event.target.files[0];
      }
      this.loadVideo();
    },

    //버튼 눌러 파일 업로드
    openFileInput() {
      this.$refs.fileInput.click();
    },

    //재생헤드 이미지 추출
    async loadVideo() {
      this.$refs.videoPlayer.src = URL.createObjectURL(this.selectedFile);
      this.imgArr = [];
      this.extractImages();
    },

    updateTime() {
      const video = this.$refs.videoPlayer;
      this.CurrentTime = video.currentTime;
    },

    getTotalTime() {
      const video = this.$refs.videoPlayer;
      this.Duration = video.duration;
    },

    togglePlayPause() {
      const video = this.$refs.videoPlayer;
      if (video.readyState >= 2) {
        if (video.paused) {
          video.play();
          this.isPlaying = true;

          this.playheadInterval = setInterval(this.playhead, 16);
        } else {
          video.pause();
          this.isPlaying = false;

          clearInterval(this.playheadInterval);
        }
      }
    },

    skipTime(seconds) {
      const video = this.$refs.videoPlayer;
      video.currentTime += seconds;
      video.play();
      this.isPlaying = true;
    },

    hanldeVideoEnded() {
      this.isPlaying = false;
    },

    extractImages() {
      const video = this.$refs.videoPlayer;
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');

      video.addEventListener('loadeddata', async () => {
        const duration = Math.floor(video.duration);

        for (let time = 0; time < duration + 1; time++) {
          video.currentTime = time;
          await video.play(); //
          canvas.width = video.videoWidth;
          canvas.height = video.videoHeight;
          ctx.drawImage(video, 0, 0, canvas.width, canvas.height);

          const imageDataURL = canvas.toDataURL('image/jpeg'); // 이미지 파일로 변환
          this.imgArr.push({ url: imageDataURL, time });
          this.lastImageTime  = time;
          video.pause();
        }
      });
      this.isPlaying = true;
    },

    startDragPlayhead(event) {
      // const dragImg= new Image();
      // dragImg.src = require('@/assets/playhead.jpeg');
      // event.dataTransfer.setDragImage(dragImg, 0, 0);
      this.isDragging = true;
      this.startX = event.clientX;
      this.initialTime = this.$refs.videoPlayer.currentTime;
      this.initialMarkerPosition = this.markerPosition;
      this.$refs.videoPlayer.pause();
    },

    dragPlayhead(event) {
      if (this.isDragging) {
        event.preventDefault();
        const deltaX = event.clientX - this.startX;
        const video = this.$refs.videoPlayer;
        const newTime = this.initialTime + (deltaX / video.clientWidth) * video.duration;
        if (newTime >= 0 && newTime <= video.duration) {
          video.currentTime = newTime;
          this.markerPosition = (newTime / video.duration) * this.timelineWidth;
        }
      }
    },

    stopDrag() {
      if (this.isDragging) {
        this.isDragging = false;
        this.markerPosition = this.initialMarkerPosition;
        this.isPlaying = true;
        this.$refs.videoPlayer.play();
      }
    },

    playhead() {
      if(this.isPlaying) {
        const video = this.$refs.videoPlayer;
        this.markerPosition = (video.currentTime / video.duration) * 1550;
      }
    },

    seekToTime(time) {
      const video = this.$refs.videoPlayer;
      video.currentTime = time;
      this.isPlaying = true;
      video.play();
    },
  },

  mounted() {
    const video = this.$refs.videoPlayer;
    setInterval(this.updateTime, 1000);
    video.addEventListener('loadedmetadata', this.getTotalTime);
  },
  beforeUnmount() {
    const video = this.$refs.videoPlayer;
    video.removeEventListener('loadedmetadata', this.getTotalTime);
  },
};
</script>

<style>
.top {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
}

.video {
  background-color: #444444;
  width: 360px;
}

.drop-area,
.viewer-panel {
  margin: 100px;
  text-align: center;
}

.drop-area {
  border: 2px dashed #444444;
  padding: 20px;
  cursor: pointer;
}

.video-time {
  text-align: center;
  color: black;
  font-size: 16px;
  margin-top: 10px;
}

.control-panel img {
  width: 50px;
}

.timeline-label {
  position: absolute;
  top: -20px;
  /* 초 표시를 타임라인 위에 표시하도록 조절 */
  left: 0;
  transform: translateX(-50%);
  /* 가운데 정렬 */
  white-space: nowrap;
  /* 텍스트가 넘칠 경우 줄 바꿈 방지 */
  color: black;
  /* 텍스트 색상 */
}

.time-label {
  font-size: 12px;
  /* 글꼴 크기 조절 */
}

.timeline-panel {
  position: relative;
  width: 100%;
  background-color: #444444;
}

.timeline {
  position: relative;
  height: 100px;
  /* width: 80%; */
  overflow:auto;
  white-space: nowrap;
}

.timeline-marker {
  position:absolute;
  cursor: grab;
  height: 100%;
}

.timeline-marker img {
  position: absolute;
  height: 100%;
  left: 0px;
}

.timeline-image {
  position: absolute;
  /*이미지 요소 나란히 배열*/
}

.timeline-image img {
  width: 100px;
  height: 100px;
}
</style>