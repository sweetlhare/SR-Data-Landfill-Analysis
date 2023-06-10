-- +goose Up
-- +goose StatementBegin
INSERT INTO public.users (id, name, "position", "role", phone, email, "password") VALUES(1, 'Иванов Сергей Демидович', 'сотрудник', 'admin', '+79876543210', 'ivanov.s@mail.ru', '333');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
