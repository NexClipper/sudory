use `sudory_prototype_r1`;

create user 'sudory'@'%' identified by 'sudory';
create user 'sudory'@'localhost' identified by 'sudory';

GRANT ALL PRIVILEGES ON sudory_prototype_r1.* TO 'sudory'@'%';
GRANT ALL PRIVILEGES ON sudory_prototype_r1.* TO 'sudory'@'localhost';