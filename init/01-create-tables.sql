-- Удаляем существующие таблицы (осторожно!)
DROP TABLE IF EXISTS task_users;
DROP TABLE IF EXISTS tasks;

-- Создание таблицы tasks с UUID
CREATE TABLE tasks (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       title TEXT NOT NULL,
                       hours INTEGER CHECK (hours > 0),
                       start_date DATE NOT NULL,
                       deadline DATE NOT NULL
);

-- Создание связующей таблицы task_users (user_id остается integer)
CREATE TABLE task_users (
                            task_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
                            user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
                            PRIMARY KEY (task_id, user_id)
);

-- Вставка тестовых задач с UUID
INSERT INTO tasks (id, title, hours, start_date, deadline) VALUES
                                                               ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Разработать API', 40, '2024-01-15', '2024-01-25'),
                                                               ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'Создать дизайн интерфейса', 30, '2024-01-10', '2024-01-20'),
                                                               ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'Написать документацию', 20, '2024-01-18', '2024-01-28')
    ON CONFLICT (id) DO NOTHING;

-- Связывание пользователей с задачами
INSERT INTO task_users (task_id, user_id) VALUES
                                              ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 1),
                                              ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 2),
                                              ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 2),
                                              ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 4),
                                              ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 4)
    ON CONFLICT DO NOTHING;