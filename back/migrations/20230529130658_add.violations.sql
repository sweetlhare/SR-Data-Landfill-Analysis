-- +goose Up
-- +goose StatementBegin

INSERT INTO violations(id, description, default_status) VALUES
(0,'Отсутствие КПП, ворот или шлагбаума', true),
(1,'Отсутствие ванн для дезинфекции колес', true),
(2,'Отсутствие весового контроля', true),
(3,'Не выполняется пересыпка тела полигона', true),
(4,'Отсутствие ограждения по периметру полигона', true),
(5,'Отсутствие обводных каналов', true),
(6,'Заболачивание и подтопление полигона', true),
(7,'Выход фильтрата за границы отведенного земельного участка', true),
(8,'Наличие тлений и возгораний на теле полигона', true),
(9,'Разлет фракций на прилегающие территории', true),
(10,'Размещение фракций, не относящихся к ТКО', true),
(11,'Наличие птиц на участке', true),
(12,'Расположение объекта на расстоянии менее 15 км от действующего аэродрома', true)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
