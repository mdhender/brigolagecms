# Notes

Handy command to find large files after crawling the website.

     find . -type f | xargs ls -l | cut -c31- | sort -n