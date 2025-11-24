# automate json tag
# 파일 안 모든 구조체 이름 추출
structs=$(grep -E '^type [A-Z][A-Za-z0-9]* struct' ./indicator.go | awk '{print $2}')

# 각 구조체에 json 태그 추가
for s in $structs; do
    gomodifytags \
      -file ./indicator.go \
      -struct $s \
      --add-tags json \
      -w
done
