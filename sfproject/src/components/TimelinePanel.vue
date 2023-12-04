<template>
    <div class="timeline-panel">
        <div class="timeline">
            <div class="timeline-image" v-for="(img, i) in imgArr" :key="i"
                :style="{ left: (i * timelineImageWidth) + 'px' }" @click="seekToTime(img.time)">
                <img :src="img.url" alt="Image">
            </div>
        </div>
    </div>
</template>

<script>
export default {
    props: ['imgArr'],
    setup(props) {
        console.log(props.imgArr);
    },
    data() {
        return {
            timelineImageWidth: 100,
        }
    },

    methods: {
        //이미지 클릭 재생 시간 변경
        sendTime(time) {
            this.$emit(time);
        },
    },
}

</script>

<style>
.timeline {
    position: relative;
    height: 100%;
}

.timeline-image {
    position: absolute;
    top: 45px;
}

.timeline-image img {
    width: 100px;
    height: 100px;
}
</style>