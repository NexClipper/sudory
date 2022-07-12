DELETE FROM `template` WHERE uuid='20000000000000000000000000000004';
DELETE FROM `template_command` WHERE uuid='20000000000000000000000000000004';

DELETE FROM `template_recipe` WHERE id=126;

UPDATE `template_command` SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern" : "^."},"chart_name":{"type":"string","pattern" : "^."},"repo_url":{"type":"string","pattern" : "^."},"namespace":{"type":"string","pattern" : "^."},"chart_version":{"type":"string"},"values":{"type":"object","additionalProperties":{"type":"string"}}},"required":["name","chart_name","repo_url","namespace"]}' WHERE uuid='20000000000000000000000000000001';
UPDATE `template_command` SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern" : "^."},"chart_name":{"type":"string","pattern" : "^."},"repo_url":{"type":"string","pattern" : "^."},"namespace":{"type":"string","pattern" : "^."},"chart_version":{"type":"string"},"values":{"type":"object","additionalProperties":{"type":"string"}}},"required":["name","chart_name","repo_url","namespace"]}' WHERE uuid='20000000000000000000000000000003';

UPDATE `template_recipe` SET `args`='{"type":"object","properties":{"name":{"type":"string"},"chart_name":{"type":"string"},"namespace":{"type":"string"},"chart_version":{"type":"string"},"values":{"type":"object","additionalProperties":{"type":"string"}}},"required":["name","chart_name","repo_url","namespace"]}' WHERE id=118;
UPDATE `template_recipe` SET `args`='{"type":"object","properties":{"name":{"type":"string"},"chart_name":{"type":"string"},"namespace":{"type":"string"},"chart_version":{"type":"string"},"values":{"type":"object","additionalProperties":{"type":"string"}}},"required":["name","chart_name","repo_url","namespace"]}' WHERE id=120;
