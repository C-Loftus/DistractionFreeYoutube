<template>
  <div>
    <Spinner v-if="results.length === 0" :errorMessage="error"></Spinner>
    <div v-else class="videos">
      <table>
        <thead>
          <tr>
            <th>Title</th>
            <th>Video</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in results" :key="index">
            <td>{{ item.title }}</td>
            <td>
              <iframe
                width="640"
                height="360"
                :src="getEmbedUrl(item.id)"
                frameborder="0"
                allowfullscreen
                loading="lazy"
                v-if="loadPlayer[index]"
              ></iframe>
              <img
                v-else
                :src="item.thumbnailLink"
                :alt="item.title"
                @click="togglePlayer(index)"
              />
            </td>
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

  // computed: {
  //   filteredVideos() {
  //     return this.results.filter((item) => item.title.includes(this.$route.query.query))
  //   }
  // },

  methods: {
    getData() {
      try {
        fetch(`${this.$apiEndpoint}/playlists`)
          .then((resp) => {
            if (!resp.ok) {
              throw new Error('Network response was not ok')
            }
            return resp.json()
          })
          .then((json) => {
            console.log(json)
            this.results = json.filter(
              (playlist) => playlist.title === this.$route.params.name
            )[0].videos
            console.log(JSON.stringify(this.results))
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
      this.loadPlayer[itemId] = true
    }
  }
}
</script>

<style scoped>
.videos {
  margin: 20px auto;
}

table {
  width: 80%;
  margin: 0 auto;
  border-collapse: collapse;
}

th,
td {
  padding: 8px;
  text-align: left;
}

th {
  background-color: #9399b2;
  text-align: center;
}

.dark tbody tr:nth-child(even) {
  background-color: #313244;
}

tbody tr:nth-child(even) {
  background-color: #a6adc8;
}

.dark th {
  background-color: #313244;
}
</style>
