import api from '@/api/extractions'

const state = {}

const getters = {}

const actions = {
  getExtractions({ commit }) {
    return api.getExtractions()
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
