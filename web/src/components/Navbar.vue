<template lang="pug">
a-layout-header.header
  a-row(type='flex', justify='center', style='height: 100%')
    a-col.container(:span='16')
      img(src='@/assets/tractor.png', alt='Tractor', style='height: 32px')
      .right
        a-menu(v-model:selectedKeys='selectedKeys', mode='horizontal')
          a-menu-item(key='connections', @click='$router.push("/connections")')
            template(#icon='')
              ApiOutlined
            | Connections
          a-menu-item(key='extractions', @click='$router.push("/extractions")')
            template(#icon='')
              VerticalAlignBottomOutlined
            | Extractions
          a-menu-item(key='sessions')
            template(#icon='')
              ThunderboltOutlined
            | Sessions
          a-sub-menu(key='more')
            template(#icon='')
              MoreOutlined
            template(#title='') More
            a-menu-item.menu-item(key='more.settings')
              template(#icon)
                SettingOutlined
              span(style='display: inline-block') Settings
            a-menu-item.menu-item(key='more.help')
              template(#icon)
                QuestionOutlined
              span(style='display: inline-block') Help
            a-menu-item.menu-item(key='more.source', @click='onSource')
              template(#icon)
                GithubOutlined
              span(style='display: inline-block') Source Code
</template>

<script setup>
import TractorIcon from '@/components/icons/TractorIcon.vue'
import {
  VerticalAlignBottomOutlined,
  MoreOutlined,
  ApiOutlined,
  ThunderboltOutlined,
  SettingOutlined,
  GithubOutlined,
  QuestionOutlined
} from '@ant-design/icons-vue'
import { ref } from 'vue'
const selectedKeys = ref([''])
const current = ref('')
const onSource = () => {
  window.open('https://github.com/bluecolor/tractor', '_blank')
}
</script>

<style lang="scss" scoped>
.header {
  position: fixed;
  width: 100%;
  z-index: 1;
  display: flex;
  flex-direction: column;
  background: #fff;
  padding: 0;
  border-bottom: 1px solid #e9e9e9;
  .container {
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    .right {
      display: flex;
      flex-grow: 1;
      flex-direction: row;
      align-items: center;
      justify-content: flex-end;
      .actions {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: flex-end;
        .action {
          margin-left: 1rem;
        }
      }
    }
  }
}
</style>

<style>
.header-right-menu {
  min-width: 200px;
}
.header-right-menu .menu-item .ant-dropdown-menu-item-icon {
  margin-right: 16px;
}
</style>
