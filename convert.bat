@echo off
cd audio

for /d %%d in (*.*) do (
	cd %%d
	for %%f in (*.mp3) do (
		if not exist "%%~nf.dca" (
			ffmpeg.exe -i "%%~nf.mp3" -f s16le -ar 48000 -ac 2 pipe:1 | dca > "%%~nf.dca"
		)
	)
	cd ..
)
