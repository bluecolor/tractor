import api from '@/api/connections'

const state = {}

const getters = {}

const actions = {
  getConnectionTypes({ commit }) {
    return api.getConnectionTypes()
  },
  getProviderTypes({ commit }) {
    return api.getProviderTypes()
  },
  createConnection({ commit }, payload) {
    return api.createConnection(payload)
  },
  getConnections({ commit }) {
    return api.getConnections()
  },
  getConnection({ commit }, { id }) {
    return api.getConnection(id)
  },
  deleteConnection({ commit }, id) {
    return api.deleteConnection(id)
  },
  updateConnection({ commit }, { id, ...payload }) {
    return api.updateConnection(id, payload)
  },
  resolveConnectorRequest({ commit }, payload) {
    return api.resolveConnectorRequest(payload)
  },
  getFields({ commit }, { connectionId, ...payload }) {
    return api.getFields(connectionId, payload)
  }
}

const mutations = {}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
