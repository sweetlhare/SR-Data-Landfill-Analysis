<script setup>
defineProps({ items: { type: Array, required: true } });
</script>

<template>
  <ul class="dictionary">
    <li v-for="item of items" :key="item.key" class="dictionary__item">
      <span class="dictionary__key">
        <strong>{{ item.key }}: </strong></span
      >
      <span class="dictionary__value"> {{ item.value }} </span>
    </li>
  </ul>
</template>

<style lang="scss" scoped>
.dictionary {
  margin-bottom: 18px;
  &__item:not(:last-child) {
    margin-bottom: 10px;
  }
  &__key {
    font-weight: 500;
    font-size: 14px;
    line-height: 18px;
  }

  &__value {
    font-weight: 400;
    font-size: 14px;
    line-height: 18px;
  }
}
</style>
