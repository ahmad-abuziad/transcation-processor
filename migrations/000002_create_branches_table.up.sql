CREATE TABLE branches (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    tenant_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    created DATETIME NOT NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE
);

INSERT INTO branches (tenant_id, name, created) 
VALUES 
    (1, 'Branch 1A', NOW()),
    (1, 'Branch 2A', NOW()),
    (1, 'Branch 3A', NOW()),
    (1, 'Branch 4A', NOW()),
    (1, 'Branch 5A', NOW()),
    
    (2, 'Branch 1B', NOW()),
    (2, 'Branch 2B', NOW()),
    (2, 'Branch 3B', NOW()),
    (2, 'Branch 4B', NOW()),
    (2, 'Branch 5B', NOW()),
    
    (3, 'Branch 1C', NOW()),
    (3, 'Branch 2C', NOW()),
    (3, 'Branch 3C', NOW()),
    (3, 'Branch 4C', NOW()),
    (3, 'Branch 5C', NOW());