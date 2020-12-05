# DDLToGo
DDLToGo is a utility to convert the BricolageCMS Postgresql DDL to Go types.

## Useful cleanup commands

    golint ./... | grep underscores | sed -e 's/.*;//' -e 's/should be//' -e 's/^ func/ /' -e 's/^ struct field/ /' -e 's/^ type/ /' | sort -u > vars.txt
    awk '{printf("%d -e 's/%s/%s/g'\n", length($1), $1, $2)}' vars.txt | sort -rn | awk '{printf("   -e '"'"'%s'"'"' \\\n", $3)}' > vars.sed
