DELETE FROM `template` WHERE uuid='00000000000000000000000000000036';
DELETE FROM `template_command` WHERE uuid='00000000000000000000000000000036';

DELETE FROM `template_recipe` WHERE method='kubernetes.services.patch.v1';
