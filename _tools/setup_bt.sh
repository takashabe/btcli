#!/bin/sh

cbt deletetable users
cbt deletetable articles

cbt createtable users
cbt createtable articles

cbt createfamily users d
cbt createfamily articles d

cbt set users 1 d:row=madoka
cbt set users 2 d:row=homura
cbt set users 3 d:row=sayaka
cbt set users 4 d:row=kyouko
cbt set users 4 d:row=anko

cbt set articles 1##1 d:title=madoka_title
cbt set articles 1##1 d:content=madoka_content
cbt set articles 2##1 d:title=homura_title
cbt set articles 2##1 d:content=homura_content
cbt set articles 2##2 d:title=homuhomu_title
cbt set articles 2##2 d:content=homuhomu_content
