DELETE FROM `template` WHERE uuid='00000000000000000000000000000040';
DELETE FROM `template_command` WHERE uuid='00000000000000000000000000000040';

DELETE FROM `template_recipe` WHERE method='kubernetes.namespaces.delete.v1';
