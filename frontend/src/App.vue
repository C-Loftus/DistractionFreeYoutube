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
        <a :href="authRoute" target="_blank"> Login </a>
        <a href="#" @click="toggleDarkMode">Toggle Dark Mode</a>
      </div>
    </div>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      authRoute: `${this.$apiEndpoint}/auth`,
      isDark: localStorage.getItem('isDark') === 'true',
      searchData: ''
    }
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
    }
  }
}
</script>

<style>
.navbar-content {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  color: #fff;
  padding: 10px 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  background-color: #333;

  justify-content: space-between; /* Align items to the start and end of the container */
  align-items: center;

  text-decoration: underline;
}
/* Align left */
.navbar-content .left {
  display: flex;
  align-items: center;
}

/* Align right */
.navbar-content .right {
  display: flex;
  align-items: center;
}

/* Main content styles */
main {
  padding-top: 60px; /* Adjust according to navbar height */
}

.navbar-content a {
  color: #fff;
  text-decoration: none;
  margin-right: 10px;
}

.dark h1 {
  color: #fff;
}

/* dark background */
.dark {
  background-color: #000;
  color: #fff;
}
</style>
