#!/usr/bin/env bash
# Scrapes code looking for @REF tags and builds a bibliography of quoted code

ref="docs/REFERENCES.md"

rm -f "${ref}"

echo "# References" >> "${ref}"
echo "" >> "${ref}"

find cmd/ -type f -exec sh -c "awk '/@REF/ {printf \"- \"; for(i=3;i<=NF;++i)printf \$i\"\"FS ; print \"\"}' {} | sort -u" \; 2>/dev/null >> "${ref}"
find internal/ -type f -exec sh -c "awk '/@REF/ {printf \"- \"; for(i=3;i<=NF;++i)printf \$i\"\"FS ; print \"\"}' {} | sort -u" \; 2>/dev/null >> "${ref}"
