use `sudory`;
create user 'sudory'@'%' identified by 'sudory';
create user 'sudory'@'localhost' identified by 'sudory';

GRANT ALL PRIVILEGES ON sudory.* TO 'sudory'@'%';
GRANT ALL PRIVILEGES ON sudory.* TO 'sudory'@'localhost';
