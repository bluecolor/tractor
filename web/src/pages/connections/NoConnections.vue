<template lang="pug">
.no-connections
  .message
    | Yo don't have any connections yet.
  a-button(type='primary', size='large')
    | Add Connection
</template>

<script setup>
</script>

<style lang="scss" scoped>
.no-connections {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  .message {
    margin-bottom: 1rem;
    font-size: 1.2rem;
    color: #999;
  }
}
</style>
