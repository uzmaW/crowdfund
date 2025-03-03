
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    goal FLOAT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    user_id INTEGER REFERENCES users(id)
);


INSERT INTO projects (id, title, description, goal, start_date, end_date, user_id)
VALUES
    (1, 'Project 1', 'Description 1', 1000.00, '2023-10-26', '2023-11-26', 1),
    (2, 'Project 2', 'Description 2', 2000.00, '2023-10-27', '2023-11-27', 2),
    (100, 'Project 100', 'Description 100', 5000.00, '2023-10-28', '2023-11-28', 3);
