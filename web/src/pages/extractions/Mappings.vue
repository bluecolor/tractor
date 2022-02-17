<template lang="pug">
a-table(:dataSource='fieldMappings', :columns='columns')
</template>

<script setup>
import { markRaw, ref } from 'vue'
const columns = markRaw([
  {
    title: 'Source',
    dataIndex: ['sourceField', 'name'],
    key: 'sourceField'
  },
  {
    title: 'Type',
    dataIndex: ['sourceField', 'type'],
    key: 'sourceType'
  },
  {
    title: 'Target',
    dataIndex: ['targetField', 'name'],
    key: 'targetField'
  },
  {
    title: 'Type',
    dataIndex: ['targetField', 'type'],
    key: 'targetType'
  }
])
const fieldMappings = ref([])
</script>
