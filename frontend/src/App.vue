<template>
  <div id="app">
    <!-- Navbar -->
    <div class="navbar-content">
      <div class="left">
        <router-link to="/">Home</router-link>
        <router-link to="/subscriptions">Subscriptions</router-link>
        <router-link to="/playlists">Playlists</router-link>
      </div>
      <form @submit.stop.prevent="submit">
        <input type="text" placeholder="Search" v-model="searchData" />
      </form>
      <div class="right">
        <a v-if="!profilePic" :href="authRoute" target="_blank"> Login </a>
        <div v-else>
          <a :href="authRoute" target="_blank">
            <img
              :src="profilePic"
              alt="Your profile picture"
              title="Your profile picture"
              class="profilePic"
            />
          </a>
        </div>
        <img
          :src="isDark ? sun : moon"
          alt="Toggle dark mode"
          @click="toggleDarkMode"
          class="toggle"
        />
      </div>
    </div>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
import moon from './assets/moon.svg'
import sun from './assets/sun.svg'
import '@fontsource/work-sans'

export default {
  name: 'App',
  data() {
    return {
      authRoute: `${this.$apiEndpoint}/auth`,
      isDark: localStorage.getItem('isDark') === 'true',
      searchData: '',
      profilePic: '',
      moon: moon,
      sun: sun
    }
  },

  created() {
    this.getData()
  },

  methods: {
    toggleDarkMode() {
      this.isDark = !this.isDark
      localStorage.setItem('isDark', this.isDark)
      document.documentElement.classList.toggle('dark')
    },
    submit() {
      console.log(this.searchData)
      if (this.searchData) {
        this.searchData = this.searchData.trim()
        this.$router.push({ path: '/search', query: { query: this.searchData } })
      }
    },

    getData() {
      fetch(`${this.$apiEndpoint}/auth/profile`)
        .then((resp) => {
          if (!resp.ok) {
            throw new Error('Network response was not ok')
          }
          return resp.json()
        })
        .then((json) => {
          console.log(json)
          this.profilePic = json.items[0].snippet.thumbnails.default.url
        })
        .catch((err) => {
          this.error = 'Error fetching playlists: ' + err.message
          console.error(err)
        })
    }
  }
}
</script>

<style>
.navbar-content {
  top: 0;
  left: 0;
  color: #cdd6f4;
  /* fixed */
  padding: 10px;
  margin-right: 0px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  background-color: #313244;
  justify-content: space-between;
  align-items: center;
  /* round corners */
  border-radius: 10px;
}

html {
  background-color: #bac2de;
  color: #313244;
  font-family: 'Work Sans';
}

/* disable style for clicked link */

.navbar-content .left {
  display: flex;
  align-items: center;
}

.navbar-content .right {
  display: flex;
  align-items: center;
}

main {
  padding-top: 60px; /* Adjust according to navbar height */
}

.navbar-content a {
  color: #cdd6f4;
  text-decoration: none;
  margin-right: 10px;
  text-decoration: underline;
}

.dark {
  background-color: #1e1e2e;
  color: #cdd6f4;
}

.profilePic,
.toggle {
  height: 2em;
  width: auto;
}

input {
  padding: 10px;
  border-radius: 5px;
  border: none;
  outline: none;
  width: 300px;
  background-color: #e2e6f3;
}

h2 {
  color: #313244;
  text-decoration: underline;
}

.dark h2 {
  color: #cdd6f4;
}
</style>
