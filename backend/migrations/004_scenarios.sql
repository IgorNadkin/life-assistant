-- 001_init.sql уже создал таблицу scenario с колонками (id, name, description, created_at).
-- Раньше здесь стоял CREATE TABLE IF NOT EXISTS с другим набором колонок (title, start_node_id) —
-- он молча не выполнялся, т.к. таблица уже существовала, и репозиторий пытался
-- писать/читать несуществующие колонки title/start_node_id. Дополняем таблицу вместо пересоздания.
ALTER TABLE scenario ADD COLUMN IF NOT EXISTS title TEXT;
ALTER TABLE scenario ADD COLUMN IF NOT EXISTS start_node_id INT;

-- Колонка name из 001_init.sql нигде не используется (domain.Scenario её не содержит),
-- но имеет NOT NULL без DEFAULT, из-за чего любой INSERT из scenario_repo.go падал
-- с "null value in column name violates not-null constraint". Убираем как мёртвую колонку.
ALTER TABLE scenario DROP COLUMN IF EXISTS name;