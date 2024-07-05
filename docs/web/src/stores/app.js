// Utilities
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', () => {
  const index = ref({})

  async function load() {
    const respIndex = await fetch('docs/__index.json')
    index.value = await respIndex.json()
  }

  return {
    index,
    load
  }
})
