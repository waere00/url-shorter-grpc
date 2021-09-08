ALTER USER postgres PASSWORD 'postgres';
ALTER USER postgres WITH CREATEDB;
CREATE TABLE links
(
    url  character varying(255) NOT NULL,
    link character varying(30)  NOT NULL
)
