CREATE TABLE driver_locations (
	id uuid,
	driver_id int,
	created_at bigint,
	lat double,
	lng double,
	PRIMARY KEY ((driver_id), created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);
