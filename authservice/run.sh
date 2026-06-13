flyway_new() {
  local dir="$(dirname "$0")/src/main/resources/db/migration"
  mkdir -p "$dir"
  local next=$(($(ls "$dir"/V*__*.sql 2>/dev/null | wc -l) + 1))
  local file="$dir/V${next}__${1}.sql"
  touch "$file"
  echo "Created: V${next}__${1}.sql"
}

case "${1:-run}" in
  new)
    shift
    flyway_new "$@"
    ;;
  run|*)
    mvn spring-boot:run
    ;;
esac
