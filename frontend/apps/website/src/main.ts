import 'uno.css'
import '@/styles/main.scss'

import { createApp } from 'vue'
import { pinia, elementPlus, initI18n, errorHandler, setupPiniaWithRouter } from './plugins'
import clipboardHelper from '@/components/ClipboardHelper'
import router, { setupRouterFilter } from '@/router'
import limitInput from '@/directives/LimitInput'

import App from './App.vue'

const run = async () => {
  const app = createApp(App)
  const i18n = await initI18n()
  app.use(i18n)
  app.use(pinia)
  app.use(elementPlus)
  app.use(errorHandler)
  app.use(router)
  app.use(clipboardHelper)
  app.use(limitInput)
  setupRouterFilter(router)
  setupPiniaWithRouter(router)

  app.mount('#app')
}

run()
