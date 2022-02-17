<template lang="pug">
PageTitle(title='New Connection')
a-form.create-connection(
  :model='state',
  name='basic',
  :label-col='{ span: 8 }',
  :wrapper-col='{ span: 8 }',
  autocomplete='off',
  @validate='onValidate',
  @submit='onSubmit'
)
  a-form-item(
    label='Name',
    name='name',
    :rules='[{ required: true, message: "Please enter a unique connection name" }]'
  )
    a-input(v-model:value='state.name')
  a-form-item(label='Type', name='type')
    a-select(ref='select', v-model:value='state.connectionTypeID')
      a-select-option(v-for='m in connectionTypes', :value='m.id') {{ m.name }}
  component(
    ref='config',
    v-if='state.connectionTypeID',
    :is='getConnectionTypeComponent()',
    :state.props='state.config'
  )
  a-form-item(:wrapper-col='{ offset: 8, span: 8 }')
    a-button(type='primary', html-type='submit', block, :loading='status.loading') Save
  a-form-item(v-if='status.type', :wrapper-col='{ offset: 8, span: 8 }')
    a-alert(:message='status.message', :type='status.type')
</template>

<script setup>
import { ref, markRaw, onBeforeMount } from 'vue'
import FileConnection from './types/file/FileConnection.vue'
import MySQLConnection from './types/mysql/MySQLConnection.vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import PageTitle from '@/components/PageTitle.vue'

const store = useStore()
const router = useRouter()
const status = ref({
  loading: false,
  type: undefined,
  message: undefined
})

const state = ref({
  name: '',
  connectionTypeID: undefined,
  config: undefined
})
const connectionTypes = ref([])
const components = markRaw({
  file: FileConnection,
  mysql: MySQLConnection
})
const config = ref(null)

const getConnectionTypeCode = (id) => {
  return connectionTypes.value.find((m) => m.id === id).code
}
const getComponentByID = (id) => {
  return components[getConnectionTypeCode(id)]
}
const getConnectionTypeComponent = () => {
  if (state.value.connectionTypeID) {
    return getComponentByID(state.value.connectionTypeID)
  }
}
const clearStatus = () => {
  status.value = {
    loading: false,
    type: undefined,
    message: undefined
  }
}

onBeforeMount(() => {
  status.value.loading = true
  store
    .dispatch('connections/getConnectionTypes')
    .then((types) => {
      connectionTypes.value = types
      if (types.length > 0) {
        state.value.connectionTypeID = types[0].id
      }
    })
    .finally(() => {
      status.value.loading = false
    })
})
const onValidate = () => {
  return true
}
const onSubmit = () => {
  state.value.config = config.value.getState()
  status.value.loading = true
  clearStatus()
  store
    .dispatch('connections/createConnection', state.value)
    .then(() => {
      message.success('Connection created successfully')
      router.push('/connections')
    })
    .catch((err) => {
      console.error(err)
      status.value = {
        type: 'error',
        message: err.message || 'Failed to create connection'
      }
    })
    .finally(() => {
      status.value.loading = false
    })
}
</script>

<style lang="scss" scoped>
.create-connection {
  padding-top: 20px;
}
</style>
