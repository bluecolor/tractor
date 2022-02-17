<template lang="pug">
a-row(type='flex', justify='space-between', style='padding-bottom: 30px')
  a-col
    .title
      | {{ title }}
  a-col(:span='4', style='text-align: right')
    a-button(type='link', @click='back', style='color: #8c8c8c')
      template(#icon)
        component(:is='icon()')
</template>

<script setup>
import { defineProps } from 'vue'
import { CloseOutlined } from '@ant-design/icons-vue'

const back = () => {
  window.history.back()
}

defineProps({
  title: {
    type: String,
    default: ''
  },
  icon: {
    type: Function,
    default: () => CloseOutlined
  }
})
</script>

<style lang="scss" scoped>
.title {
  font-size: 20px;
  font-weight: 500;
  color: #333;
}
</style>
