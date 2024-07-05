<template>
  <v-container class="fill-height">
    <v-responsive
      class="align-centerfill-height mx-auto"
      max-width="900"
    >
      <div class="my-auto text-center">
        <v-img
          height="250"
          src="@/assets/jetscript_logo.png"
        />

        <h1 class="text-h3 font-weight-bold">JetScript</h1>
        <p class="mt-8">
          Jetscript is a simple, modern, and powerful programming language that is designed to be easy to learn and use.
        </p>

        <div class="text-center w-50 py-8 mx-auto">
          <v-btn prepend-icon="mdi-github" text="Github" class="mr-3"/>
          <v-btn prepend-icon="mdi-book" text="Docs" color="primary" class="ml-3"/>
        </div>
      </div>




      <div class="py-8" />

    </v-responsive>
  </v-container>
</template>

<script setup>
import { ref, computed } from 'vue'

</script>

<style scoped>
.benthos_link {
  color: red;
  text-decoration: none;
}
</style>
