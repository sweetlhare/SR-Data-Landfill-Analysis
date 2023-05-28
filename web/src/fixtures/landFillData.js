export default {
  id: 4375643,
  regionCode: 1,
  illegal: false,
  city: 'г. Санкт-Петербург',
  type: 'Полигон твердых бытовых отходов',
  address: 'ул. Экологичная, 3',
  coordinates: [55.78435, 37.54345],
  satellitePreview: '/img/383.png',
  manager: {
    id: 1,
    name: 'Иванова Алла Владимировна',
    position: 'Директор',
    phone: '+7 987 654-32-10',
    email: 'ivanova.a@mail.ru',
  },
  violations_сount: 5,
  surveys: [
    {
      id: 23,
      date: 'Fri May 01 2023 08:30:00 GMT+0300',
      rawImages: [
        'api/{sampleid}/img1.jpg',
        'api/{sampleid}/img2.jpg',
        'api/{sampleid}/img3.jpg',
        'api/{sampleid}/img4.jpg',
        'api/{sampleid}/img5.jpg',
      ],
      ai_images: ['/img/banner_main.png', '/img/banner_main.png'],
      audits: [
        {
          id: 1,
          date: 'Fri May 01 2023 12:30:00 GMT+0300',
          auditor: {
            id: 1,
            name: 'Иванов Сергей Демидович',
            position: 'сотрудник',
            phone: '+7 987 654-32-10',
            email: 'ivanov.s@mail.ru',
          },
          violations: [
            {
              id: 1,
              title: 'Отсутствие КПП, ворот или шлагбаума',
              status: false,
            },
            {
              id: 2,
              title: 'Не выполнение пересыпки ТКО инертными материалами',
              status: false,
            },
            {
              id: 3,
              title: 'Выход фильтрата за границы отведенного земельного участка',
              status: false,
            },
            {
              id: 4,
              title: 'Отсутствие сооружения дезинфекции транспорта',
              status: false,
            },
            {
              id: 5,
              title: 'Наличие тлений и возгораний на теле полигона',
              status: false,
            },
            {
              id: 6,
              title: 'Отсутствие весового контроля',
              status: false,
            },
            {
              id: 7,
              title: 'Отсутствие весового контроля',
              status: false,
            },
            {
              id: 8,
              title: 'Заболачивание и подтопление полигона',
              status: true,
            },
            {
              id: 9,
              title:
                'Наличие птиц на полигонах (по фотоснимкам), в том числе определение их количества',
              status: false,
            },
            {
              id: 10,
              title: 'Разлет фракций на прилегающие территории',
              status: true,
            },
            {
              id: 11,
              title: 'Отсутствие ограждения по периметру полигона',
              status: true,
            },
          ],
        },
      ],
    },
    {
      id: 456,
      date: 'Fri May 01 2023 08:30:00 GMT+0300',
      rawImages: [
        'api/{sampleid}/img1.jpg',
        'api/{sampleid}/img2.jpg',
        'api/{sampleid}/img3.jpg',
        'api/{sampleid}/img4.jpg',
        'api/{sampleid}/img5.jpg',
      ],
      ai_images: ['/img/banner_main.png', '/img/banner_main2.jpg'],
      audits: [
        {
          id: 45,
          date: 'Fri May 01 2023 12:30:00 GMT+0300',
          auditor: {
            id: 65,
            name: 'Иванов Сергей Демидович',
            position: 'сотрудник',
            phone: '+7 987 654-32-10',
            email: 'ivanov.s@mail.ru',
          },
          violations: [
            {
              id: 1,
              title: 'Отсутствие КПП, ворот или шлагбаума',
              status: false,
            },
            {
              id: 2,
              title: 'Не выполнение пересыпки ТКО инертными материалами',
              status: false,
            },
            {
              id: 3,
              title: 'Выход фильтрата за границы отведенного земельного участка',
              status: false,
            },
            {
              id: 4,
              title: 'Отсутствие сооружения дезинфекции транспорта',
              status: false,
            },
            {
              id: 5,
              title: 'Наличие тлений и возгораний на теле полигона',
              status: false,
            },
            {
              id: 6,
              title: 'Отсутствие весового контроля',
              status: false,
            },
            {
              id: 7,
              title: 'Отсутствие весового контроля',
              status: false,
            },
            {
              id: 8,
              title: 'Заболачивание и подтопление полигона',
              status: true,
            },
            {
              id: 9,
              title:
                'Наличие птиц на полигонах (по фотоснимкам), в том числе определение их количества',
              status: false,
            },
            {
              id: 10,
              title: 'Разлет фракций на прилегающие территории',
              status: true,
            },
            {
              id: 11,
              title: 'Отсутствие ограждения по периметру полигона',
              status: false,
            },
          ],
        },
      ],
    },
  ],
};
