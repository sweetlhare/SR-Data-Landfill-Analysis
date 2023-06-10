<script setup></script>

<template>
  <header class="layout">
    <img src="/icons/logo.svg" alt="Логотип" />
    <div class="header__user-info">
      <img src="" alt="" />
      <div>
        <p>Сергей Иванов</p>
        <button>Личный кабинет</button>
      </div>
    </div>
  </header>
</template>

<style lang="scss" scoped>
header {
  display: flex;
  justify-content: space-between;
  padding: 36px;

  &__user-info {
    img {
      border-radius: 50%;
      width: 75px;
      height: 75px;
      margin-right: 15px;
    }

    p {
      font-weight: 500;
      font-size: 18px;
      line-height: 23px;
    }
    span {
      font-weight: 400;
      font-size: 14px;
      line-height: 18px;
    }
  }
}
</style>
