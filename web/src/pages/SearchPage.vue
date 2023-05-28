<script setup>
import BaseHeader from '../components/BaseHeader.vue';

const options = [
  'Показать полигоны с нарушениями на карте',
  'Показать полигоны без нарушений на карте',
  'Показать незаконные свалки',
];

const violations = [
  {
    violation_title: 'Отсутствие КПП, ворот или шлагбаума',
    violation_description:
      'Отсутствие КПП, ворот или шлагбаума в ограде полигона у производственно-бытового здания.',
  },
  {
    violation_title: 'Отсутствие сооружения дезинфекции транспорта',
    violation_description:
      'Отсутствие дезинфицирующей установки для очистки ходовой части мусоровозов – железобетонные ванны.',
  },
];
const images = [];
</script>

<template>
  <main class="layout search-page">
    <BaseHeader></BaseHeader>

    <h2 class="h2">Поиск по параметрам</h2>

    <div class="input_block">
      <label class="input_block__label" for="name">Выберите тип поиска</label>
      <select id="name" class="input_block__input">
        <option v-for="opt in options" :value="opt">{{ opt }}</option>
      </select>
    </div>

    <ul class="search-page__list">
      <li class="search-page__list-item">
        <div class="input_block">
          <label class="input_block__label" for="name">Поиск по регионам</label>
          <input id="name" class="input_block__input" placeholder="Введите регион" />
        </div>
      </li>
      <li class="search-page__list-item">
        <div class="input_block">
          <label class="input_block__label" for="name">Поиск по городам</label>
          <input id="name" class="input_block__input" placeholder="Введите город" />
        </div>
      </li>
      <li class="search-page__list-item">
        <div class="input_block">
          <label class="input_block__label" for="name">Поиск по адресу</label>
          <input id="name" class="input_block__input" placeholder="Введите адрес" />
        </div>
      </li>
      <li class="search-page__list-item">
        <div class="input_block">
          <label class="input_block__label" for="name">Поиск по координатам</label>
          <input id="name" class="input_block__input" placeholder="55°45′07″ N, 37°36′56″  E" />
        </div>
      </li>
      <li class="search-page__list-item">
        <div class="input_block">
          <label class="input_block__label" for="name">Поиск по названию полигона</label>
          <input id="name" class="input_block__input" placeholder="Введите название полигона" />
        </div>
      </li>
    </ul>

    <div class="search-page__buttons">
      <button class="btn btn_primary-outlined">Показать полигоны на карте</button>
      <button class="btn btn_primary-outlined">CRM для полигонов РФ</button>
    </div>

    <section class="search-page__content">
      <img src="" alt="" />
      <div class="content__polygon">
        <h4 class="h4">Название полигона</h4>
        <p><strong>Адрес:</strong> адрес полигона</p>
        <p><strong>Наименование объектом размещения отходов:</strong> lorem ipsum dolor sit</p>
        <p><strong>Координаты:</strong> 55°45′07″ N, 37°36′56″ E</p>

        <ul class="polygon__violations-status">
          <p><strong>Статусы нарушений:</strong></p>

          <li v-for="item in violations" class="violations-status__item">
            {{ item.violation_title }}
          </li>
        </ul>

        <ul class="polygon__violations-list">
          <p>
            <strong>Нарушения: {{ violations.length }}</strong>
          </p>

          <li v-for="item in violations" class="violations-list__item">
            {{ item.violation_description }}
          </li>
        </ul>

        <ul class="polygon__images-list">
          <li v-for="item in images" class="images-list__item">
            <img :src="item" alt="Изображение полигона" />
          </li>
        </ul>
      </div>
    </section>
  </main>
</template>

<style lang="scss">
.input_block {
  width: 359px;
}
.search-page {
  &__list {
    display: flex;
    align-items: center;
    margin-top: 18px;

    .input_block {
      width: 240px;

      &__input {
        border-radius: 0;
      }
    }

    &-item:not(:last-child) {
      margin-right: 24px;
    }

    &-item:first-child .input_block__input {
      border-top-left-radius: 12px;
      border-bottom-left-radius: 12px;
    }
    &-item:last-child .input_block__input {
      border-top-right-radius: 12px;
      border-bottom-right-radius: 12px;
    }
  }

  &__buttons {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 28px;

    button:first-child {
      margin-right: 48px;
    }
  }
  &__content {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 50px;
    margin-top: 28px;

    .content__polygon {
      border: 1px solid #cfcfcf;
      border-radius: 12px;
      padding: 28px;

      h4 {
        margin-bottom: 10px;
      }

      p {
        margin-bottom: 10px;
      }

      .violations-status__item {
        font-weight: 400;
        font-size: 14px;
        line-height: 18px;
        color: #ffffff;
        padding: 4px 10px;
        background: #ff4141;
        border-radius: 50px;
        margin-bottom: 10px;
        width: fit-content;
      }

      .polygon__violations-list {
        list-style-type: decimal;

        .violations-list__item {
          margin-bottom: 10px;
          font-weight: 400;
          font-size: 14px;
          line-height: 18px;
        }
      }
      .polygon__images-list {
        margin-top: 28px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        flex-wrap: wrap;

        .images-list__item {
          img {
            width: 110px;
            height: 110px;
          }
        }
      }
    }
  }
}
</style>
