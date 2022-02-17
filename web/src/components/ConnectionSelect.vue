<template lang="pug">
a-select(
  v-model:value='value',
  placeholder='Select source connection',
  style='width: 100%',
  @change='onChange'
)
  a-select-option(v-for='conn in connections', :key='conn.id', :value='conn.id')
    .div(style='display: flex; justify-content: space-between')
      span(style='display: inline-block')
        | {{ conn.name }}
      span(style='display: inline-block; color: #ccc; font-size: 0.8rem')
        | {{ conn.connectionType.name }}
</template>

<script setup>
import { defineProps, ref, defineEmits } from 'vue'

const props = defineProps({
  connections: {
    type: Array,
    default: () => []
  },
  value: {
    type: Number,
    default: null
  }
})
const emit = defineEmits(['change'])
const onChange = (id) => {
  emit('change', id)
}
</script>
