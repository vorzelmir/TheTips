CREATE USER 'repluser'@'%' IDENTIFIED BY 'secret';
GRANT REPLICATION SLAVE ON *.* TO 'repluser'@'%';
CREATE DATABASE primary_db;
