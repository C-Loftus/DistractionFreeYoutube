<template>
  <div>
    <spinner v-if="results.length === 0"></spinner>
    <div v-else class="videos">
      <table>
        <thead>
          <tr>
            <th>Title</th>
            <th>Description</th>
            <th>Video</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in results" :key="index">
            <td>{{ item.snippet.title }}</td>
            <td>{{ item.snippet.description }}</td>
            <td>
              <iframe
                width="320"
                height="180"
                :src="getEmbedUrl(item.id.videoId)"
                frameborder="0"
                allowfullscreen
                loading="lazy"
              ></iframe>
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
      error: ''
    }
  },
  components: {
    Spinner
  },
  updated() {
    this.getData()
  },
  methods: {
    getData() {
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
    },
    // Utility function to get embed URL
    getEmbedUrl(videoId) {
      return `https://www.youtube.com/embed/${videoId}`
    }
  }
}
</script>

<style scoped>
.videos {
  margin: 20px auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th,
td {
  padding: 8px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

th {
  background-color: #f2f2f2;
}

th:nth-child(even) {
  background-color: #fafafa;
}

.dark tbody tr:nth-child(even) {
  background-color: #333;
}

tbody tr:nth-child(even) {
  background-color: #fafafa;
}

/* dark mode child even */
.dark th:nth-child(even) {
  background-color: #333;
}

.dark th {
  background-color: #333;
}
</style>
