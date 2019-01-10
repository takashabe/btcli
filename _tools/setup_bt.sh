#!/bin/sh

project=${1:-test-project}
instance=${2:-test-instance}
creds=${3:-dummy}

cbt -project $project -instance $instance -creds $creds deletetable users
cbt -project $project -instance $instance -creds $creds deletetable articles
cbt -project $project -instance $instance -creds $creds createtable users
cbt -project $project -instance $instance -creds $creds createtable articles
cbt -project $project -instance $instance -creds $creds createfamily users d
cbt -project $project -instance $instance -creds $creds createfamily articles d
cbt -project $project -instance $instance -creds $creds set users 1 d:row=madoka
cbt -project $project -instance $instance -creds $creds set users 2 d:row=homura
cbt -project $project -instance $instance -creds $creds set users 3 d:row=sayaka
cbt -project $project -instance $instance -creds $creds set users 4 d:row=kyouko
cbt -project $project -instance $instance -creds $creds set users 4 d:row=anko
cbt -project $project -instance $instance -creds $creds set articles 1##1 d:title=madoka_title
cbt -project $project -instance $instance -creds $creds set articles 1##1 d:content=madoka_content
cbt -project $project -instance $instance -creds $creds set articles 2##1 d:title=homura_title
cbt -project $project -instance $instance -creds $creds set articles 2##1 d:content=homura_content
cbt -project $project -instance $instance -creds $creds set articles 2##2 d:title=homuhomu_title
cbt -project $project -instance $instance -creds $creds set articles 2##2 d:content=homuhomu_content
