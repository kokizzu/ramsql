#!/usr/bin/env bash

# USAGE: ${0} {--lexeme <LEXEME>}+ [--name <NAME>]

# bash best practices
set -o errexit
set -o pipefail
set -o nounset

# set magic variables
_SCRIPTPATH=$(dirname "${0}")
_SCRIPTNAME=$(basename "${0}")

##
## Command-line Options
##

LEXEMES=()
NAME=
POSITIONAL_OPTIONS=()

while [[ $# -gt 0 ]]; do
  key="$1"

  case ${key} in
    --lexeme)
      LEXEMES+=("${2}")
      shift 2
      ;;
    --name)
      NAME="${2}"
      shift 2
      ;;
    *) # unknown option
      POSITIONAL_OPTIONS+=("${key}")
      shift
      ;;
  esac
done

# Check LEXEMES parameter
if [[ ${#LEXEMES[@]} -le 0 ]]; then
  echo "[ERROR] at least one --lexeme <STRING> option is required"
  exit 1
fi

# Check NAME parameter
if [[ "${NAME}" == "" ]]; then
  NAME="$(tr '[:lower:]' '[:upper:]' <<< ${LEXEMES[0]:0:1})${LEXEMES[0]:1}"
fi

##
## Main
##

main() {
  _output_file="lexer.matcher.${NAME}.go"

  echo "package parser" > ${_output_file}
  echo "" >> ${_output_file}
  echo "func (l *lexer) Match${NAME}Token() bool {" >> ${_output_file}

  echo -n "  return" >> ${_output_file}
  for i in $(seq 0 $(( ${#LEXEMES[@]} - 1))); do
    if [[ ${i} -gt 0 ]]; then
      echo -ne " ||\n    " >> ${_output_file}
    fi

    if [[ ${#LEXEMES[${i}]} -gt 1 ]]; then
      echo -n " l.Match([]byte(\"${LEXEMES[${i}]}\"), ${NAME}Token)" >> ${_output_file}
    else
      echo -n " l.MatchSingle('${LEXEMES[${i}]}', ${NAME}Token)" >> ${_output_file}
    fi
  done
  echo >> ${_output_file}

  echo "}" >> ${_output_file}
}

main
