#!/bin/bash
###########################################################################
# simple script to clone the BricolageCMS website
###########################################################################
ADJUST_EXTENSION=--adjust-extension               # Save files with .html on the end.
CONVERT_LINKS=--convert-links                     # Update links to still work in the static version.
DOMAINS=--domains                                 # Do not follow links outside this domain.
NO_PARENT=--no-parent                             # Don't follow links outside the directory you pass in.
PAGE_REQUISITES=--page-requisites                 # Get all assets/elements (CSS/JS/images).
RECURSIVE=--recursive                             # Download the whole site.
RESTRICT_FILE_NAMES=--restrict-file-names=windows # Modify filenames to work in Windows as well.
#SPAN_HOSTS=--span-hosts                           # Include necessary assets from offsite as well.

###########################################################################
#
wget \
  $RECURSIVE \
  $PAGE_REQUISITES \
  $ADJUST_EXTENSION \
  $SPAN_HOSTS \
  $CONVERT_LINKS \
  $RESTRICT_FILE_NAMES \
  $DOMAINS bricolagecms.org \
  $NO_PARENT \
  bricolagecms.org
