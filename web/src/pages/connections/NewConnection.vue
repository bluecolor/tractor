<template lang="pug">
.title
  | New Connection
a-form.create-connection(
  :model='state',
  name='basic',
  :label-col='{ span: 8 }',
  :wrapper-col='{ span: 8 }',
  autocomplete='off',
  @validate='onValidate'
)
  a-form-item(
    label='Name',
    name='name',
    :rules='[{ required: true, message: "Please enter a unique connection name" }]'
  )
    a-input(v-model:value='state.name')
  a-form-item(label='Type', name='type')
    a-select(ref='select', v-model:value='state.type')
      a-select-option(value='file') File
      a-select-option(value='mysql') MySQL
  component(:is='components[state.type]', :state.path='state.config')
  a-form-item(:wrapper-col='{ offset: 8, span: 8 }')
    a-button(type='primary', html-type='submit', block) Save
</template>

<script setup>
import { ref, markRaw } from 'vue'
import FileConnection from './types/file/FileConnection.vue'
import MySQLConnection from './types/mysql/MySQLConnection.vue'

const onValidate = () => {
  return true
}

const components = markRaw({
  file: FileConnection,
  mysql: MySQLConnection
})

const state = ref({
  name: '',
  type: 'file',
  config: undefined
})
</script>

<style lang="scss" scoped>
.create-connection {
  padding-top: 20px;
}
.title {
  font-size: 20px;
  font-weight: 500;
  padding-bottom: 20px;
  color: #333;
}
</style>
