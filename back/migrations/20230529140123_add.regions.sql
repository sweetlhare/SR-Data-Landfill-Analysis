-- +goose Up
-- +goose StatementBegin
INSERT INTO public.regions (id, "name") VALUES(4, 'Волгоградская область');
INSERT INTO public.regions (id, "name") VALUES(5, 'Вологодская область');
INSERT INTO public.regions (id, "name") VALUES(6, 'Забайкальский край');
INSERT INTO public.regions (id, "name") VALUES(7, 'Кабардино-Балкарская Республика');
INSERT INTO public.regions (id, "name") VALUES(8, 'Калужская область');
INSERT INTO public.regions (id, "name") VALUES(9, 'Карачаево-Черкесская Республика');
INSERT INTO public.regions (id, "name") VALUES(10, 'Кировская область');
INSERT INTO public.regions (id, "name") VALUES(11, 'Костромская область');
INSERT INTO public.regions (id, "name") VALUES(12, 'Курганская область');
INSERT INTO public.regions (id, "name") VALUES(13, 'Магаданская область');
INSERT INTO public.regions (id, "name") VALUES(14, 'Нижегородская область');
INSERT INTO public.regions (id, "name") VALUES(15, 'Новосибирская область');
INSERT INTO public.regions (id, "name") VALUES(16, 'Омская область');
INSERT INTO public.regions (id, "name") VALUES(17, 'Пензенская область');
INSERT INTO public.regions (id, "name") VALUES(18, 'Приморский край');
INSERT INTO public.regions (id, "name") VALUES(19, 'Псковская область');
INSERT INTO public.regions (id, "name") VALUES(20, 'Республика Бурятия');
INSERT INTO public.regions (id, "name") VALUES(21, 'Республика Дагестан');
INSERT INTO public.regions (id, "name") VALUES(22, 'Республика Ингушетия');
INSERT INTO public.regions (id, "name") VALUES(23, 'Республика Коми');
INSERT INTO public.regions (id, "name") VALUES(24, 'Республика Марий Эл');
INSERT INTO public.regions (id, "name") VALUES(25, 'Республика Мордовия');
INSERT INTO public.regions (id, "name") VALUES(26, 'Республика Северная Осетия- Алания');
INSERT INTO public.regions (id, "name") VALUES(27, 'Республика Татарстан');
INSERT INTO public.regions (id, "name") VALUES(28, 'Республика Хакасия');
INSERT INTO public.regions (id, "name") VALUES(29, 'Самарская область');
INSERT INTO public.regions (id, "name") VALUES(30, 'Свердловская область');
INSERT INTO public.regions (id, "name") VALUES(31, 'Ставропольский край');
INSERT INTO public.regions (id, "name") VALUES(32, 'Тверская область');
INSERT INTO public.regions (id, "name") VALUES(33, 'Тульская область');
INSERT INTO public.regions (id, "name") VALUES(34, 'Чеченская Республика');
INSERT INTO public.regions (id, "name") VALUES(35, 'Чувашская Республика');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
