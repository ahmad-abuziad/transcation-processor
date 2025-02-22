CREATE TABLE products (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    tenant_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    created DATETIME NOT NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

INSERT INTO products (tenant_id, name, created) 
VALUES 
    (1, 'Product 1A', NOW()),
    (1, 'Product 2A', NOW()),
    (1, 'Product 3A', NOW()),
    (1, 'Product 4A', NOW()),
    (1, 'Product 5A', NOW()),

    (2, 'Product 1B', NOW()),
    (2, 'Product 2B', NOW()),
    (2, 'Product 3B', NOW()),
    (2, 'Product 4B', NOW()),
    (2, 'Product 5B', NOW()),

    (3, 'Product 1C', NOW()),
    (3, 'Product 2C', NOW()),
    (3, 'Product 3C', NOW()),
    (3, 'Product 4C', NOW()),
    (3, 'Product 5C', NOW());