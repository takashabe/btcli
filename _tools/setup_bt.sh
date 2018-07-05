#!/bin/sh

cbt -project $1 -instance $2 -creds $3 deletetable users
cbt -project $1 -instance $2 -creds $3 deletetable articles
cbt -project $1 -instance $2 -creds $3 createtable users
cbt -project $1 -instance $2 -creds $3 createtable articles
cbt -project $1 -instance $2 -creds $3 createfamily users d
cbt -project $1 -instance $2 -creds $3 createfamily articles d
cbt -project $1 -instance $2 -creds $3 set users 1 d:row=madoka
cbt -project $1 -instance $2 -creds $3 set users 2 d:row=homura
cbt -project $1 -instance $2 -creds $3 set users 3 d:row=sayaka
cbt -project $1 -instance $2 -creds $3 set users 4 d:row=kyouko
cbt -project $1 -instance $2 -creds $3 set users 4 d:row=anko
cbt -project $1 -instance $2 -creds $3 set articles 1##1 d:title=madoka_title
cbt -project $1 -instance $2 -creds $3 set articles 1##1 d:content=madoka_content
cbt -project $1 -instance $2 -creds $3 set articles 2##1 d:title=homura_title
cbt -project $1 -instance $2 -creds $3 set articles 2##1 d:content=homura_content
cbt -project $1 -instance $2 -creds $3 set articles 2##2 d:title=homuhomu_title
cbt -project $1 -instance $2 -creds $3 set articles 2##2 d:content=homuhomu_content
