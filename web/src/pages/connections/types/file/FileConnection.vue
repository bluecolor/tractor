<template lang="pug">
a-form-item(label='Provider', name='provider')
  a-select(v-model:value='state.provider')
    a-select-option(v-for='m in providerTypes', :value='m.code') {{ m.name }}
a-form-item(label='Format', name='format')
  a-select(v-model:value='state.format')
    a-select-option(value='csv') CSV
a-form-item(
  label='Path',
  name='path',
  :rules='[{ required: true, message: "Please enter file path", validator: validatePath }]'
)
  a-input(v-model:value='state.path')
</template>

<script setup>
import { defineProps, onBeforeMount, ref, defineExpose } from 'vue'
import { useStore } from 'vuex'

const providerTypes = ref([])
const store = useStore()

onBeforeMount(() => {
  store.dispatch('connections/getProviderTypes').then((types) => {
    providerTypes.value = types
  })
})

const validatePath = (rule, value) => {
  if (props.state.path.length === 0) {
    return Promise.reject(new Error('Please enter file path'))
  } else {
    return Promise.resolve()
  }
}

const props = defineProps({
  state: {
    type: Object,
    default: () => ({
      path: '',
      provider: 's3',
      format: 'csv'
    })
  }
})

const getState = () => {
  return props.state
}
defineExpose({ getState })
</script>
