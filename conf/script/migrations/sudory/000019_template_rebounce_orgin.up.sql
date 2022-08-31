-- change system service value 
UPDATE template SET `origin` = 'system', `updated` = NOW() WHERE `uuid` = '99990000000000000000000000000001';
