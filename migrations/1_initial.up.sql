
CREATE EXTENSION IF NOT EXISTS postgis SCHEMA public;
BEGIN;

CREATE TABLE clients (
	id UUID NOT NULL UNIQUE PRIMARY KEY,
	name VARCHAR,
	email VARCHAR NOT NULL UNIQUE,
	password VARCHAR NOT NULL,
	loc_id VARCHAR,
	point public.geography(POINT,4326),
	loc_name VARCHAR,
	loc_type VARCHAR
);

CREATE TABLE locations (

	loc_id VARCHAR,
	point public.geography(POINT,4326),
	loc_name VARCHAR,
	loc_type VARCHAR
);

END;