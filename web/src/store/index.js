import { createStore } from 'vuex'
import connections from './modules/connections'
import extractions from './modules/extractions'

export default createStore({
  state: {},
  mutations: {},
  actions: {},
  modules: {
    connections,
    extractions
  }
})
