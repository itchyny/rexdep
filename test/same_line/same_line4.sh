rexdep --pattern '^dependency\("\d+", *"(\d+)"(?:, *"(\d+)")?(?:, *"(\d+)")?\)' \
        --module '^dependency\("(\d+)", *"\d+"(?:, *"\d+")?(?:, *"\d+")?\)' --reverse 2 3
