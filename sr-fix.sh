#!/bin/bash
targetSampleRate=48000
for i in $(ls audio/) ; do
  sr=$(ffprobe -hide_banner -loglevel error -show_entries stream=sample_rate audio/${i} | awk 'FNR == 2 {print}')
  echo "audio/${i} ${sr}"
  if [[ ${sr} != *"${targetSampleRate}"* ]];then
    echo "transcode ${i}, set sample rate ${targetSampleRate}hz"
    mv audio/${i} audio/tempFile.mp3
    ffmpeg -i audio/tempFile.mp3 -hide_banner -loglevel error -af aresample=resampler=soxr -ar ${targetSampleRate} audio/${i}
    rm audio/tempFile.mp3
  fi
done
