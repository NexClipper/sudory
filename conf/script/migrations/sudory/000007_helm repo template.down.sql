DELETE FROM `template` WHERE uuid IN ('20000000000000000000000000000005', '20000000000000000000000000000006', '20000000000000000000000000000007');
DELETE FROM `template_command` WHERE uuid IN ('20000000000000000000000000000005', '20000000000000000000000000000006', '20000000000000000000000000000007');

DELETE FROM `template_recipe` WHERE id IN (127, 128, 129);

UPDATE `template_command` SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"}},"required":["name","chart_name","repo_url","namespace"]}' WHERE uuid='20000000000000000000000000000001';
UPDATE `template_command` SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"}},"required":["name","chart_name","repo_url","namespace"]}' WHERE uuid='20000000000000000000000000000003';

UPDATE `template_recipe` SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"}},"required":["name","chart_name","repo_url","namespace"]}' WHERE id=118;
UPDATE `template_recipe` SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"}},"required":["name","chart_name","repo_url","namespace"]}' WHERE id=120;
