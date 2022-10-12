DELETE FROM `template` WHERE uuid='00000000000000000000000000000037';
DELETE FROM `template_command` WHERE uuid='00000000000000000000000000000037';

DELETE FROM `template_recipe` WHERE method='kubernetes.secrets.patch.v1';
