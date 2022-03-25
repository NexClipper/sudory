--	https://docs.google.com/spreadsheets/d/1-Knmp6SEO_uKsJYztUV9Qlgywo-tazU3TGlbhacev-w/edit#gid=0
	use sudory;
--	environment
	REPLACE INTO environment (`uuid`, `name`, `summary`, `value`, `id`, `created`) VALUES ('e2db6f6b08e94cb58bc6a35e244aaa29', 'bearer-token-signature-secret', 'bearer-토큰 시그니처 생성 시크릿 default=\'\'', '', '1', '2022-03-28 13:41:47');
	REPLACE INTO environment (`uuid`, `name`, `summary`, `value`, `id`, `created`) VALUES ('0f5658f37f2b45d881f19c7f56ea2e23', 'bearer-token-expiration-time', 'bearer-토큰 만료 시간 (month) default=\'6\'', '', '2', '2022-03-28 13:41:47');
	REPLACE INTO environment (`uuid`, `name`, `summary`, `value`, `id`, `created`) VALUES ('77f7b2aeb0aa4254ad073ae7743291ab', 'client-session-signature-secret', '클라이언트 세션 시그니처 생성 시크릿 default=\'\'', '', '3', '2022-03-28 13:41:47');
	REPLACE INTO environment (`uuid`, `name`, `summary`, `value`, `id`, `created`) VALUES ('af9a14a58b254d13ae69c065a27811b6', 'client-session-expiration-time', '클라이언트 세션 만료 시간 (초) default=\'60\'', '60', '4', '2022-03-28 13:41:47');
	REPLACE INTO environment (`uuid`, `name`, `summary`, `value`, `id`, `created`) VALUES ('75531e760ee6423cb3a050ddcc83e275', 'client-config-poll-interval', '클라이언트 poll interval (초) default=\'15\'', '15', '5', '2022-03-28 13:41:47');
	REPLACE INTO environment (`uuid`, `name`, `summary`, `value`, `id`, `created`) VALUES ('4e55651f63814b648f7284ab9113cf67', 'client-config-loglevel', '클라이언트 log level [\'debug\', \'info\', \'warn\', \'error\', \'fatal\'] default=\'debug\'', 'debug', '6', '2022-03-28 13:41:47');
	