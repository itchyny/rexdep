rexdep --pattern '^\d+ +depends +on +(\d+)(?:, *(\d+))?(?:, *(\d+))?' \
        --module '^(\d+) +depends +on +\d+(?:, *\d+)?(?:, *\d+)?' --reverse 1
