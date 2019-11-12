#!/bin/bash

for i in "$@"
do
	case $i in
		--dev|--development)
			MODE="development"
			shift
			;;
		--prod|--production)
			MODE="production"
			shift
			;;
		--clean)
			CLEAN=true
			shift
			;;
		-y|--yes)
			YES=true
			shift
			;;
		--help)
			HELP=true
			shift
			;;
		*)
			;;
	esac
done

if [ -z "$MODE" ]
then
	echo "Usage:"
	echo "    $0 [MODE] [OPTIONS]"
	echo ""
	echo "Description:"
	echo "    Convert all mp3 files in \"audio\" directory to dca files for easy streaming to Discord."
	echo ""
	echo "Modes:"
	echo "    --dev             Keep original files (mp3) after conversion."
	echo "    --production      Remove original files (mp3) after conversion."
	echo ""
	echo "Options:"
	echo "    --clean           Remove dca file if corresponding mp3 file is missing."
	exit 0
fi
echo "Running in $MODE mode."

if [ "$MODE" == "production" ] || [ "$CLEAN" == "true" ] && [ "$YES" != "true" ]
then
	echo ""
	echo "!!! THIS WILL REMOVE FILES PERMANENTLY !!!"
	echo "Execute without parameters to view help. Pass \"-y\" flag to skip this notice."
	exit 0
fi

cd audio
for f in *
do
	if [ ! -d "$f" ]; then
		continue
	fi

	if [ "$MODE" == "development" ] && [ "$CLEAN" == "true" ]; then
		echo "clean up: $f"
		for ff in "$f"/*.dca
		do
			[ -e "$ff" ] || continue
			if [ ! -f "${ff%.*}.mp3" ]; then
				echo "clean: $ff"
				rm "$ff"
			fi
		done
	fi

	for ff in "$f"/*.mp3
	do
		[ -e "$ff" ] || echo "nothing to convert: $f"
		[ -e "$ff" ] || continue

		if [ -f "${ff%.*}.dca" ]; then
			echo "exists: ${ff%.*}.dca"
		else
			echo "convert: $ff"
			ffmpeg -i "$ff" -f s16le -ar 48000 -ac 2 pipe:1 2>/dev/null | dca > "${ff%.*}.dca"
		fi

		if [ "$MODE" == "production" ] && [ "$?" -eq 0 ]; then
			echo "remove: $ff"
			rm "$ff"
		fi
	done
done
