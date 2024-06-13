<template>
  <div class="subscriptions-container">
    <Spinner v-if="!items" :errorMessage="errorMessage"></Spinner>
    <div v-else class="subscriptions-table">
      <div class="table-header">
        <div class="header-item">Title</div>
        <div class="header-item">Thumbnail</div>
      </div>
      <div class="table-body">
        <div v-for="(item, index) in items" :key="index" class="table-row">
          <div class="table-item">{{ item.title }}</div>
          <div class="table-item">
            <img :src="item.thumbnailLink" :alt="item.title" width="100" height="100" />
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
      items: null,
      intervalId: null,
      errorMessage: null
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
      fetch(`${this.$apiEndpoint}/subscriptions`)
        .then((resp) => resp.json())
        .then((json) => {
          this.items = json
        })
        .catch((err) => {
          console.error(err)
          this.errorMessage = err
        })
    }
  }
}
</script>

<style scoped>
.subscriptions-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.subscriptions-table {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 800px;
  margin: 20px 0;
}

.table-header,
.table-row {
  display: flex;
  width: 100%;
}

.header-item,
.table-item {
  flex: 1;
  padding: 10px;
  text-align: left;
}

.table-header {
  background-color: #a6adc8;
  font-weight: bold;
}

.table-row:nth-child(even) {
  background-color: #cdd6f4;
}

/* dark mode child even */
.dark .table-row:nth-child(even) {
  background-color: #45475a;
}

.dark .table-header {
  background-color: #45475a;
}
</style>
