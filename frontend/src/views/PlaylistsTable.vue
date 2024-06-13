<template>
  <div>
    <spinner v-if="playlists.length === 0" :errorMessage="error"></spinner>
    <div v-else class="playlist-container">
      <div v-for="(playlist, index) in playlists" :key="index" class="playlist">
        <router-link :to="`/playlists/${playlist.title}`">
          <h2>{{ playlist.title }}</h2>
          <img :src="playlist.thumbnailLink" :alt="playlist.title" />
        </router-link>

        <div class="videos"></div>
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
  justify-content: center;
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
