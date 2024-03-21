#!/bin/bash

function get_commits_per_tag() {
    local commits_per_tag=""
    if [ -n "$2" ]; then
        # Maybe --no-merges could be appropriate
        commits_per_tag=$(git log $1..$2 \
            --numstat \
            --pretty=format:'{"hash":"%H","date":"%ad",%nMSG_TOKEN: %s' \
            $@ | \
            perl -lawne '
                if ( not defined $F[0] ) {
                    # TODO: how can I handle this better ??
                }
                elsif ( $F[0] =~ m/^{"hash":"/ ) {
                    print qq#]},@F#
                } elsif ( $F[0] =~ m/^MSG_TOKEN:$/ ) {
                    shift @F;
                    $sanityzed_msg = join( " ", @F );
                    $sanityzed_msg =~ s#"#\\"#g;
                    print qq#"message":"$sanityzed_msg","changes":[#
                } else {
                    print qq#{"insertions":$F[0],"deletions":$F[1],"path":"$F[2]"},#
                }' | \
            perl -wpe 'END{print "]}"}' | \
            tr -d '\n' | \
            perl -wpe 's#^\[]},\{#[{#g' | \
            # I don't like this solution
            perl -wpe 's#,"path":""#,"path":"#g' | \
            # This either
            perl -wpe 's#""},#"},#g' | \
            perl -wpe 's#(]|}),\s*(]|})#$1$2#g' | \
            # Nor this one
            # tr -d '\\' | \
            perl -wpe 's#{"insertions":-,"deletions":-,"path":"#{"insertions":0,"deletions":0,"path":"#g')
    fi
    echo $commits_per_tag
}

while [[ $# -gt 0 ]]; do
  case $1 in
    --url)
      git clone --bare --single-branch $2 && cd $(basename "$_").git
      shift # past argument
      shift # past value
      ;;
    esac
done

# https://stackoverflow.com/questions/57333971/git-log-numstat-has-weird-data

url=$(git config --get remote.origin.url | sed 's/.git$//')
name=$(echo $url | sed 's/.*\/[^ ]*\/\([^.]*\).*/\1/')

roots=$(git rev-list --max-parents=0 HEAD \
    --pretty=format:'{"hash":"%H","date":%at}' \
    $@ | \
    perl -lawne '
        if (not defined $F[1]) {
            print qq#@F,#
        }' | \
    perl -wpe 'BEGIN{print "["}; END{print "]"}' | \
    tr -d '\n' | \
    perl -wpe 's#(]|}),\s*(]|})#$1$2#g')

tags=($(git for-each-ref --sort=-creatordate | grep -E refs/tags/v[0-9]+.[0-9]+.0$))

for (( i=0; i<${#tags[@]}; i+=3))
do
    # This could be a function (or even another file) and tagName could be local
        tagName=$(echo "${tags[$i+2]}" | perl -wpe 's#refs\/tags\/##g')
        commits=$(get_commits_per_tag ${tags[$i+3]} ${tags[$i]})
        ## Check if directory exists or create it
        echo "{\"name\":\"$tagName\",\"type\":\"${tags[$i+1]}\",\"hash\":\"${tags[$i]}\",\"commits\":[$commits]}" | \
            perl -wpe 's#"commits":\[]},\{#"commits":\[\{#g' > ../../data/$name/$tagName\_shodo_report.json
done

go run ../../../main.go -- $name
