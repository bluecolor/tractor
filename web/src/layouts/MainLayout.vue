<template lang="pug">
a-layout
  navbar
router-view
</template>

<script>
import Navbar from '../components/Navbar.vue'
export default {
  components: {
    Navbar
  },
  data() {
    return {
      collapsed: false
    }
  }
}
</script>
