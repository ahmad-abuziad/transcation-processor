CREATE TABLE tenants (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    created DATETIME NOT NULL
);

INSERT INTO tenants (name, created) 
VALUES 
    ('Tenant A', NOW()),
    ('Tenant B', NOW()),
    ('Tenant C', NOW());