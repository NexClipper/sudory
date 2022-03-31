use `sudory`;

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `sudory` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;

create user 'sudory'@'%' identified by 'sudory';
create user 'sudory'@'localhost' identified by 'sudory';

GRANT ALL PRIVILEGES ON sudory.* TO 'sudory'@'%';
GRANT ALL PRIVILEGES ON sudory.* TO 'sudory'@'localhost';
