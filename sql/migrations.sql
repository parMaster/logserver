DROP TABLE rawdata;

CREATE TABLE rawdata (
	id bigint NOT NULL,
	dt timestamp with time zone NOT NULL,
	topic character varying NOT NULL,
	message character varying NOT NULL
);

CREATE SEQUENCE tempdata_id_seq 
	START WITH 1 
	INCREMENT BY 1 
	NO MINVALUE 
	NO MAXVALUE 
	CACHE 1; 

DROP TABLE tempdata;

CREATE TABLE tempdata (
	id			bigint NOT NULL DEFAULT nextval('tempdata_id_seq'::regclass), 
	year		int NOT NULL,
	month		int NOT NULL,
	day			int NOT NULL,
	hour		int NOT NULL,
	minute		int NOT NULL,
	min_temp	DECIMAL NOT NULL,
	avg_temp	DECIMAL NOT NULL,
	max_temp	DECIMAL NOT NULL,
	strval		character varying NOT NULL,
	sensor		character varying NOT NULL,
	CONSTRAINT tempdata_pk PRIMARY KEY (id) 
);

DELETE from tempdata;