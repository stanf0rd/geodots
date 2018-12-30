-- drop
DROP TABLE IF EXISTS dots;

-- dots
CREATE TABLE dots (
  name       citext     NOT NULL,
  lat        real       NOT NULL,
  lon        real       NOT NULL
);
