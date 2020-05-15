CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER SYSTEM SET shared_buffers = '128MB';

CREATE EXTENSION postgis;
CREATE EXTENSION btree_gist;
CREATE TABLE IF NOT EXISTS Cafe
    (
      CafeID               SERIAL PRIMARY KEY,
      CafeName             TEXT,
      Address              TEXT,
      Description          TEXT,
      StaffID              INT,
      OpenTime             TIME,
      CloseTime            TIME,
      Photo                TEXT,
			location GEOMETRY(POINT, 4326),
      location_str text
    );
CREATE INDEX cafe_staff_id_idx on cafe(StaffID);
CREATE INDEX cafe_location_idx ON cafe USING GIST (location);

CREATE TABLE IF NOT EXISTS Staff
(
    StaffID  SERIAL PRIMARY KEY,
    Name     text,
    Email    text UNIQUE,
    Password bytea,
    EditedAt timestamp,
    Photo    text,
    IsOwner  boolean,
    CafeID   integer,
	Position text
);

CREATE TABLE IF NOT EXISTS ApplePass
(
  ApplePassID SERIAL PRIMARY KEY,
  CafeID      int NOT NULL,
  Type        text NOT NULL,
  LoyaltyInfo jsonb,
  published   bool,
  Design      JSONB NOT NULL,
  Icon     	  bytea NOT NULL,
  Icon2x   	  bytea NOT NULL,
  Logo        bytea NOT NULL,
  Logo2x      bytea NOT NULL,
  Strip    	  bytea,
  Strip2x     bytea,
  FOREIGN KEY (CafeID) REFERENCES Cafe (CafeID),
  UNIQUE (CafeID, Type, published)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS Customer
(
    CustomerID           uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
	CafeID               INT,
	Type                 text,
	Points               jsonb,
	surveyresult         jsonb
);

CREATE TABLE IF NOT EXISTS uuidcaferepository
(
uuid varchar(255) PRIMARY KEY,
cafeid integer
);

CREATE TABLE ApplePassMeta (
    CafeID int primary key,
    meta   jsonb
);

CREATE TABLE IF NOT EXISTS surveyTemplate(
    cafeID int references cafe(cafeid) ON DELETE CASCADE UNIQUE,
    surveyTemplate jsonb,
    cafeOwnerId int references staff(staffid) ON DELETE CASCADE
);

CREATE FUNCTION NotEmpty(value1 bytea, value2 bytea) RETURNS bytea AS $$
    BEGIN
        IF length(value1) <> 0 THEN
            RETURN value1;
        end if;
        RETURN value2;
    END;
    $$ LANGUAGE plpgsql;

CREATE FUNCTION NotEmpty(value1 text, value2 jsonb) RETURNS jsonb AS $$
    BEGIN
        IF value1::text <> '' THEN
            RETURN value1::jsonb;
        end if;
        RETURN value2;
    END
    $$ LANGUAGE plpgsql;

CREATE FUNCTION NotEmpty(value1 text, value2 text) RETURNS text AS $$
    BEGIN
        IF value1::text <> '' THEN
            RETURN value1::text;
        end if;
        RETURN value2;
    END
    $$ LANGUAGE plpgsql;



Create TABLE IF NOT EXISTS statistics_table
(
    jsonData   jsonb,
    time       timestamp,
    clientUUID uuid references customer (CustomerID),
    staffId    int,
    cafeId     int
);

CREATE INDEX stat_idx on statistics_table(cafeId,time);
