CREATE TABLE banner_content
(
	id SERIAL,
	title TEXT,
	banner_text TEXT,
	url TEXT,
	feature_id INTEGER,
	PRIMARY KEY(id)
);


CREATE TABLE banner
(
	id INTEGER,
	tag_id INTEGER,
	feature_id INTEGER,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	is_active BOOLEAN,
	PRIMARY KEY(tag_id, feature_id),
	FOREIGN KEY(id) REFERENCES banner_content(id)
);