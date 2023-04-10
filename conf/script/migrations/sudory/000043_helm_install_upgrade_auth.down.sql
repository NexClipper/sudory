UPDATE `template_command`         SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE uuid='20000000000000000000000000000001';
UPDATE `template_recipe`          SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE method='helm.install';
UPDATE `template_v2`           SET  `inputs`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE uuid='20000000000000000000000000000001';
UPDATE `template_command_v2`   SET  `inputs`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE name='helminstall';

UPDATE `template_command`         SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"reuse_values":{"type":"boolean"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE uuid='20000000000000000000000000000003';
UPDATE `template_recipe`          SET `args`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"reuse_values":{"type":"boolean"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE method='helm.upgrade';
UPDATE `template_v2`            SET `inputs`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"reuse_values":{"type":"boolean"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE uuid='20000000000000000000000000000003';
UPDATE `template_command_v2`    SET `inputs`='{"type":"object","properties":{"name":{"type":"string","pattern":"^."},"chart_name":{"type":"string","pattern":"^."},"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},"namespace":{"type":"string","pattern":"^."},"chart_version":{"type":"string"},"values":{"type":"object"},"reuse_values":{"type":"boolean"},"timeout":{"type":"integer","minimum":1}},"oneOf":[{"required":["name","chart_name","repo_url","namespace"]},{"required":["name","chart_name","repo_name","namespace"]}]}' WHERE name='helm.upgrade';

UPDATE `template_command`         SET `args`='{"type":"object","properties":{"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},},"required":["repo_url","repo_name"]}' WHERE uuid='20000000000000000000000000000005';
UPDATE `template_recipe`          SET `args`='{"type":"object","properties":{"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},},"required":["repo_url","repo_name"]}' WHERE method='helm.repo_add';
UPDATE `template_v2`            SET `inputs`='{"type":"object","properties":{"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},},"required":["repo_url","repo_name"]}' WHERE uuid='20000000000000000000000000000005';
UPDATE `template_command_v2`    SET `inputs`='{"type":"object","properties":{"repo_url":{"type":"string","pattern":"^."},"repo_name":{"type":"string","pattern":"^."},},"required":["repo_url","repo_name"]}' WHERE name='helm.repo_add';