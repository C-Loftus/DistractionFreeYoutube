<template>
  <div class="spinner" v-if="!errorMessage">
    <div class="bounce1"></div>
    <div class="bounce2"></div>
    <div class="bounce3"></div>
  </div>
  <p v-else>{{ errorMessage }}</p>
</template>

<script>
export default {
  // eslint-disable-next-line vue/multi-word-component-names
  name: 'spinner',

  props: {
    errorMessage: String
  }
}
</script>

<style>
.spinner {
  margin: 10px auto 0;
  width: 90px;
  text-align: center;
}
p {
  text-align: center;
}

.spinner > div {
  width: 28px;
  height: 28px;
  background-color: #313244;

  border-radius: 100%;
  display: inline-block;
  -webkit-animation: sk-bouncedelay 1.4s infinite ease-in-out both;
  animation: sk-bouncedelay 1.4s infinite ease-in-out both;
}

.spinner .bounce1 {
  -webkit-animation-delay: -0.32s;
  animation-delay: -0.32s;
}

.spinner .bounce2 {
  -webkit-animation-delay: -0.16s;
  animation-delay: -0.16s;
}

@-webkit-keyframes sk-bouncedelay {
  0%,
  80%,
  100% {
    -webkit-transform: scale(0);
  }
  40% {
    -webkit-transform: scale(1);
  }
}

@keyframes sk-bouncedelay {
  0%,
  80%,
  100% {
    -webkit-transform: scale(0);
    transform: scale(0);
  }
  40% {
    -webkit-transform: scale(1);
    transform: scale(1);
  }
}

.dark .spinner > div {
  background-color: #cdd6f4;
}
</style>
