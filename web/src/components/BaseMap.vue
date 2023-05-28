<script setup>
import { yandexMap, ymapMarker } from 'vue-yandex-maps';
defineProps({ mapCenter: { type: Array, required: true } });

const settings = {
  apiKey: 'e70506a0-5c6f-499f-9fdf-ffc9170481e9',
  lang: 'ru_RU',
  coordorder: 'latlong',
  enterprise: false,
  version: '2.1',
};
</script>

<template>
  <yandex-map
    :coords="mapCenter"
    zoom="16"
    :controls="['zoomControl', 'rulerControl']"
    :settings="settings"
  >
    <ymap-marker :coords="mapCenter" />
  </yandex-map>
</template>

<style lang="scss" scoped></style>
