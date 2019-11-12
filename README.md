# NiksiBot

This project is forked from [discordapp/airhornbot](https://github.com/discordapp/airhornbot).

NiksiBot plays sound clips to Discord's voice channels, it can be used like Airhorn Bot to play short clips, or fill the channel with music tracks.

## Installation

**Note!** You should have working Go environment before proceeding.

Installation is easy, just clone the repository and install dependencies with the ``go get .`` command on the project directory. To compile the code to a single file, run ``go build .`` in the same directory. You also need to provide your bot token to the bot, you can get one from [Discord Developer Portal](https://discordapp.com/developers/applications/) if you don't have one.

## Adding sound clips

NiksiBot organizes clips to collections. You should have directory called ``audio``, where each sub-directory represents a collection. Every clip should be in one of those sub-directories.

All clips must be converted to [.dca](https://github.com/bwmarrin/dca) files, this can be done easily with provided convert scripts, just make sure you have [ffmpeg](https://ffmpeg.org/) and [dca command line tool](https://github.com/bwmarrin/dca/tree/master/cmd/dca) installed.

When NiksiBot is started, it automatically builds collections based on the contents on ``audio`` directory. (this differs from Airhorn Bot, where each collection is specified in the code) If audio directory is modified, NiksiBot should be restarted to update collections. Remember, NiksiBot looks only files with ``.dca`` extension.

## Usage

**Start the bot with the following command:**
```
./niksibot -t "BOT_TOKEN"
```

The bot uses queue to manage plays, so every time clip is requested, it is added to the queue. Bot will play clips in order (FIFO) from the queue, and disconnects from voice when the queue exhausts.

**Use the bot with the following commands**:
```
Queue random clip from collection
!<COLLECTION>

Queue specific clip from collection
!<COLLECTION> <CLIP>

Skip currently playing clip
!skip

Disconnect and clear queue
!dd

Display list of recently played clips
!history
```

### RNG4EVER Mode

When set in ``RNG4EVER`` mode, the bot will play clips from collection until disconnected with command. Other clips can still be queued, and those are prioritized over random clips. The bot can be set in ``RNG4EVER`` mode with command:
```
!<COLLECTION> rng4ever
```
