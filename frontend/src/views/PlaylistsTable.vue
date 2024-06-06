<template>
  <div>
    <spinner v-if="playlists.length === 0"></spinner>
    <div v-else class="playlist-container">
      <div v-for="(playlist, index) in playlists" :key="index" class="playlist">
        <h2>{{ playlist.title }}</h2>
        <div class="videos">
          <div v-for="(video, vIndex) in playlist.videos" :key="vIndex" class="video">
            <iframe
              width="560"
              height="315"
              :src="getEmbedUrl(video.id)"
              frameborder="0"
              loading="lazy"
              allowfullscreen
            ></iframe>
            <p>{{ video.title }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Spinner from '@/components/Spinner.vue'

export default {
  data() {
    return {
      playlists: [],
      error: ''
    }
  },
  components: {
    Spinner
  },
  created() {
    this.getData()
  },
  methods: {
    getData() {
      fetch(`${this.$apiEndpoint}/playlists`)
        .then((resp) => {
          if (!resp.ok) {
            throw new Error('Network response was not ok')
          }
          return resp.json()
        })
        .then((json) => {
          this.playlists = json
        })
        .catch((err) => {
          this.error = 'Error fetching playlists: ' + err.message
        })
    },
    getEmbedUrl(videoId) {
      return `https://www.youtube.com/embed/${videoId}`
    }
  }
}
</script>

<style scoped>
.playlist-container {
  display: flex;
  flex-wrap: wrap;
}

.playlist {
  margin-bottom: 20px;
}

.videos {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
}

.video {
  margin: 10px;
  text-align: center;
}

.video p {
  margin-top: 5px;
}

h2 {
  text-align: center;
}
</style>
