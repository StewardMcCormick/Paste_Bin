ALTER TABLE api_key
ADD COLUMN key_prefix varchar(15) NOT NULL default 'pb_live'
