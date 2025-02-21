CREATE DATABASE transaction_processor_db;
CREATE USER 'dbuser'@'%' IDENTIFIED BY 'pass';
GRANT ALL PRIVILEGES ON transaction_processor_db.* TO 'dbuser'@'%';
FLUSH PRIVILEGES;