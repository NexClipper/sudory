#!/bin/bash
cat sudory.sql.create | awk 'BEGIN{FS=OFS="AUTO_INCREMENT="}{if(NF<1){print;next;}
  gsub(/^[0-9]+/,"AUTO_INCREMENT=0",$2);
  print $1 $2
}' > sudory.sql.create
