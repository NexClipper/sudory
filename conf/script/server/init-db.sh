#!/usr/bin/env bash

DB_HOST=$1
DB_PORT=$2
DB_SCHEME=$3
SQL_PATH=$4
EXPORT_PATH=$5
ROOT_USERNAME=$6
ROOT_PASSWORD=$7
SERVER_USERNAME=$8
SERVER_PASSWORD=$9

echo DB_HOST=${DB_HOST}
echo DB_PORT=${DB_PORT}
echo DB_SCHEME=${DB_SCHEME}
echo SQL_PATH=${SQL_PATH}
echo EXPORT_PATH=${EXPORT_PATH}
echo ROOT_USERNAME=${ROOT_USERNAME}
echo ROOT_PASSWORD=${ROOT_PASSWORD}
echo SERVER_USERNAME=${SERVER_USERNAME}
echo SERVER_PASSWORD=${SERVER_PASSWORD}

EXPORT_FILE=${EXPORT_PATH}/${DB_SCHEME}_$(date +%Y%m%d%H%M).sql

CMD_PRE="mysql --user=${ROOT_USERNAME} --password=${ROOT_PASSWORD} --host=${DB_HOST} --port=${DB_PORT}"


apk update
apk add mariadb-client


EXISTS=$(${CMD_PRE} --execute "show databases" | grep "${DB_SCHEME}")


# if [[ ${EXISTS} != "" && ${EXPORT_PATH} != "" ]] ; then
# 	echo "=============== start export for backup scheme ==============="
# 	mkdir ${EXPORT_PATH}
# 	CMD=$(mysqldump --user=${ROOT_USERNAME} --password=${ROOT_PASSWORD} --host=${DB_HOST} --port=${DB_PORT} -e --single-transaction -c ${DB_SCHEME} > ${EXPORT_FILE})
# 	echo "=============== complete export for backup scheme ==============="
# fi


if [[ ${EXISTS} != "" ]] ; then
	cat > ${SQL_PATH}.execute <<- EOM
		USE \`${DB_SCHEME}\`;
	EOM
	
	cat ${SQL_PATH}.modify >> ${SQL_PATH}.execute
	
	SQL_PATH="${SQL_PATH}.execute"
else
	cat > ${SQL_PATH}.execute <<- EOM
		CREATE DATABASE IF NOT EXISTS \`${DB_SCHEME}\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
		USE \`${DB_SCHEME}\`;
		CREATE USER IF NOT EXISTS \`${SERVER_USERNAME}\`@\`%\` IDENTIFIED BY '${SERVER_PASSWORD}';
		GRANT ALL PRIVILEGES ON \`${DB_SCHEME}\`.* to \`${SERVER_USERNAME}\`@\`%\`;
	EOM
	
	cat ${SQL_PATH}.create >> ${SQL_PATH}.execute
	cat ${SQL_PATH}.truncate >> ${SQL_PATH}.execute
	cat ${SQL_PATH}.template.insert >> ${SQL_PATH}.execute
	cat ${SQL_PATH}.recipe.insert >> ${SQL_PATH}.execute
	
	SQL_PATH="${SQL_PATH}.execute"
fi

cat ${SQL_PATH} > ${EXPORT_PATH}/${DB_SCHEME}.execute.sql
echo SQL_PATH=${SQL_PATH}
echo ${EXPORT_PATH}/execute.sql


if [ -s "${EXPORT_FILE}" ] ; then
    echo "=============== start import scheme ==============="
    CMD=$(${CMD_PRE} -f --execute "source ${SQL_PATH}")
    echo "=============== complete import scheme ==============="
elif [[ ${EXISTS} == "" || ${EXPORT_PATH} == "" ]] ; then
	echo "=============== start import scheme ==============="
    CMD=$(${CMD_PRE} -f --execute "source ${SQL_PATH}")
    echo "=============== complete import scheme ==============="
else
	echo "Import failed due to schema export failed."
	exit 1
fi
