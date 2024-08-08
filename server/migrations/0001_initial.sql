DROP TABLE IF EXISTS ferment_run CASCADE;
DROP TABLE IF EXISTS temp_log;

SET timezone = 'America/New_York';

CREATE TABLE ferment_run (
    ferment_run_id INT GENERATED ALWAYS AS IDENTITY,
    name VARCHAR,
    start_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(ferment_run_id)
);

CREATE INDEX ferment_run_name ON ferment_run(name);

CREATE TABLE temp_log (
    temp_log_id INT GENERATED ALWAYS AS IDENTITY,
    ferment_run_id INT NOT NULL,
    temp DECIMAL(5,3) NOT NULL,
    time_stamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY(temp_log_id),
    CONSTRAINT fk_ferment_run FOREIGN KEY(ferment_run_id) REFERENCES ferment_run(ferment_run_id) ON DELETE CASCADE
);
