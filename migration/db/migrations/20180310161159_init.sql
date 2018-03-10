
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE passengers (
	id SERIAL PRIMARY KEY,
	name TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE enum_driver_state AS ENUM (
	'available',
	'busy'
);

CREATE TABLE drivers (
	id SERIAL PRIMARY KEY,
	name TEXT,
	state enum_driver_state,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE driver_locations (
	id SERIAL PRIMARY KEY,
	driver_id INTEGER,
	lat FLOAT,
	lng FLOAT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS driver_locations;
DROP TABLE IF EXISTS drivers;
DROP TYPE IF EXISTS enum_driver_state;
DROP TABLE IF EXISTS passengers;
