INSERT INTO graph_node (id, type, text) VALUES
(1, 'question', 'У тебя есть действующий паспорт?'),
(2, 'action', 'Подать заявление на замену паспорта'),
(3, 'action', 'Оформить новый паспорт'),
(4, 'end', 'Процесс завершён');

INSERT INTO graph_edge (from_node, to_node, condition) VALUES
(1, 2, 'yes'),
(1, 3, 'no'),
(2, 4, 'next'),
(3, 4, 'next');