#!/usr/bin/env bash

function _wh_completions()
{
  latest="${COMP_WORDS[$COMP_CWORD]}"
  prev="${COMP_WORDS[$COMP_CWORD - 1]}"
  words=""
  case "${prev}" in
    wh)
      words="help report start stop"
      ;;
    report)
      words="day week month year account"
      ;;
    *)
      ;;
  esac
  COMPREPLY=($(compgen -W "$words" -- $latest))
  return 0
}

complete -F _wh_completions wh