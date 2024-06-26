<template>
  <div :key="forceUpdateKey">
    <Spinner v-if="results.length === 0" :errorMessage="error"></Spinner>
    <div v-else class="videos">
      <table>
        <thead>
          <tr>
            <th class="video-column">Video</th>
            <th>Title</th>
            <th>Description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in results" :key="index">
            <td class="video-column">
              <iframe
                ref="iframe"
                width="640"
                height="360"
                :src="getEmbedUrl(item.id.videoId)"
                frameborder="0"
                allowfullscreen
                loading="lazy"
                v-if="loadPlayer[index]"
                @click="toggleFullScreen($event)"
              ></iframe>
              <!-- get image -->
              <img
                v-else
                :src="item.snippet.thumbnails.high.url"
                :alt="item.snippet.title"
                @click="togglePlayer(index)"
              />
            </td>
            <td>{{ item.snippet.title }}</td>
            <td>{{ item.snippet.description }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import Spinner from '@/components/Spinner.vue'

export default {
  data() {
    return {
      results: [],
      error: '',
      loadPlayer: {},
      forceUpdateKey: 0
    }
  },
  components: {
    Spinner
  },
  mounted() {
    this.getData()
  },
  methods: {
    getData() {
      this.forceUpdateKey += 1
      try {
        fetch(`${this.$apiEndpoint}/search?query=${this.$route.query.query}`)
          .then((resp) => {
            if (!resp.ok) {
              throw new Error('Network response was not ok')
            }
            return resp.json()
          })
          .then((json) => {
            this.results = json // Assuming 'items' is the key containing the array of videos
          })
          .catch((err) => {
            this.error = 'Error fetching playlists: ' + err.message
          })
      } catch {
        this.error = 'Error fetching playlists: '
      }
    },
    // Utility function to get embed URL
    getEmbedUrl(videoId) {
      return `https://www.youtube.com/embed/${videoId}`
    },
    togglePlayer(itemId) {
      this.$set(this.loadPlayer, itemId, true)
    },
    toggleFullScreen(event) {
      const iframe = event.target
      if (iframe.requestFullscreen) {
        iframe.requestFullscreen()
      } else if (iframe.mozRequestFullScreen) {
        // Firefox
        iframe.mozRequestFullScreen()
      } else if (iframe.webkitRequestFullscreen) {
        // Chrome, Safari and Opera
        iframe.webkitRequestFullscreen()
      } else if (iframe.msRequestFullscreen) {
        // IE/Edge
        iframe.msRequestFullscreen()
      }
    }
  }
}
</script>

<style scoped>
.videos {
  margin: 20px auto;
}

table {
  width: 60%;
  border-collapse: collapse;
  margin: 0 auto;
}

th,
td {
  padding: 10px;
  text-align: left;
}

th.video-column,
td.video-column {
  width: 60%;
}
th {
  text-align: center;
}
</style>
