<template lang="pug">
PageTitle(title='New Extraction')
a-divider(orientation='left') Connections
a-row.extraction-connections(type='flex', :gutter='16')
  a-col.source-conn(:flex='1', :xs='24', :sm='24', :md='12', :lg='12')
    ConnectionSelect(
      :connections='connections',
      :value='sourceConnectionId',
      @change='onSourceConnectionChange'
    )
    a-collapse(ghost)
      a-collapse-panel.options(v-if='sourceOptComponent', key='1', header='Options')
        component(:is='sourceOptComponent', :connection='getConnection(sourceConnectionId)')
  a-col.target-conn(:flex='1', :xs='24', :sm='24', :md='12', :lg='12')
    ConnectionSelect(
      :connections='connections',
      :value='targetConnectionId',
      @change='onTargetConnectionChange'
    )
    a-collapse(v-model:activeKey='activeKey', ghost)
      a-collapse-panel.options(v-if='targetOptComponent', key='1', header='Options')
        component(:is='targetOptComponent')
a-divider(orientation='left') Mappings
.actions(style='display: flex; justify-content: space-between')
  .left
    a-dropdown
      template(#overlay)
        a-menu
          a-menu-item(key='1')
            | Delete selected
          a-menu-item(key='2')
            | Add new
          a-menu-item(key='3')
            | Sync source fields
          a-menu-item(key='4')
            | Sync target fields
      a-button
        | Actions
        DownOutlined
  .right
    a-select(ref='select', v-model:value='extractionMode', style='width: 120px')
      a-select-option(value='append') Append
      a-select-option(value='truncate') Truncate

Mappings(style='margin-top: 10px')
</template>

<script setup>
import PageTitle from '@/components/PageTitle.vue'
import { useStore } from 'vuex'
import { ref, onBeforeMount, markRaw, watch, shallowRef } from 'vue'
import MySQLOptions from './options/mysql/MySQLOptions.vue'
import CsvOptions from './options/file/csv/CsvOptions.vue'
import Mappings from './Mappings.vue'
import { DownOutlined } from '@ant-design/icons-vue'
import ConnectionSelect from '@/components/ConnectionSelect.vue'

const store = useStore()
const connections = ref([])
const sourceConnectionId = ref(null)
const targetConnectionId = ref(null)
const activeKey = ref(null)
const extractionMode = ref('append')
const selectedMappings = ref([])
const sourceOptComponent = shallowRef(null)
const targetOptComponent = shallowRef(null)
const status = ref({
  loading: false,
  error: null,
  message: null
})
const options = markRaw({
  mysql: MySQLOptions,
  file: {
    csv: CsvOptions
  }
})

const clearStatus = () => {
  status.value = {
    loading: false,
    error: null,
    message: null
  }
}
const getComponent = (id) => {
  const connection = connections.value.find((c) => c.id === id)
  const { connectionType, config } = connection
  const { code } = connectionType
  if (connection) {
    if (code === 'file') {
      const { format } = config
      return options.file[format]
    } else {
      return options[code]
    }
  } else {
    return undefined
  }
}
const getConnection = (id) => {
  return connections.value.find((c) => c.id === id)
}
const onSourceConnectionChange = (id) => {
  sourceConnectionId.value = id
  sourceOptComponent.value = getComponent(id)
}
const onTargetConnectionChange = (id) => {
  targetConnectionId.value = id
  targetOptComponent.value = getComponent(id)
}
const fetchConnections = () => {
  status.value.loading = true
  store
    .dispatch('connections/getConnections')
    .then((result) => {
      connections.value = result
    })
    .finally(() => {
      clearStatus()
    })
}

onBeforeMount(() => {
  fetchConnections()
})
</script>

<style lang="scss" scoped>
.options {
  margin-top: 10px;
}
</style>

<style lang="scss">
.extraction-connections {
  .ant-collapse-ghost > .ant-collapse-item > .ant-collapse-content > .ant-collapse-content-box,
  .ant-collapse > .ant-collapse-item > .ant-collapse-header {
    padding-left: 0;
    padding-right: 0;
  }
}
</style>
