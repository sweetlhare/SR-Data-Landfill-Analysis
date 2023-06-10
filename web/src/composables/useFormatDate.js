export function useFormatDate(string, short = true, includeTime = false) {
  const date = new Date(string),
    day = ('0' + date.getDate()).slice(-2),
    month = ('0' + (Number(date.getMonth()) + 1)).slice(-2),
    year = date.getFullYear();

  const time = ('0' + date.getHours()).slice(-2) + ':' + ('0' + date.getMinutes()).slice(-2);

  if (short) return day + '.' + month + '.' + year + (includeTime ? ' ' + time : '');

  const monthIndex = Number(date.getMonth()),
    monthTexts = [
      'января',
      'февраля',
      'марта',
      'апреля',
      'мая',
      'июня',
      'июля',
      'августа',
      'сентября',
      'октября',
      'ноября',
      'декабря',
    ];

  return day + ' ' + monthTexts[monthIndex] + ' ' + year + (includeTime ? ' ' + time : '');
}
