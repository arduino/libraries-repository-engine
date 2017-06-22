#!/bin/bash

ISSUE=$1
REPO=arduino/Arduino

if [ -z $ISSUE ]; then
	echo Please specify the issue number
	echo $0 6543
	exit 1
fi

REPO_URL=`curl -s -f https://api.github.com/repos/$REPO/issues/$ISSUE | jq -r .body | tr "\\r\\n" " " | egrep -o 'https?://[^ )]+' | head -n1`.git

echo Detected repo url is: $REPO_URL
echo Press enter to continue
read


TMP=`mktemp -u`

git clone $REPO_URL $TMP
if [ $? != 0 ]; then
	echo "Failed to fetch $REPO_URL."
	exit 1
fi
cd $TMP
TAGS=`git tag`
echo TAGS=$TAGS

if [ -z "$TAGS" ]; then
	echo "ERROR: Failed to detect TAGS."
	exit 2
fi

git config advice.detachedHead false

for TAG in `git tag`; do
	git checkout "$TAG"
	N=`cat library.properties | grep "name="`
	if [ $? != 0 ]; then
		echo "Invalid library.properties in tag $TAG"
		continue
	fi
	NAME=${N#*=}
	echo "Found name: $NAME"
done

if [ -z "$NAME" ]; then
	echo "ERROR: Failed to detect library NAME."
	exit 2
fi

echo ""
echo ">> Library content:"
ls -l

cd -

rm -rf $TMP

echo ""
echo ">> possible duplicates on repositories.txt:"
grep "$NAME" -i repositories.txt

ADDED_LINE="$REPO_URL|Contributed|$NAME"

echo ""
echo ">> Line to be added:"
echo $ADDED_LINE
echo Press enter to proceed
read

echo $ADDED_LINE >> repositories.txt
git add -u repositories.txt
echo -e "Added lib $NAME\n\nhttps://github.com/arduino/Arduino/issues/$ISSUE" | git commit -F -

