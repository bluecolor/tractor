/* eslint-disable camelcase */

import api from '@/api/connection'

const state = {}

const getters = {}

const actions = {
  getConnectionTypes({ commit }) {
    return api.getConnectionTypes()
  },
  getProviderTypes({ commit }) {
    return api.getProviderTypes()
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
