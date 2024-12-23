CREATE TABLE IF NOT EXISTS project (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    details TEXT NOT NULL,
    team_id BIGINT NOT NULL,
    CONSTRAINT project_team_id_fk FOREIGN KEY (team_id) REFERENCES team(id) ON DELETE CASCADE
);