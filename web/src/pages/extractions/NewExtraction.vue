<template lang="pug">
PageTitle(title='New Extraction')
a-divider(orientation='left') Connections
a-row.extraction-connections(type='flex', :gutter='16')
  a-col.source-conn(:flex='1', :xs='24', :sm='24', :md='12', :lg='12')
    ConnectionSelect(:connections='connections', v-model='sourceConnectionId')
    a-collapse(v-if='sourceDatasetComponent', ghost)
      a-collapse-panel.options(key='1', header='Options')
        component(
          :is='sourceDatasetComponent',
          v-model='sourceDataset',
          :connection='sourceConnection'
        )
  a-col.target-conn(:flex='1', :xs='24', :sm='24', :md='12', :lg='12')
    ConnectionSelect(:connections='connections', v-model='targetConnectionId')
//-     a-collapse(v-model:activeKey='activeKey', ghost)
//-       a-collapse-panel.options(v-if='targetOptComponent', key='1', header='Options')
//-         component(:is='targetDatasetComponent')
//- a-divider(orientation='left') Mappings
//- .actions(style='display: flex; justify-content: space-between')
//-   .left
//-     a-dropdown
//-       template(#overlay)
//-         a-menu
//-           a-menu-item(key='1')
//-             | Delete selected
//-           a-menu-item(key='2')
//-             | Add new
//-           a-menu-item(key='3')
//-             | Sync source fields
//-           a-menu-item(key='4')
//-             | Sync target fields
//-       a-button
//-         | Actions
//-         DownOutlined
//-   .right
//-     a-select(ref='select', v-model:value='extractionMode', style='width: 120px')
//-       a-select-option(value='append') Append
//-       a-select-option(value='truncate') Truncate

//- Mappings(style='margin-top: 10px')
</template>

<script>
import { mapActions } from 'vuex'
import PageTitle from '@/components/PageTitle.vue'
import ConnectionSelect from '@/components/ConnectionSelect.vue'
import MySQLDataset from './datasets/mysql/MySQLDataset.vue'
import CsvDataset from './datasets/file/csv/CsvDataset.vue'

export default {
  components: {
    PageTitle,
    ConnectionSelect
  },
  data() {
    return {
      sourceDataset: null,
      targetDataset: null,
      connections: [],
      sourceConnectionId: null,
      targetConnectionId: null,
      activeKey: null,
      extractionMode: 'append',
      datasetComponents: {
        mysql: MySQLDataset,
        file: {
          csv: CsvDataset
        }
      }
    }
  },
  computed: {
    sourceDatasetComponent() {
      console.log(this.sourceConnectionId)
      if (!this.sourceConnectionId) return null
      return this.getComponent(this.sourceConnectionId)
    },
    targetDatasetComponent() {
      if (!this.sourceConnectionId) return null
      return this.getComponent(this.targetConnectionId)
    }
  },
  methods: {
    ...mapActions('connections', ['getConnections']),
    onSourceConnectionChange() {},
    onTargetConnectionChange() {},
    fetchConnections() {
      this.getConnections()
        .then((result) => {
          this.connections = result
        })
        .finally(() => {})
    },
    getComponent(connectionId) {
      if (!connectionId) {
        return null
      }
      const connection = this.connections.find((c) => c.id === connectionId)
      const { connectionType, config } = connection
      const { code } = connectionType
      if (connection) {
        if (code === 'file') {
          const { format } = config
          return this.datasetComponents.file[format]
        } else {
          return this.datasetComponents[code]
        }
      }
    }
  },
  beforeMount() {
    this.fetchConnections()
  }
}
</script>

<script setup>
// import PageTitle from '@/components/PageTitle.vue'
// import { useStore } from 'vuex'
// import { ref, onBeforeMount, markRaw, watch, shallowRef, computed } from 'vue'
// import MySQLDataset from './datasets/mysql/MySQLDataset.vue'
// import CsvDataset from './datasets/file/csv/CsvDataset.vue'
// import Mappings from './Mappings.vue'
// import { DownOutlined } from '@ant-design/icons-vue'
// import ConnectionSelect from '@/components/ConnectionSelect.vue'

// const store = useStore()
// const connections = ref([])
// const sourceConnectionId = ref(null)
// const targetConnectionId = ref(null)
// const activeKey = ref(null)
// const extractionMode = ref('append')
// const selectedMappings = ref([])
// const sourceOptComponent = shallowRef(null)
// const targetOptComponent = shallowRef(null)
// const status = ref({
//   loading: false,
//   error: null,
//   message: null
// })
// const datasetComponent = markRaw({
//   mysql: MySQLDataset,
//   file: {
//     csv: CsvDataset
//   }
// })
// const clearStatus = () => {
//   status.value = {
//     loading: false,
//     error: null,
//     message: null
//   }
// }

// const getComponent = (id) => {
//   const connection = connections.value.find((c) => c.id === id)
//   const { connectionType, config } = connection
//   const { code } = connectionType
//   if (connection) {
//     if (code === 'file') {
//       const { format } = config
//       return datasetComponent.file[format]
//     } else {
//       return datasetComponent[code]
//     }
//   }
// }
// const getConnection = (id) => {
//   return connections.value.find((c) => c.id === id)
// }
// const onSourceConnectionChange = (id) => {
//   sourceConnectionId.value = id
//   sourceOptComponent.value = getComponent(id)
// }
// const onTargetConnectionChange = (id) => {
//   targetConnectionId.value = id
//   targetOptComponent.value = getComponent(id)
// }
// const fetchConnections = () => {
//   status.value.loading = true
//   store
//     .dispatch('connections/getConnections')
//     .then((result) => {
//       connections.value = result
//     })
//     .finally(() => {
//       clearStatus()
//     })
// }

// onBeforeMount(() => {
//   fetchConnections()
// })
//
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
