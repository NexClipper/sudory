DELETE FROM `template` WHERE uuid IN ('50000000000000000000000000000013', '50000000000000000000000000000014');
DELETE FROM `template_command` WHERE uuid IN ('50000000000000000000000000000013', '50000000000000000000000000000014');

DELETE FROM `template_recipe` WHERE method IN ('openstack.compute.servers.reboot', 'openstack.compute.servers.resize');