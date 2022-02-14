<template lang="pug">
a-layout-header.header
  a-row(type='flex', justify='center', style='height: 100%')
    a-col.container(:span='16')
      img(src='@/assets/tractor.png', alt='Tractor', style='height: 32px')
      .right
        .actions
          a-tooltip(title='Connections')
            a-button.action(shape='circle', @click='$router.push("connections")')
              template(#icon)
                ApiOutlined
          a-tooltip(title='Extractions')
            a-button.action(shape='circle')
              template(#icon)
                VerticalAlignBottomOutlined
          a-tooltip(title='Sessions')
            a-button.action(shape='circle')
              template(#icon)
                ThunderboltOutlined
          a-tooltip(title='More')
            a-dropdown.action.drop-down(placement='bottomRight', trigger='click')
              a-button(shape='circle')
                template(#icon)
                  MoreOutlined
              template(#overlay='')
                a-menu.header-right-menu
                  a-menu-item.menu-item
                    template(#icon)
                      SettingOutlined
                    span(style='display: inline-block') Settings
                  a-menu-item.menu-item
                    template(#icon)
                      QuestionOutlined
                    span(style='display: inline-block') Help
                  a-menu-item.menu-item
                    template(#icon)
                      GithubOutlined
                    span(style='display: inline-block') Source
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
</script>

<style lang="scss" scoped>
.header {
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
