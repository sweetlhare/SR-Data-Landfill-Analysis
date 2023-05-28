<script setup>
import { useRoute } from 'vue-router';
import BaseHeader from './components/BaseHeader.vue';

const route = useRoute();
</script>

<template>
  <BaseHeader v-if="route.name !== 'not-found' && route.name !== 'sign-in'"></BaseHeader>
  <router-view></router-view>
</template>
