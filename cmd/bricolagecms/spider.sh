#!/bin/bash
###########################################################################
# brigolagecms/cmd/bricolagecms/spider.sh
#
# Copyright (c) 2020 Michael Henderson
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
#  SOFTWARE.
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
