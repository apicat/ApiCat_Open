import { onBeforeUnmount, onMounted, onUpdated, shallowRef, watch } from 'vue'
import { EVENT_CODE } from '@apicat/shared'

import type { Ref } from 'vue'
import { off, on } from '../utils'
import type { Nullable } from '../utils'
import type TreeStore from './tree-store'
import { useNamespace } from '@/hooks/useNamespace'

interface UseKeydownOption {
  el$: Ref<HTMLElement>
}
export function useKeydown({ el$ }: UseKeydownOption, store: Ref<TreeStore>) {
  const namespaceRef = ref('el')
  const ns = useNamespace('tree', namespaceRef)

  const treeItems = shallowRef<Nullable<HTMLElement>[]>([]) as any
  const checkboxItems = shallowRef<Nullable<HTMLElement>[]>([]) as any

  onMounted(() => {
    initTabIndex()
    on(el$.value, 'keydown', handleKeydown as any)
  })

  onBeforeUnmount(() => {
    off(el$.value, 'keydown', handleKeydown as any)
  })

  onUpdated(() => {
    treeItems.value = Array.from(el$.value.querySelectorAll('[role=treeitem]'))
    checkboxItems.value = Array.from(el$.value.querySelectorAll('input[type=checkbox]'))
  })

  watch(checkboxItems, (val) => {
    val.forEach((checkbox: any) => {
      checkbox.setAttribute('tabindex', '-1')
    })
  })

  const handleKeydown = (ev: KeyboardEvent): void => {
    const currentItem = ev.target as HTMLElement
    if (!currentItem.className.includes(ns.b('node')))
      return
    const code = ev.code
    treeItems.value = Array.from(el$.value.querySelectorAll(`.${ns.is('focusable')}[role=treeitem]`))
    const currentIndex = treeItems.value.indexOf(currentItem)
    let nextIndex
    if ([EVENT_CODE.up, EVENT_CODE.down].includes(code)) {
      ev.preventDefault()
      if (code === EVENT_CODE.up) {
        nextIndex = currentIndex === -1 ? 0 : currentIndex !== 0 ? currentIndex - 1 : treeItems.value.length - 1
        const startIndex = nextIndex
        while (true) {
          if (store.value.getNode(treeItems.value[nextIndex].dataset.key).canFocus)
            break
          nextIndex--
          if (nextIndex === startIndex) {
            nextIndex = -1
            break
          }
          if (nextIndex < 0)
            nextIndex = treeItems.value.length - 1
        }
      }
      else {
        nextIndex = currentIndex === -1 ? 0 : currentIndex < treeItems.value.length - 1 ? currentIndex + 1 : 0
        const startIndex = nextIndex
        while (true) {
          if (store.value.getNode(treeItems.value[nextIndex].dataset.key).canFocus)
            break
          nextIndex++
          if (nextIndex === startIndex) {
            nextIndex = -1
            break
          }
          if (nextIndex >= treeItems.value.length)
            nextIndex = 0
        }
      }
      nextIndex !== -1 && treeItems.value[nextIndex].focus()
    }
    if ([EVENT_CODE.left, EVENT_CODE.right].includes(code)) {
      ev.preventDefault()
      currentItem.click()
    }
    const hasInput = currentItem.querySelector('[type="checkbox"]') as Nullable<HTMLInputElement>
    if ([EVENT_CODE.enter, EVENT_CODE.space].includes(code) && hasInput) {
      ev.preventDefault()
      hasInput.click()
    }
  }

  const initTabIndex = (): void => {
    treeItems.value = Array.from(el$.value.querySelectorAll(`.${ns.is('focusable')}[role=treeitem]`))
    checkboxItems.value = Array.from(el$.value.querySelectorAll('input[type=checkbox]'))
    const checkedItem = el$.value.querySelectorAll(`.${ns.is('checked')}[role=treeitem]`)
    if (checkedItem.length) {
      checkedItem[0].setAttribute('tabindex', '0')
      return
    }
    treeItems.value[0]?.setAttribute('tabindex', '0')
  }
}
