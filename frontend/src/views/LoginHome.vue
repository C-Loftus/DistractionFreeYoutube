<template>
  <h1>Simple Youtube</h1>
  <img :src="playbutton" alt="Play button" />
</template>
<script>
import playbutton from '../assets/playbutton.svg'

export default {
  data() {
    return {
      playbutton
    }
  }
}
</script>
<style scoped>
h1 {
  text-align: center;
}

img {
  display: block;
  margin: 0 auto;
  width: 200px;
  height: 200px;
}
</style>
