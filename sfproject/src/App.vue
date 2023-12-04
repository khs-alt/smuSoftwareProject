<template>
  <div id="app">
    <VideoEditor/>
  </div>
</template>

<script>
import VideoEditor from "./components/VideoEditor.vue";
export default {
  data() {
    return {
    };
  },

  components: {
    VideoEditor,
  },

};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}
</style>
