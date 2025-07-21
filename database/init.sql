CREATE TYPE media_types AS ENUM ('Image', 'Video', 'Audio', 'Book');

CREATE TABLE IF NOT EXISTS Users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  pass_hash TEXT NOT NULL,
  is_admin BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS Sessions (
  username TEXT PRIMARY KEY,
  session_id TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS Posts (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  media_type media_types NOT NULL,
  filename TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS UserPosts (
  user_id INT REFERENCES Users(id) ON DELETE CASCADE,
  post_id INT REFERENCES Posts(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, post_id)
);

CREATE INDEX IF NOT EXISTS idx_userposts_user_id ON UserPosts(user_id);
CREATE INDEX IF NOT EXISTS idx_userposts_post_id ON UserPosts(post_id);

CREATE TABLE IF NOT EXISTS UserFavs (
  user_id INT REFERENCES Users(id) ON DELETE CASCADE,
  post_id INT REFERENCES Posts(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, post_id)
);

CREATE INDEX IF NOT EXISTS idx_userfavs_user_id ON UserFavs(user_id);
CREATE INDEX IF NOT EXISTS idx_userfavs_post_id ON UserFavs(post_id);

CREATE TABLE IF NOT EXISTS Tags (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS PostTags (
  post_id INT REFERENCES Posts(id) ON DELETE CASCADE,
  tag_id INT REFERENCES Tags(id) ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_posttags_tag_id ON PostTags(tag_id);
CREATE INDEX IF NOT EXISTS idx_posttags_post_id ON PostTags(post_id);

CREATE TABLE IF NOT EXISTS Artists (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS ArtistPosts (
  artist_id INT REFERENCES Artists(id) ON DELETE CASCADE,
  post_id INT REFERENCES Posts(id) ON DELETE CASCADE,
  PRIMARY KEY (artist_id, post_id)
);

CREATE INDEX IF NOT EXISTS idx_artistposts_artist_id ON ArtistPosts(artist_id);
CREATE INDEX IF NOT EXISTS idx_artistposts_post_id ON ArtistPosts(post_id);