# We Think

```mysql
DROP DATABASE IF EXISTS brigolage_cms;
CREATE DATABASE brigolage_cms;
SHOW DATABASES;
DROP USER IF EXISTS 'brigolage'@'localhost';
CREATE USER 'brigolage'@'localhost' IDENTIFIED BY 'daisies.are.crunchy';
SELECT USER, HOST FROM mysql.user;
GRANT ALL PRIVILEGES ON brigolage_cms.* TO 'brigolage'@'localhost';
SHOW GRANTS FOR 'brigolage'@'localhost';
FLUSH PRIVILEGES;
```