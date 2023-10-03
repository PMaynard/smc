for i in data/ctrees/*.ctrees; do 
  [[ -f "$i" ]] || continue
  # dmtooling sand "$i" "${i%.ctrees}.graphml"
  dmtooling convert "$i" "${i%.ctrees}.dependencymodel"
  dmtooling convert "$i" "${i%.ctrees}.fripp"
  # graphml2gv "${i%.ctrees}.graphml" | dot -Tpdf -o "${i%.ctrees}.pdf"
done
