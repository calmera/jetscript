<template>
  <v-container class="fill-height">
    <pre>{{index}}</pre>
  </v-container>
</template>

<script setup>
import {useAppStore} from "@/stores/app";

const appStore = useAppStore();
const {index} = storeToRefs(appStore);
</script>
