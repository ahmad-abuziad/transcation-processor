CREATE TABLE sales_transactions (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    tenant_id INTEGER NOT NULL,
    branch_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity_sold INTEGER NOT NULL,
    price_per_unit DECIMAL(10,2) NOT NULL,
    log_timestamp TIMESTAMP NOT NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
