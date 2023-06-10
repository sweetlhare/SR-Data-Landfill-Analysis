<script setup>
// rows - массив объектов с произвольными полями, в том числе обязательным полем id
// columns - массив вида [{ title, field, template? }]
// При отсутствии поля template у элемента массива в ячейке показывается содержимое row[column.field]
// В противном случае в ячейке показывается содержимое шаблона, переданного в шаблон с именем column.template
defineProps({ rows: { type: Array, required: true }, columns: { type: Array, required: true } });
</script>

<template>
  <table class="table">
    <thead class="table__header">
      <tr class="table__header-row">
        <th v-for="column of columns" :key="column.field" class="table__cell table__cell_header">
          {{ column.title }}
        </th>
      </tr>
    </thead>
    <tbody class="table__body">
      <tr v-for="(row, index) of rows" :key="row.id" class="table__row">
        <td v-for="column of columns" :key="column.field" class="table__cell">
          <slot :name="column.template" :value="row[column.field] ?? index">{{
            row[column.field]
          }}</slot>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<style lang="scss" scoped>
.table {
  border-spacing: 0;
  background-color: #f7f7f7;
  border-radius: 8px;
  overflow-y: auto;
  max-height: 507px;
  min-width: 638px;
}

.table__cell {
  padding: 12px 10px;
  font-size: 14px;
  line-height: 21px;
  color: #7d7d7d;
}
.table__cell_header {
  &:first-child {
    border-top-left-radius: 8px;
  }
  &:last-child {
    border-top-right-radius: 8px;
  }
}

.table__row {
  margin: 0 16px;

  &:hover {
    background-color: var(--white);
  }

  &:last-child {
    & .table__cell:first-child {
      border-bottom-left-radius: 8px;
    }
    & .table__cell:last-child {
      border-bottom-right-radius: 8px;
    }
  }
}

.table__header {
  border-top-right-radius: 8px;
  border-top-left-radius: 8px;
  background-color: #ececec;
  font-style: normal;
  font-weight: 400;
  text-align: left;
  vertical-align: bottom;
}
</style>
