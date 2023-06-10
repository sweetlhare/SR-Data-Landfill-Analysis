<script setup>
defineProps({ items: { type: Array, required: true } });
</script>

<template>
  <ul class="breadcrumbs">
    <li v-for="item of items" :key="item.name" class="breadcrumbs__item">
      <RouterLink
        :to="{ name: item.routeName, params: item.params, query: item.query }"
        class="breadcrumbs__link"
        :href="item.routeName ? '#' : null"
      >
        {{ item.title }}
      </RouterLink>
    </li>
  </ul>
</template>

<style lang="scss" scoped>
.breadcrumbs {
  display: flex;
  align-items: center;
  flex-wrap: wrap;

  &__item {
    position: relative;

    &:not(:last-child) {
      padding-right: 16px;

      &::after {
        content: '/';
        position: absolute;
        top: 50%;
        right: 6px;
        transform: translateY(-50%);
        color: gray;
      }
    }
  }

  &__link {
    font-size: 14px;
    line-height: 1;
    color: #a1a1a1;
    transition: color 0.3s ease-in-out;

    &[href]:hover,
    &[href]:focus {
      color: black;
    }

    .breadcrumbs__item:last-child > & {
      color: black;
    }
  }
}
</style>
