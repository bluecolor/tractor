<template lang="pug">
a-layout.main-layout
  a-layout-header
    navbar
  a-layout-content
    a-row.content(type='flex', justify='center', style='height: 100%')
      a-col.dyn-container(:span='16')
        router-view
</template>

<script setup>
import Navbar from '../components/Navbar.vue'
</script>

<style lang="scss" scoped>
.ant-layout-header {
  background: #fff;
  padding: 0;
}
.main-layout {
  .ant-layout-content {
    .content {
      padding-top: 20px;
      padding-bottom: 20px;
      height: 100%;
      background: #eeeeee;
      min-height: calc(100vh - 64px);
      .dyn-container {
        border-radius: 4px;
        background: #fff;
        padding: 24px;
      }
    }
  }
}
</style>
