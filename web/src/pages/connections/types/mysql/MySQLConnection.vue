<template lang="pug">
a-form-item(label='Host', name='host', :rules='[{ required: true, validator: validateHost }]')
  a-input(v-model:value='state.host')
a-form-item(label='Port', name='port', :wrapper-col='{ span: 8 }')
  a-input-number(v-model:value='state.port', :min='0', :max='65535', style='width: 100%')
a-form-item(label='Database', name='database')
  a-input(v-model:value='state.database')
a-form-item(label='Username', name='username')
  a-input(v-model:value='state.username')
a-form-item(label='Password', name='password')
  a-input-password(v-model:value='state.password', autoComplete='off')
a-form-item(label='Params', name='params')
  a-input(v-model:value='state.params')
</template>

<script setup>
import { defineProps } from 'vue'

const validateHost = (rule) => {
  if (props.state.host.length === 0) {
    return Promise.reject(new Error(rule.message))
  } else {
    return Promise.resolve()
  }
}

const props = defineProps({
  state: {
    type: Object,
    default: () => ({
      host: '',
      port: 3306
    })
  }
})
</script>
