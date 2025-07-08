CREATE TABLE IF NOT EXISTS Users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  access_level INT NOT NULL -- 0 general user, 1 admin
);

CREATE TABLE IF NOT EXISTS Content (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  media_type INT NOT NULL, -- 0 image
  file_name VARCHAR(255) NOT NULL,
  artist VARCHAR(255) NULL
);

CREATE TABLE IF NOT EXISTS Tags (
  tag_name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS PostTags (
  post_id INT REFERENCES Content(id),
  tag_name VARCHAR(255) REFERENCES Tags(tag_name)
);

INSERT INTO Tags VALUES ('fighting');
INSERT INTO Tags VALUES ('walking');
INSERT INTO Tags VALUES ('city');
INSERT INTO Tags VALUES ('beach');
INSERT INTO Tags VALUES ('forest');
INSERT INTO Tags VALUES ('tattoo');
INSERT INTO Tags VALUES ('rainy');
INSERT INTO Tags VALUES ('night');
INSERT INTO Tags VALUES ('sunny');
