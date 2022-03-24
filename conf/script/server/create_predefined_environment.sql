use sudory;

DELETE FROM `environment`;

INSERT INTO `environment` (`uuid`, `name`, `summary`, `value`, `created`, `updated`, `deleted`) VALUES ('cc6eeb942b9a4a9ca34dc4dfabc54275', 'cluster-token-signature-secret', '클러스터 토큰 시그니처 생성 시크릿 default=\'\'', '', NOW(), NULL, NULL);
INSERT INTO `environment` (`uuid`, `name`, `summary`, `value`, `created`, `updated`, `deleted`) VALUES ('77f7b2aeb0aa4254ad073ae7743291ab', 'client-token-signature-secret', '클라이언트 토큰 시그니처 생성 시크릿 default=\'\'', '', NOW(), NULL, NULL);
INSERT INTO `environment` (`uuid`, `name`, `summary`, `value`, `created`, `updated`, `deleted`) VALUES ('e2db6f6b08e94cb58bc6a35e244aaa29', 'bearer-token-signature-secret', 'bearer 토큰 시그니처 생성 시크릿 default=\'\'', '', NOW(), NULL, NULL);
INSERT INTO `environment` (`uuid`, `name`, `summary`, `value`, `created`, `updated`, `deleted`) VALUES ('af9a14a58b254d13ae69c065a27811b6', 'client-session-expiration-time', '클라이언트 세션 만료 시간 (초) default=\'60\'', '60', NOW(), NULL, NULL);
INSERT INTO `environment` (`uuid`, `name`, `summary`, `value`, `created`, `updated`, `deleted`) VALUES ('75531e760ee6423cb3a050ddcc83e275', 'client-config-poll-interval', '클라이언트 poll interval (초) default=\'15\'', '15', NOW(), NULL, NULL);
INSERT INTO `environment` (`uuid`, `name`, `summary`, `value`, `created`, `updated`, `deleted`) VALUES ('4e55651f63814b648f7284ab9113cf67', 'client-config-loglevel', '클라이언트 log level [\'debug\', \'info\', \'warn\', \'error\', \'fatal\'] default=\'debug\'', 'debug', NOW(), NULL, NULL);
