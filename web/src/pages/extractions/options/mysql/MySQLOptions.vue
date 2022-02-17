<template lang="pug">
a-form.connection-options(:label-col='{ span: 8 }', :wrapper-col='{ span: 16 }', labelAlign='left')
  a-form-item(label='Database', name='database')
    a-input(v-model:value='database')
  a-form-item(
    label='Table',
    name='table',
    :rules='[{ required: true, message: "Please enter table name" }]'
  )
    a-input(v-model:value='database')
</template>

<script setup>
import { ref, defineProps } from 'vue'

defineProps({
  database: {
    type: String,
    default: ''
  },
  table: {
    type: String,
    default: ''
  }
})
</script>

<style lang="scss">
.connection-options {
  .ant-form-item {
    margin-bottom: 12px;
  }
}
</style>
