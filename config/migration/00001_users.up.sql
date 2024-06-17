BEGIN;

-- Set Timezone
SET TIME ZONE 'Asia/Bangkok';

-- Create the Role table
CREATE TABLE "role" (
    Role_ID SERIAL PRIMARY KEY,
    Role_Title VARCHAR(40) NOT NULL UNIQUE
);

-- Insert the role table
INSERT INTO "role" (
    Role_Title
)
VALUES
    ('Manager'),
    ('Admin'),
    ('Customer');

-- Create the User table
CREATE TABLE "user" (
    User_ID INT PRIMARY KEY,
    User_Role_ID INT REFERENCES "role"(Role_ID) NOT NULL,
    User_Name VARCHAR(40), 
    User_Email VARCHAR(50) UNIQUE NOT NULL,
    User_Password TEXT NOT NULL
);

-- Insert the User table
INSERT INTO "user" (
    User_ID,
    User_Role_ID,
    User_Name,
    User_Email,
    User_Password

)
VALUES
    (1, 1, 'Gordon Freeman', 'gordon_freeman@gmail.com', '$2a$10$DAi7ije26J6vGiZ8EknTK.Go8VsH3/CerbE9QJTEbk3HnF6S0/h9O');

-- Create the Oauth table
CREATE TABLE "oauth" (
    Oauth_ID SERIAL PRIMARY KEY,
    Oauth_User_ID INT REFERENCES "user"(User_ID) NOT NULL,
    Access_Token VARCHAR(300) NOT NULL,
    Reflesh_Token VARCHAR(300) NOT NULL
);

CREATE TABLE "producttype" (
    ProdType_Code INT PRIMARY KEY,
    ProdType_Name VARCHAR(40) NOT NULL
);

INSERT INTO "producttype" (ProdType_Code, ProdType_Name) VALUES (1, 'Food');
INSERT INTO "producttype" (ProdType_Code, ProdType_Name) VALUES (2, 'Drink');
INSERT INTO "producttype" (ProdType_Code, ProdType_Name) VALUES (3, 'Snack');

COMMIT;