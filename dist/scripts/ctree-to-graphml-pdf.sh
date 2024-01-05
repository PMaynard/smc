for i in data/ctrees/*.ctrees; do 
  [[ -f "$i" ]] || continue
  # smc sand "$i" "${i%.ctrees}.graphml"
  smc convert "$i" "${i%.ctrees}.dependencymodel"
  smc convert "$i" "${i%.ctrees}.fripp"
  # graphml2gv "${i%.ctrees}.graphml" | dot -Tpdf -o "${i%.ctrees}.pdf"
done
